/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/10 下午 09:49
 */

package a1

import (
	"bufio"
	"github.com/Rehtt/archive/model"
	"io"
)

type A1 struct {
	isEncrypt bool
}

func (a *A1) Compress(in *bufio.Reader, out io.Writer, encrypt bool) error {
	//TODO implement me
	panic("implement me")
}

func (a *A1) Uncompress(in *bufio.Reader, out string) error {
	//TODO implement me
	panic("implement me")
}

func (a *A1) CheckPackage(in *bufio.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (a A1) Version() string {
	return "A1"
}

var a1 = new(A1)

func init() {
	model.Register(a1)
}

func parse(in *bufio.Reader) {

}
