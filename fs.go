package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
)

type FS interface {
	Get(name string) (data []byte, err error)
	Put(name string, data []byte) (err error)
	Cmd(ctx context.Context, name string, arg ...string) (*exec.Cmd, error)
}

type LocalFS struct {
}

func (f *LocalFS) Get(name string) (data []byte, err error) {
	return ioutil.ReadFile(name)
}
func (f *LocalFS) Put(name string, data []byte) (err error) {
	return ioutil.WriteFile(name, data, 0666)
}
func (f *LocalFS) Cmd(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	return exec.CommandContext(ctx, name, arg...), nil
}

type ServeFS struct {
	LocalFS
	fd             net.Listener
	donec, donesrv chan bool
}

func Serve(netw, addr string) (*ServeFS, error) {
	fd, err := net.Listen(netw, addr)
	if err != nil {
		return nil, err
	}
	s := &ServeFS{
		LocalFS: LocalFS{},
		fd:      fd,
		donec:   make(chan bool, 1),
		donesrv: make(chan bool),
	}
	s.donec <- true

	//	go s.run()
	go func() {
		for {
			select {
			case <-s.donesrv:
				break
			default:
				conn, err := fd.Accept()
				if err != nil {
					log.Printf("accept: %s\n", err)
					continue
				}
				go s.handle(&client{conn, make(chan []byte), make(chan []byte)})
			}
		}
	}()
	return s, nil
}

type client struct {
	conn net.Conn
	rx   chan []byte
	tx   chan []byte
}

func (s *ServeFS) Close() error {
	select {
	case ok := <-s.donec:
		if ok {
			close(s.donec)
			close(s.donesrv)
		}
	default:
	}
	return nil
}

func (s *ServeFS) handle(c *client) {
	bio := bufio.NewReader(c.conn)
	defer c.conn.Close()
	for {
		hdr := make([]byte, 3)
		_, err := io.ReadAtLeast(bio, hdr, len(hdr))
		if err != nil {
			log.Printf("invalid header: %s", err)
		}
		select {
		case <-s.donesrv:
			return
		default:
		}
		switch string(hdr) {
		case "Get", "Put", "Cmd":
			ln, err := bio.ReadSlice('\n')
			if err != nil {
				log.Printf("readslice: %s\n", err)
				break
			}
			ln = bytes.TrimSpace(ln)
			switch string(hdr) {
			case "Get":
				data, err := s.LocalFS.Get(string(ln))
				if err != nil {
					log.Printf("get: %s\n", err)
				}

				err = binary.Write(c.conn, binary.BigEndian, int64(len(data)))
				if err != nil {
					log.Printf("get: write len: %s\n", err)
				}

				_, err = c.conn.Write(data)
				if err != nil {
					log.Printf("get: write: %s\n", err)
				}

			case "Put":
				n := int64(0)
				err = binary.Read(bio, binary.BigEndian, &n)

				if err != nil {
					log.Printf("put: %s\n", err)
				}

				if n < 0 {
					log.Printf("put: len<0\n")
					return
				}

				data, err := ioutil.ReadAll(io.LimitReader(bio, n))
				if err != nil {
					log.Printf("put: data read err: %s\n", err)
				}

				err = s.LocalFS.Put(string(ln), data)
				if err != nil {
					log.Printf("put: local: %s", err)
					return
				}

			}
		default:
			log.Printf("bad cmd: %s", hdr)
		}
	}
}
