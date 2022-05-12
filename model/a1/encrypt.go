/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/11 21:51
 */

package a1

import (
	"bytes"
	"encoding/json"
	"github.com/Rehtt/archive/utils"
	"github.com/Rehtt/archive/utils/aes"
	"github.com/Rehtt/archive/utils/rsa"
	"io"
)

type encrypt struct {
	w   io.Writer
	buf bytes.Buffer
	a1  *A1
}

func (e *encrypt) Write(p []byte) (n int, err error) {
	e.buf.Write(p)
	for e.buf.Len() > blockSize {
		out, err := aes.AesCBCEncrypt(e.buf.Next(blockSize), e.a1.aes.Key, e.a1.aes.Iv)
		if err != nil {
			return 0, err
		}
		_, err = e.w.Write(out)
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
func (e *encrypt) Close() error {
	e.a1.lastLength = e.buf.Len()
	if e.a1.lastLength != 0 {
		e.buf.Write(utils.Random(blockSize - e.buf.Len()))
		out, err := aes.AesCBCEncrypt(e.buf.Bytes(), e.a1.aes.Key, e.a1.aes.Iv)
		if err != nil {
			return err
		}
		e.w.Write(out)
	}
	return nil
}

func (a *A1) newEncrypt(out io.Writer) (*encrypt, error) {
	k, err := json.Marshal(a.aes)
	if err != nil {
		return nil, err
	}
	o, err := rsa.Encrypt(k, a.rsaKey.Key)
	if err != nil {
		return nil, err
	}
	_, err = out.Write(o)
	if err != nil {
		return nil, err
	}
	return &encrypt{
		w:  out,
		a1: a,
	}, nil
}
