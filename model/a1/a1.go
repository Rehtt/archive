/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/10 下午 09:49
 */

package a1

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rehtt/archive/model"
	"github.com/Rehtt/archive/utils"
	"github.com/Rehtt/archive/utils/rsa"
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

	aes aesConfig
}
type aesConfig struct {
	Iv   []byte `json:"iv"`
	Key  []byte `json:"key"`
	Salt []byte `json:"salt"`
}

func (a *A1) Compress(in string, out *os.File, encrypt *model.Encrypt) error {
	var head bytes.Buffer
	defer func() {
		out.Seek(0, 0)
		buf := make([]byte, 3)
		s := blockSize - a.lastLength
		if s < 256 {
			buf[0] = byte(s)
		} else {
			buf[0] = byte(255)
			buf[1] = byte(s) - buf[0]
			if int(buf[1])-255 > 0 {
				buf[2] = buf[1] - 255
			}
		}
		head.Write(buf)
		head.Write(bytes.Repeat([]byte{1}, headSize-head.Len()))
		out.Write(head.Bytes())
	}()
	head.Write([]byte{'A', 1})

	var w io.Writer
	if encrypt.Key != nil {
		head.Write([]byte{1})

		a.aes.Salt = utils.Random(16)
		a.aes.Iv = utils.Random(16)
		a.aes.Key = utils.Random(32)
		ew, err := a.newEncrypt(w)
		if err != nil {
			return err
		}
		defer ew.Close()
	} else {
		head.Write([]byte{0})
		w = out
	}
	x, err := xz.NewWriter(w)
	if err != nil {
		return err
	}
	defer x.Close()

	err = utils.Tar(in, out)
	if err != nil {
		return err
	}
	return nil
}

// todo 解压bug修复
func (a *A1) Uncompress(in *bufio.Reader, out string, encrypt *model.Encrypt) error {
	err := a.uncompressParse(in, encrypt)
	if err != nil {
		return err
	}
	var r io.Reader
	if a.isEncrypt {
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

func (a *A1) uncompressParse(in *bufio.Reader, encrypt *model.Encrypt) error {
	head := make([]byte, headSize)
	_, err := io.ReadFull(in, head)
	if err != nil {
		return errors.New("file not recognized")
	}

	a.lastLength = int(head[3]) + int(head[4]) + int(head[5])
	a.isEncrypt = head[2] == 1
	fmt.Println(a.lastLength, a.isEncrypt)
	if !a.isEncrypt {
		return err
	}

	blockData := make([]byte, blockSize)
	_, err = io.ReadFull(in, blockData)
	if err != nil {
		return errors.New("file error")
	}
	out, err := rsa.Decrypt(blockData, encrypt.Key)
	if err != nil {
		return err
	}
	err = json.Unmarshal(out, &a.aes)
	return err
}
