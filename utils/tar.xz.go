/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:52
 */

package utils

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var ShowProcessFilePath bool

func Compress(file *os.File, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if ShowProcessFilePath {
		fmt.Println(filepath.Join(prefix, info.Name()))
	}

	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			// 跳过符号链接，避免循环链接导致的过深解析
			if fi.Mode()&os.ModeSymlink != 0 {
				continue
			}
			childPath := filepath.Join(file.Name(), fi.Name())
			f, err := os.Open(childPath)
			if err != nil {
				return err
			}
			err = Compress(f, prefix, tw)
			f.Close()
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = prefix + "/" + header.Name
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

func UnTar(r io.Reader, to string) error {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		if to != "" {
			file, err := CreateFile(to + "/" + hdr.Name)
			if err != nil {
				return err
			}
			if ShowProcessFilePath {
				fmt.Println(to + "/" + hdr.Name)
			}
			io.Copy(file, tr)
		}
	}
	return nil
}

func CreateFile(name string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(name), 0o755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
