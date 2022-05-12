/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/10 下午 09:49
 */

package model

import (
	"bufio"
	"errors"
	"github.com/Rehtt/archive/utils"
	"os"
)

type Model interface {
	Version() string
	Compress(in string, out *os.File, encrypt *Encrypt) error
	Uncompress(in *bufio.Reader, out string, encrypt *Encrypt) error
	CheckPackage(in *bufio.Reader, encrypt *Encrypt) error
}

type Encrypt struct {
	Key []byte
}

var (
	model = make(map[string]Model)
)

func Register(m Model) {
	model[m.Version()] = m
}

func Compress(inFile, outFile, compressModel string, encrypt *Encrypt) error {
	oFile, err := utils.CreateFile(outFile)
	if err != nil {
		return err
	}
	defer oFile.Close()

	if m, ok := model[compressModel]; ok {
		err = m.Compress(inFile, oFile, encrypt)
		return err
	}
	return errors.New("未知编码")
}

func Uncompress(inFile, outFile string, encrypt *Encrypt) error {
	iFile, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer iFile.Close()

	buf := bufio.NewReader(iFile)
	version, err := ParseVersion(buf)
	if err != nil {
		return err
	}
	if m, ok := model[version]; ok {
		return m.Uncompress(buf, outFile, encrypt)
	}
	return errors.New("未知编码")
}

func CheckPackage(inFile string, encrypt *Encrypt) error {
	iFile, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer iFile.Close()

	buf := bufio.NewReader(iFile)
	version, err := ParseVersion(buf)
	if err != nil {
		return err
	}
	if m, ok := model[version]; ok {
		return m.CheckPackage(buf, encrypt)
	}
	return errors.New("未知编码")

}
