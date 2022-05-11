/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/10 下午 09:49
 */

package model

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Model interface {
	Version() string
	Compress(in *bufio.Reader, out io.Writer, encrypt bool) error
	Uncompress(in *bufio.Reader, out string) error
	CheckPackage(in *bufio.Reader) error
}

var (
	model = make(map[string]Model)
)

func Register(m Model) {
	model[m.Version()] = m
}

func Compress(inFile, outFile, compressModel string, encrypt bool) error {
	iFile, err := os.Open(filepath.Clean(inFile))
	if err != nil {
		return err
	}
	defer iFile.Close()

	oFile, err := os.Open(filepath.Clean(outFile))
	if err != nil {
		return err
	}
	defer oFile.Close()

	buf := bufio.NewReader(iFile)

	if m, ok := model[compressModel]; ok {
		return m.Compress(buf, oFile, encrypt)
	}
	return errors.New("未知编码")
}

func Uncompress(inFile, outFile string) error {
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
		return m.Uncompress(buf, outFile)
	}
	return errors.New("未知编码")
}

func CheckPackage(inFile string) error {
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
		return m.CheckPackage(buf)
	}
	return errors.New("未知编码")

}
