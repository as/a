package fs

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

const MaxFilename = 65536 // We're in Windows and Linux country

var (
	ErrNameLong = fmt.Errorf("name too long (>%d)", MaxFilename)
)

type info struct {
	Size  int64
	Mode  os.FileMode
	Isdir bool
}
type remoteFileInfo struct {
	info
	name string
	mod  time.Time
}

func (r *remoteFileInfo) String() string {
	return fmt.Sprintf("%s	%s	%d	%s", r.info.Mode, r.mod, r.info.Size, r.name)
}

func (r *remoteFileInfo) clone(fi os.FileInfo) {
	if fi == nil {
		return
	}
	r.info.Size = fi.Size()
	r.info.Mode = fi.Mode()
	r.mod = fi.ModTime()
	r.info.Isdir = fi.IsDir()
	r.name = fi.Name()
}

func (r *remoteFileInfo) WriteBinary(dst io.Writer) (err error) {
	if err = binary.Write(dst, binary.BigEndian, r.info); err != nil {
		return err
	}

	if err = writeString(dst, r.name); err != nil {
		return err
	}

	bintime, err := r.mod.MarshalBinary()
	if err != nil {
		return err
	}
	return writeBytes(dst, bintime)
}

func (r *remoteFileInfo) ReadBinary(src io.Reader) error {
	err := binary.Read(src, binary.BigEndian, &r.info)
	if err != nil {
		return err
	}

	name, err := readString(src, MaxFilename)
	if err != nil {
		if err == ErrStrlen {
			return ErrNameLong
		}
		return err
	}
	r.name = name

	modtime, err := readString(src, 64)
	if err != nil {
		return err
	}

	if err = r.mod.UnmarshalBinary([]byte(modtime)); err != nil {
		return err
	}

	return err
}

func (r *remoteFileInfo) Name() string       { return r.name }
func (r *remoteFileInfo) Size() int64        { return r.info.Size }
func (r *remoteFileInfo) Mode() os.FileMode  { return r.info.Mode }
func (r *remoteFileInfo) ModTime() time.Time { return r.mod }
func (r *remoteFileInfo) IsDir() bool        { return r.info.Isdir }
func (r *remoteFileInfo) Sys() interface{}   { return nil }
