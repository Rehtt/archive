/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:52
 */

package utils

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

//func Tar(in string, out io.Writer) error {
//	inFile, err := os.Open(in)
//	if err != nil {
//		return err
//	}
//	tw := tar.NewWriter(out)
//	defer fmt.Println(tw.Close())
//	err = compress(inFile, "", tw)
//	return err
//}
func Compress(file *os.File, prefix string, tw *tar.Writer) error {
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
			f, err := os.Open(filepath.Join(file.Name(), fi.Name()))
			if err != nil {
				return err
			}
			err = Compress(f, prefix, tw)
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
			filename := filepath.Join(to, hdr.Name)
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
