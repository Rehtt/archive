/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:52
 */

package utils

import (
	"archive/tar"
	"archive/utils/rsa"
	"bufio"
	"github.com/ulikunitz/xz"
	"io"
	"log"
	"os"
	"strings"
)

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

	var w io.Writer
	if isEncrypt {
		outFile.Write([]byte("Rsa!"))
		r := rsa.NewWriter(outFile)
		defer func() {
			if err := r.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
		w = r
	} else {
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

	buf := bufio.NewReader(inFile)
	var r io.Reader
	// 检查头部
	temp, err := buf.Peek(4)
	if err != nil {
		return err
	}
	if string(temp) == "Rsa!" {
		buf.Discard(4)
		r = rsa.NewReader(buf)
	} else {
		r = buf
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
			file, err := createFile(filename)
			if err != nil {
				return err
			}
			io.Copy(file, tr)
		}

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

func CheckPackage(tarFile string) error {
	return Uncompress(tarFile, "")
}
