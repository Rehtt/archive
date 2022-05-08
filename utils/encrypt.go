/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/8 下午 02:25
 */

package utils

import (
	"archive/utils/aes"
	"archive/utils/rsa"
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/big"
)

var (
	publicKey  []byte
	privateKey []byte
)

const separate = "Aes!"
const size = 496

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
	for w.buf.Len() > size {
		out, err := aes.AesCBCEncrypt(w.buf.Next(size), w.Key, w.Iv)
		if err != nil {
			return 0, err
		}
		_, err = w.w.Write(out)
		if err != nil {
			return 0, err
		}
	}
	return len(p), err

	//out, err := aes.Encrypt(p, w.Key, w.Iv)
	//if err != nil {
	//	return
	//}
	//_, err = w.w.Write(out)
	//if err != nil {
	//	return 0, err
	//}
	//_, err = w.w.Write([]byte(separate))
	//n = len(p)
	//return
}
func (w *Writer) Close() error {
	if w.buf.Len() != 0 {
		var a, b byte
		if w.buf.Len() != size {
			s := size - w.buf.Len()
			if s < 256 {
				a = byte(s)
			} else {
				a = byte(255)
				b = byte(s - 255)
			}
			w.buf.Write(bytes.Repeat([]byte("0"), s))

		}
		out, err := aes.AesCBCEncrypt(w.buf.Bytes(), w.Key, w.Iv)
		if err != nil {
			return err
		}
		_, err = w.w.Write(out)
		if err != nil {
			return err
		}
		w.w.Write([]byte{a, b})
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
	if r.Key == nil {
		return r.r.Read(p)
	}

	size := len(p)
	if r.buf.Len() < size {
		buf := make([]byte, 512)
		var temp bytes.Buffer
		for r.buf.Len() < size {
			nn, err := io.ReadFull(r.r, buf)
			if err != nil {
				return 0, err
			}
			switch nn {
			case 512:
				out, err := aes.AesCBCDncrypt(buf, r.Key, r.Iv)
				if err != nil {
					if temp.Len() != 0 {
						return 0, err
					}
					temp.Write(out)
					continue
				}
				r.buf.Write(out)
			case 2:
				out, err := aes.AesCBCDncrypt(temp.Next(int(buf[0]+buf[1])), r.Key, r.Iv)
				if err != nil {
					return 0, err
				}
				r.buf.Write(out)
			default:
				return 0, errors.New("file error")
			}
		}
	}

	//size := len(p)
	//var buf bytes.Buffer
	//for r.buf.Len() < size {
	//	temp, _ := r.r.ReadBytes(separate[3])
	//	tempSize := len(temp)
	//	if tempSize == 0 {
	//		return 0, errors.New("[1]文件错误")
	//	}
	//	buf.Write(temp)
	//	if tempSize > 4 && string(temp[tempSize-4:]) == separate {
	//		out, err := aes.Dncrypt(buf.Next(buf.Len()-4), r.Key, r.Iv)
	//		if err != nil {
	//			return 0, err
	//		}
	//		r.buf.Write(out)
	//		buf.Reset()
	//	}
	//}
	n = copy(p, r.buf.Next(size))
	return
}

func NewReader(r io.Reader) *Reader {
	reader := &Reader{}
	buf := bufio.NewReader(r)
	m, err := buf.Peek(512)
	if err == nil {
		out, err := rsa.Decrypt(m, privateKey)
		if err == nil {
			err = json.Unmarshal(out, &reader.password)
			if err == nil {
				buf.Discard(512)
			}
		}
	}
	reader.r = buf
	return reader
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
