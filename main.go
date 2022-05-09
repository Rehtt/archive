/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:10
 */

package main

import (
	"flag"
	"fmt"
	"github.com/Rehtt/archive/utils"
	"github.com/Rehtt/archive/utils/rsa"
	"io/ioutil"
	"log"
)

var (
	gen = flag.Bool("generateRsaKey", false, "生成公钥和私钥")

	archiveMode = flag.Bool("a", false, "压缩")
	encrypt     = flag.Bool("e", false, "开启加密")

	unArchiveMode = flag.Bool("ua", false, "解压，会自动检查是否加密")
	check         = flag.Bool("check", false, "检测压缩包的错误")

	inFile  = flag.String("in", "", "输入文件")
	outFile = flag.String("out", "", "输出文件")

	keyFile = flag.String("inKey", "", "解密指定私钥，加密指定公钥")
)

func main() {
	flag.Parse()

	// 生成证书
	if *gen {
		priKey, pubKey, err := rsa.GenerateRsaKey()
		if err != nil {
			log.Fatalln(err)
		}
		err = ioutil.WriteFile("private.pem", priKey, 644)
		if err != nil {
			log.Fatalln(err)
		}
		err = ioutil.WriteFile("public.pem", pubKey, 644)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	// 加密
	if *encrypt {
		data, err := ioutil.ReadFile(*keyFile)
		if err != nil {
			log.Fatalln(err)
			flag.Usage()
			return
		}
		utils.InitEncrypt(nil, data)
	}

	// 输入位置
	if *inFile == "" {
		fmt.Println("输入为空")
		flag.Usage()
		return
	}

	// 检查错误
	if *check {
		data, err := ioutil.ReadFile(*keyFile)
		if err == nil {
			utils.InitEncrypt(data, nil)
		}
		err = utils.CheckPackage(*inFile)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("没有错误")
		return
	}

	// 输出位置
	if *outFile == "" {
		fmt.Println("输出为空")
		flag.Usage()
		return
	}

	// 压缩
	if *archiveMode {
		err := utils.Compress(*inFile, *outFile, *encrypt)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	// 解压
	if *unArchiveMode {

		data, err := ioutil.ReadFile(*keyFile)
		if err == nil {
			utils.InitEncrypt(data, nil)
		}
		err = utils.Uncompress(*inFile, *outFile)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
	flag.Usage()
}
