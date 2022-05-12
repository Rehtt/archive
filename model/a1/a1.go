/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/10 下午 09:49
 */

package a1

import (
	"archive/tar"
	"bufio"
	"errors"
	"github.com/Rehtt/archive/model"
	"github.com/Rehtt/archive/utils"
	"github.com/ulikunitz/xz"
	"io"
	"os"
)

const (
	headSize  = 16
	blockSize = 512
)

type A1 struct {
	lastLength int
	isEncrypt  bool
	rsaKey     *model.Encrypt

	aes aesConfig
}
type aesConfig struct {
	Iv   []byte `json:"iv"`
	Key  []byte `json:"key"`
	Salt []byte `json:"salt"`
}

func (a *A1) Compress(in string, out *os.File, encrypt *model.Encrypt) error {
	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	head := make([]byte, headSize)
	defer func() {
		out.Seek(0, 0)
		s := a.lastLength
		if s < 256 {
			head[3] = byte(s)
		} else {
			head[3] = byte(255)
			head[4] = byte(s) - head[3]
			if int(head[4])-255 > 0 {
				head[5] = head[4] - 255
			}
		}
		out.Write(head)
	}()
	head[0] = 'A'
	head[1] = 1

	out.Write(head)

	var w io.Writer
	if encrypt.Key != nil {
		head[2] = 1

		a.rsaKey = encrypt
		a.aes.Salt = utils.Random(16)
		a.aes.Iv = utils.Random(16)
		a.aes.Key = utils.Random(32)
		ew, err := a.newEncrypt(out)
		if err != nil {
			return err
		}
		defer ew.Close()
		w = ew
	} else {
		head[2] = 0
		w = out
	}
	x, err := xz.NewWriter(w)
	if err != nil {
		return err
	}
	defer x.Close()

	tw := tar.NewWriter(x)
	defer tw.Close()
	err = utils.Compress(inFile, "", tw)
	return err
}

func (a *A1) Uncompress(in *bufio.Reader, out string, encrypt *model.Encrypt) error {
	head := make([]byte, headSize)
	_, err := io.ReadFull(in, head)
	if err != nil {
		return errors.New("file not recognized")
	}

	a.lastLength = int(head[3]) + int(head[4]) + int(head[5])
	a.isEncrypt = head[2] == 1

	var r io.Reader
	if a.isEncrypt {
		a.rsaKey = encrypt
		r, err = a.newDecrypt(in)
		if err != nil {
			return err
		}
	} else {
		r = in
	}
	x, err := xz.NewReader(r)
	if err != nil {
		return err
	}
	return utils.UnTar(x, out)
}

func (a *A1) CheckPackage(in *bufio.Reader, encrypt *model.Encrypt) error {
	return a.Uncompress(in, "", encrypt)
}

func (A1) Version() string {
	return "A1"
}

var a1 = new(A1)

func init() {
	model.Register(a1)
}
