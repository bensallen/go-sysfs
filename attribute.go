package sysfs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Attribute struct {
	Path string
	File *os.File
}

func (attrib *Attribute) Exists() bool {
	f, _ := fileExists(attrib.Path)
	return f
}

func (attrib *Attribute) Open(flag int, perm os.FileMode) (err error) {
	attrib.File, err = os.OpenFile(attrib.Path, flag, perm)
	return err
}

func (attrib *Attribute) OpenRW() (err error) {
	return attrib.Open(os.O_RDWR|syscall.O_NONBLOCK, 0666)
}

func (attrib *Attribute) OpenRO() (err error) {
	return attrib.Open(os.O_RDONLY|syscall.O_NONBLOCK, 0444)
}

func (attrib *Attribute) Close() (err error) {
	err = attrib.File.Close()
	attrib.File = nil
	return err
}

func (attrib *Attribute) Ioctl(request, arg uintptr) (result uintptr, errno syscall.Errno, err error) {
	if attrib.File == nil {
		err = attrib.OpenRW()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	result, _, errno = syscall.Syscall(syscall.SYS_IOCTL, attrib.File.Fd(), request, arg)
	return result, errno, err
}

func (attrib *Attribute) Read() (str string, err error) {
	if attrib.File == nil {
		err = attrib.OpenRO()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	attrib.File.Seek(0, os.SEEK_SET)
	data, err := ioutil.ReadAll(attrib.File)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(data)), nil
}

func (attrib *Attribute) Write(value string) (err error) {
	if attrib.File == nil {
		err = attrib.OpenRW()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	attrib.File.Seek(0, os.SEEK_SET)
	_, err = attrib.File.WriteString(value)
	return err
}

func (attrib *Attribute) Print(value interface{}) (err error) {
	if attrib.File == nil {
		err = attrib.OpenRO()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	attrib.File.Seek(0, os.SEEK_SET)
	_, err = fmt.Fprint(attrib.File, value)
	return err
}

func (attrib *Attribute) Scan(value interface{}) (err error) {
	if attrib.File == nil {
		err = attrib.OpenRO()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	attrib.File.Seek(0, os.SEEK_SET)
	_, err = fmt.Fscan(attrib.File, value)
	return err
}

func (attrib *Attribute) Printf(format string, args ...interface{}) (err error) {
	if attrib.File == nil {
		err = attrib.OpenRO()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	attrib.File.Seek(0, os.SEEK_SET)
	_, err = fmt.Fprintf(attrib.File, format, args...)
	return err
}

func (attrib *Attribute) Scanf(format string, args ...interface{}) (err error) {
	if attrib.File == nil {
		err = attrib.OpenRO()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	attrib.File.Seek(0, os.SEEK_SET)
	_, err = fmt.Fscanf(attrib.File, format, args...)
	return err
}

func (attrib *Attribute) ReadAllBytes() (data []byte, err error) {
	if attrib.File == nil {
		err = attrib.OpenRO()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	attrib.File.Seek(0, os.SEEK_SET)
	return ioutil.ReadAll(attrib.File)
}

func (attrib *Attribute) WriteBytes(data []byte) (err error) {
	if attrib.File == nil {
		err = attrib.OpenRW()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	_, err = attrib.File.WriteAt(data, 0)
	return err
}

func (attrib *Attribute) ReadBytes(offset int64, count int) (data []byte, n int, err error) {
	if attrib.File == nil {
		err = attrib.OpenRO()
		if err != nil {
			return
		}
		defer func() {
			e := attrib.Close()
			if err == nil {
				err = e
			}
		}()
	}
	data = make([]byte, count)
	n, err = attrib.File.ReadAt(data, offset)
	return data, n, err
}

func (attrib *Attribute) ReadByte() (byte, error) {
	data, _, err := attrib.ReadBytes(0, 1)
	return data[0], err
}

func (attrib *Attribute) WriteByte(value byte) (err error) {
	return attrib.WriteBytes([]byte{value})
}

func (attrib *Attribute) ReadInt() (value int, err error) {
	s, err := attrib.Read()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(s))
}

func (attrib *Attribute) WriteInt(value int) (err error) {
	return attrib.Write(strconv.Itoa(value))
}
