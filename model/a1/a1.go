/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/10 下午 09:49
 */

package a1

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/Rehtt/archive/model"
	"github.com/Rehtt/archive/utils"
	"github.com/ulikunitz/xz"
	"io"
	"os"
)

const (
	headSize  = 14
	blockSize = 512
)

type A1 struct {
	isEncrypt bool
	rsaKey    *model.Encrypt

	aes aesConfig
}
type aesConfig struct {
	Iv   []byte `json:"iv"`
	Key  []byte `json:"key"`
	Salt []byte `json:"salt"`
}

type Head struct {
	IsEncryption bool
}

func (a *A1) Compress(in string, out *os.File, encrypt *model.Encrypt) error {
	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	head := Head{
		IsEncryption: encrypt.Key != nil,
	}
	size := headSize - binary.Size(head)
	if size < 0 {
		return errors.New("head too big")
	}
	binary.Write(out, model.HeadVersionEndian, head)
	if size > 0 {
		out.Write(make([]byte, size))
	}

	//defer func() {
	//	out.Seek(int64(a.Version().Len()), 0)
	//	var tmp bytes.Buffer
	//	binary.Write(&tmp, model.HeadVersionEndian, head)
	//	out.Write(tmp.Bytes())
	//}()

	var w io.Writer
	if encrypt.Key != nil {
		head.IsEncryption = true
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
	readHead := make([]byte, headSize)
	_, err := io.ReadFull(in, readHead)
	if err != nil {
		return errors.New("file not recognized")
	}
	var head Head
	if err = binary.Read(bytes.NewReader(readHead), model.HeadVersionEndian, &head); err != nil {
		return err
	}

	a.isEncrypt = head.IsEncryption

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

func (A1) Version() model.HeadVersion {
	return model.HeadVersion{
		Protocol: 'A',
		Version:  1,
	}
}

var a1 = new(A1)

func init() {
	model.Register(a1)
}
