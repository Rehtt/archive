/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:11
 */

package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"strings"
)

var pubKey []byte
var pirKey []byte
var open bool

func InitKey(privateKey, publicKey []byte) {
	pubKey = publicKey
	pirKey = privateKey
	open = true
}

var separate = []byte("Ras")

type Reader struct {
	r   io.Reader
	buf bytes.Buffer
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if !open {
		return r.r.Read(p)
	}
	if r.buf.Len() == 0 {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r.r)
		if err != nil {
			return 0, err
		}
		for _, v := range bytes.Split(buf.Bytes(), separate) {
			if len(v) == 0 {
				continue
			}
			out, err := Decrypt(v)
			if err != nil {
				return 0, err
			}
			r.buf.Write(out)
		}
	}
	n = copy(p, r.buf.Next(len(p)))
	return
}

type Writer struct {
	w   io.Writer
	buf bytes.Buffer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if !open {
		return w.w.Write(p)
	}
	w.buf.Write(p)
	return len(p), nil
}
func (w *Writer) Close() error {
	for w.buf.Len() != 0 {
		var wBuf bytes.Buffer
		wBuf.Write(separate)
		out, err := Encrypt(w.buf.Next(501))
		if err != nil {
			return err
		}
		wBuf.Write(out)
		_, err = io.Copy(w.w, &wBuf)
		if err != nil {
			return err
		}
	}
	w.buf.Reset()
	return nil
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func Encrypt(src []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil || !strings.Contains(strings.ToUpper(block.Type), "PUBLIC KEY") {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	result, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), src)
	return result, err
}
func Decrypt(cip []byte) ([]byte, error) {
	block, _ := pem.Decode(pirKey)
	if block == nil || !strings.Contains(strings.ToUpper(block.Type), "PRIVATE KEY") {
		return nil, errors.New("failed to decode PEM block containing PRIVATE key")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, cip)
	return plainText, err
}
func GenerateRsaKey() (priKey, pubKey []byte, err error) {

	// 1. 使用rsa中的GenerateKey方法生成私钥

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	// 2. 通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串

	derText := x509.MarshalPKCS1PrivateKey(privateKey)

	// 3. 要组织一个pem.Block(base64编码)

	block := pem.Block{
		Type:  "rsa private key", // 这个地方写个字符串就行
		Bytes: derText,
	}

	// 4. pem编码
	priKey = pem.EncodeToMemory(&block)

	// ============ 公钥 ==========

	// 1. 从私钥中取出公钥

	publicKey := privateKey.PublicKey

	// 2. 使用x509标准序列化

	derstream, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, nil, err
	}

	// 3. 将得到的数据放到pem.Block中

	block = pem.Block{
		Type:  "rsa public key",
		Bytes: derstream,
	}

	// 4. pem编码
	pubKey = pem.EncodeToMemory(&block)
	return
}
