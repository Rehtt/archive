/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/10 下午 10:16
 */

package model

import (
	"bufio"
	"encoding/binary"
	"io"
)

var HeadVersionEndian = binary.BigEndian

func ParseVersion(i *bufio.Reader) (HeadVersion, error) {
	var out HeadVersion
	err := binary.Read(i, HeadVersionEndian, &out)
	return out, err
}

func WriteVersion(i io.Writer, version HeadVersion) error {
	return binary.Write(i, HeadVersionEndian, version)
}
