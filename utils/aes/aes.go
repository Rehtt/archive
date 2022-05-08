/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/8 下午 02:24
 */

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesCBCEncrypt(rawData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	//rawData = PKCS7Padding(rawData, blockSize)
	cipherText := make([]byte, blockSize+len(rawData))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)
	return cipherText, nil
}

func AesCBCDncrypt(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	encryptData = encryptData[blockSize:]
	if len(encryptData)%blockSize != 0 {
		fmt.Println(len(encryptData))
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	//encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

func Encrypt(rawData, key, iv []byte) (out []byte, err error) {
	out, err = AesCBCEncrypt(rawData, key, iv)
	if err != nil {
		return
	}
	return
}

func Dncrypt(data, key, iv []byte) ([]byte, error) {
	dnData, err := AesCBCDncrypt(data, key, iv)
	if err != nil {
		return nil, err
	}
	return dnData, nil
}
