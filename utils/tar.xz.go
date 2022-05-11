/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:52
 */

package utils

import (
	"archive/tar"
	"errors"
	"github.com/ulikunitz/xz"
	"io"
	"os"
	"path/filepath"
)

var version = []byte{'A', 1}
var filling = 16 - 3

func Compress(targetPath, dest string, isEncrypt bool) error {
	outFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outFile.Close()

	file, err := os.Open(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 编码版本
	outFile.Write(version)

	var w io.Writer
	if isEncrypt {
		outFile.Write([]byte{1})
		outFile.Write(make([]byte, filling))

		r, err := NewWriter(outFile)
		if err != nil {
			return err
		}
		defer r.Close()
		w = r
	} else {
		outFile.Write([]byte{0})
		outFile.Write(make([]byte, filling))

		w = outFile
	}

	xw, err := xz.NewWriter(w)
	if err != nil {
		return err
	}
	defer xw.Close()
	tw := tar.NewWriter(xw)
	defer tw.Close()
	return compress(file, "", tw)
}
func compress(file *os.File, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func Uncompress(tarFile, dest string) error {
	inFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	head := make([]byte, 16)
	inFile.Read(head)
	if head[0] != version[0] || head[1] != version[1] {
		return errors.New("file version error")
	}
	isEn := head[2] == 1

	var r io.Reader
	if isEn {
		rr, err := NewReader(inFile)
		if err != nil {
			return err
		}
		r = rr
	} else {
		r = inFile
	}

	xr, err := xz.NewReader(r)
	if err != nil {
		return err
	}
	tr := tar.NewReader(xr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		if dest != "" {
			filename := dest + hdr.Name
			file, err := CreateFile(filename)
			if err != nil {
				return err
			}
			io.Copy(file, tr)
		}

	}
	return nil
}

func CreateFile(name string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(name), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func CheckPackage(tarFile string) error {
	return Uncompress(tarFile, "")
}
