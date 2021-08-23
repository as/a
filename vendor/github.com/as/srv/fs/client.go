package fs

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
)

type Client struct {
	conn net.Conn
	bio  *bufio.ReadWriter
}

func (c *Client) WriteHeader(msg string, name string) error {
	log.Printf("writeheader: msg %q, name %q", msg, name)
	_, err := c.Write([]byte(msg + name + "\n"))
	return err
}

func (c *Client) Write(p []byte) (n int, err error) {
	return c.bio.Write(p)
}
func (c *Client) Read(p []byte) (n int, err error) {
	return c.bio.Read(p)
}
func (c *Client) Flush() error {
	return c.bio.Flush()
}
func (c *Client) Close() error {
	return c.conn.Close()
}

func Dial(netw, addr string) (*Client, error) {
	conn, err := net.Dial(netw, addr)
	if err != nil {
		log.Printf("dial: %s\n", err)
		return nil, err
	}
	return &Client{
		conn: conn,
		bio: bufio.NewReadWriter(
			bufio.NewReader(conn),
			bufio.NewWriter(conn),
		),
	}, nil
}

func (f *Client) Stat(name string) (fi os.FileInfo, err error) {
	f.WriteHeader("Sta", name)

	if err = f.Flush(); err != nil {
		return nil, err
	}

	r := &remoteFileInfo{}
	return r, r.ReadBinary(f)
}
func (f *Client) Get(name string) (data []byte, err error) {
	f.WriteHeader("Get", name)
	if err = f.Flush(); err != nil {
		return nil, err
	}

	n := int64(0)
	err = binary.Read(f, binary.BigEndian, &n)
	if err != nil {
		log.Printf("get: write len: %s\n", err)
	}
	return ioutil.ReadAll(io.LimitReader(f, n))
}
func (f *Client) Put(name string, data []byte) (err error) {
	println("put", name)
	f.WriteHeader("Put", name)

	if err = writeBytes(f, data); err != nil {
		defer f.Flush()
		return err
	}

	return f.Flush()
}
func (f *Client) Cmd(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	return nil, nil
}
