/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:52
 */

package utils

import (
	"archive/tar"
	"fmt"
	"github.com/ulikunitz/xz"
	"io"
	"os"
	"strings"
)

func Compress(targetPath string, dest string) error {
	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()

	file, err := os.Open(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	r := NewWriter(d)
	defer func() {
		if err := r.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	xw, err := xz.NewWriter(r)
	if err != nil {
		return err
	}
	defer xw.Close()
	tw := tar.NewWriter(xw)
	defer tw.Close()
	return compress(file, "", tw)
	//for _, file := range files {
	//	err := compress(file, "", tw)
	//	if err != nil {
	//		return err
	//	}
	//}
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

func DeTar(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	r := NewReader(srcFile)
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
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
