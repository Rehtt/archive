/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/8 下午 02:25
 */

package utils

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"github.com/Rehtt/archive/utils/aes"
	"github.com/Rehtt/archive/utils/rsa"
	"io"
	"log"
	"math/big"
)

var (
	publicKey  []byte
	privateKey []byte
)

const Size = 512

type password struct {
	Iv   []byte `json:"iv"`
	Key  []byte `json:"key"`
	Salt []byte `json:"salt"`
}

type Writer struct {
	w   io.Writer
	buf bytes.Buffer
	password
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if w.Key == nil {
		return w.w.Write(p)
	}
	w.buf.Write(p)
	for w.buf.Len() > Size {
		out, err := aes.AesCBCEncrypt(w.buf.Next(Size), w.Key, w.Iv)
		if err != nil {
			return 0, err
		}
		_, err = w.w.Write(out)
		if err != nil {
			return 0, err
		}
	}
	return len(p), err
}
func (w *Writer) Close() error {
	if w.buf.Len() != 0 {
		var a, b, c byte
		if w.buf.Len() != Size {
			s := Size - w.buf.Len()
			if s < 256 {
				a = byte(s)
			} else {
				a = byte(255)
				b = byte(s) - a
				if int(b)-255 > 0 {
					c = b - 255
				}
			}
			w.buf.Write(random(s))

		}
		out, err := aes.AesCBCEncrypt(w.buf.Bytes(), w.Key, w.Iv)
		if err != nil {
			return err
		}
		_, err = w.w.Write(out)
		if err != nil {
			return err
		}
		w.w.Write([]byte{a, b, c})
	}
	return nil
}

func NewWriter(w io.Writer) (*Writer, error) {
	passwd := password{
		Iv:   random(16),
		Key:  random(32),
		Salt: random(16),
	}
	data, _ := json.Marshal(passwd)
	out, err := rsa.Encrypt(data, publicKey)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(out)
	if err != nil {
		return nil, err
	}
	return &Writer{w: w, password: passwd}, nil
}

type Reader struct {
	r   *bufio.Reader
	buf bytes.Buffer
	password
}

func (r *Reader) Read(p []byte) (n int, err error) {
	size := len(p)
	if r.buf.Len() < size {
		buf := make([]byte, Size)
		for r.buf.Len() < size {
			nn, err := io.ReadFull(r.r, buf)
			if err != nil {
				return 0, err
			}
			out, err := aes.AesCBCDncrypt(buf, r.Key, r.Iv)
			if err != nil {
				return 0, err
			}
			temp, err := r.r.Peek(4)
			if len(temp) != 0 && err != nil {
				nn = Size - int(temp[0]) + int(temp[1]) + int(temp[2])
				err = nil
			}
			r.buf.Write(out[:nn])
		}
	}
	n = copy(p, r.buf.Next(size))
	return
}

func NewReader(r io.Reader) (*Reader, error) {
	buf := make([]byte, 512)
	r.Read(buf)
	out, err := rsa.Decrypt(buf, privateKey)
	if err != nil {
		return nil, err

	}

	reader := &Reader{r: bufio.NewReader(r)}
	err = json.Unmarshal(out, &reader.password)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func InitEncrypt(private, public []byte) {
	privateKey = private
	publicKey = public
}

func random(l int) []byte {
	var buf bytes.Buffer
	for buf.Len() < l {
		n, err := rand.Int(rand.Reader, big.NewInt(256))
		if err != nil {
			log.Fatalln(err)
		}
		buf.WriteByte(byte(n.Int64()))
	}
	return buf.Next(l)
}
