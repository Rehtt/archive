/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/10 下午 10:16
 */

package model

import (
	"bufio"
	"fmt"
)

func ParseVersion(i *bufio.Reader) (string, error) {
	temp, err := i.Peek(2)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%d", string(temp[0]), temp[1]), nil
}
