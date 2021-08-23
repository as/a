package fs

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
)

type Server struct {
	Local
	fd             net.Listener
	donec, donesrv chan bool
}

type client struct {
	conn net.Conn
	bio  *bufio.ReadWriter
	rx   chan []byte
	tx   chan []byte
}

func (c *client) Write(p []byte) (n int, err error) {
	return c.bio.Write(p)
}
func (c *client) Read(p []byte) (n int, err error) {
	return c.bio.Read(p)
}
func (c *client) Flush() error {
	err := c.bio.Flush()
	if err != nil {
		log.Printf("flush: %s\n", err)
	}
	return err
}

func Serve(netw, addr string) (*Server, error) {
	fd, err := net.Listen(netw, addr)
	if err != nil {
		return nil, err
	}
	s := &Server{
		Local:   Local{},
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
				return
			default:
				conn, err := fd.Accept()
				if err != nil {
					log.Printf("accept: %s\n", err)
					continue
				}
				bio := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
				go s.handle(&client{conn, bio, make(chan []byte), make(chan []byte)})
			}
		}
	}()
	return s, nil
}

func (s *Server) Close() error {
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

func (s *Server) handle(c *client) {
	defer c.conn.Close()
	for {
		select {
		case <-s.donesrv:
			return
		default:
		}
		hdr := make([]byte, 3)
		_, err := io.ReadAtLeast(c, hdr, len(hdr))
		if err != nil {
			return
			log.Printf("invalid header: %s", err)
		}
		switch string(hdr) {
		case "Get", "Put", "Cmd", "Sta":
			ln, err := c.bio.ReadSlice('\n')
			if err != nil {
				log.Printf("readslice: %s\n", err)
				break
			}
			ln = bytes.TrimSpace(ln)
			switch string(hdr) {
			case "Sta":
				fi, err := s.Local.Stat(string(ln))
				if err != nil {
					log.Printf("sta: %s\n", err)
				}

				r := new(remoteFileInfo)
				r.clone(fi)

				err = r.WriteBinary(c)
				if err != nil {
					log.Printf("sta: writebinary: %s", err)
				}
				c.Flush()

			case "Get":
				data, err := s.Local.Get(string(ln))
				if err != nil {
					log.Printf("get: %s\n", err)
				}

				err = writeBytes(c, data)
				c.Flush()

			case "Put":

				data, err := readBytes(c, 1e12)
				if err != nil {
					log.Printf("put: data read err: %s\n", err)
				}

				log.Printf("server: put %q with data %q\n", string(ln), data)
				err = s.Local.Put(string(ln), data)
				if err != nil {
					log.Printf("put: local: %s", err)
				}
			}
		default:
			log.Printf("bad cmd: %s", hdr)
			return
		}
	}
}
