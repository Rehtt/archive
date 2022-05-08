/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:11
 */

package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

//var pubKey []byte
//var pirKey []byte
//var open bool

const size = 4096
const plaintextSize = 501

//func InitKey(privateKey, publicKey []byte) {
//	pubKey = publicKey
//	pirKey = privateKey
//	open = true
//}
//
//type Reader struct {
//	r   io.Reader
//	buf bytes.Buffer
//}
//
//func (r *Reader) Read(p []byte) (n int, err error) {
//	// 按需解密
//	s := len(p)
//	if r.buf.Len() < s {
//		// rsa密文位数等于公钥位数
//		buf := make([]byte, size/8)
//
//		for r.buf.Len() < s {
//			nn, err := io.ReadFull(r.r, buf)
//			if err != nil {
//				return 0, err
//			}
//			if nn == 0 {
//				break
//			} else if nn != size/8 {
//				return 0, errors.New("file error")
//			}
//			out, err := Decrypt(buf)
//			if err != nil {
//				return 0, err
//			}
//			r.buf.Write(out)
//		}
//		if r.buf.Len() < s {
//			s = r.buf.Len()
//		}
//	}
//
//	n = copy(p, r.buf.Next(s))
//	return
//}
//
//type Writer struct {
//	w   io.Writer
//	buf bytes.Buffer
//}
//
//func (w *Writer) Write(p []byte) (n int, err error) {
//	if !open {
//		return w.w.Write(p)
//	}
//
//	w.buf.Write(p)
//	for w.buf.Len() >= plaintextSize {
//		out, err := Encrypt(w.buf.Next(plaintextSize))
//		if err != nil {
//			return 0, err
//		}
//		_, err = w.w.Write(out)
//		if err != nil {
//			return 0, err
//		}
//	}
//	return len(p), nil
//}
//func (w *Writer) Close() error {
//	if w.buf.Len() != 0 {
//		out, err := Encrypt(w.buf.Bytes())
//		if err != nil {
//			return err
//		}
//		_, err = w.w.Write(out)
//		if err != nil {
//			return err
//		}
//		w.buf.Reset()
//	}
//	return nil
//}
//
//func NewReader(r io.Reader) *Reader {
//	return &Reader{r: r}
//}
//func NewWriter(w io.Writer) *Writer {
//	return &Writer{w: w}
//}

func Encrypt(src, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil || !strings.Contains(strings.ToUpper(block.Type), "PUBLIC KEY") {
		return nil, errors.New("未指定公钥，或公钥错误")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	result, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), src)
	return result, err
}
func Decrypt(cip, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil || !strings.Contains(strings.ToUpper(block.Type), "PRIVATE KEY") {
		return nil, errors.New("未指定私钥，或私钥错误")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, cip)
	return plainText, err
}
func GenerateRsaKey() (priKey, pubKey []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, nil, err
	}
	derText := x509.MarshalPKCS1PrivateKey(privateKey)
	block := pem.Block{
		Type:  "rsa private key",
		Bytes: derText,
	}
	priKey = pem.EncodeToMemory(&block)

	publicKey := privateKey.PublicKey
	derstream, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, nil, err
	}

	block = pem.Block{
		Type:  "rsa public key",
		Bytes: derstream,
	}
	pubKey = pem.EncodeToMemory(&block)
	return
}
