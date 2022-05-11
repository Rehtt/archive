/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/11 21:53
 */

package a1

import (
	"bufio"
	"bytes"
	"github.com/Rehtt/archive/utils/aes"
	"io"
)

type decrypt struct {
	r   *bufio.Reader
	buf bytes.Buffer

	a1 *A1
}

func (d *decrypt) Read(p []byte) (n int, err error) {
	size := len(p)
	if d.buf.Len() < size {
		buf := make([]byte, blockSize)
		for d.buf.Len() < size {
			n, err = io.ReadFull(d.r, buf)
			if err != nil {
				return 0, err
			}
			out, err := aes.AesCBCDncrypt(buf, d.a1.aes.Key, d.a1.aes.Iv)
			if err != nil {
				return 0, err
			}
			_, err = d.r.Peek(1)
			if err != nil {
				n = d.a1.lastLength
			}
			d.buf.Write(out[:n])
		}
	}
	n = copy(p, d.buf.Next(size))
	return
}

func (a *A1) newDecrypt(r *bufio.Reader) (*decrypt, error) {
	return &decrypt{
		r:  r,
		a1: a,
	}, nil
}
