/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/8 下午 02:25
 */

package utils

import (
	"bytes"
	"crypto/rand"
	"log"
	"math/big"
)

func Random(l int) []byte {
	var buf bytes.Buffer
	for buf.Len() < l {
		n, err := rand.Int(rand.Reader, big.NewInt(256))
		if err != nil {
			log.Fatalln(err)
		}
		buf.WriteByte(byte(n.Int64()))
	}
	return buf.Next(l)
}
