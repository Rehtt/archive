/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/5/7 14:10
 */

package main

import (
	"archive/utils"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	gen = flag.Bool("generateRsaKey", false, "生成公钥和私钥")

	archiveMode = flag.Bool("a", false, "压缩")
	encrypt     = flag.Bool("e", false, "开启加密")

	unArchiveMode = flag.Bool("ua", false, "解压")
	decrypt       = flag.Bool("d", false, "开启解密")
	check         = flag.Bool("check", false, "检测压缩包的错误")

	inFile  = flag.String("in", "", "输入文件")
	outFile = flag.String("out", "", "输出文件")

	keyFile = flag.String("inKey", "", "解密指定私钥，加密指定公钥")
)

func main() {
	flag.Parse()

	if *gen {
		priKey, pubKey, err := utils.GenerateRsaKey()
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
	if *inFile == "" {
		fmt.Println("输入为空")
		flag.Usage()
		return
	}

	if *check {
		if *decrypt {
			data, err := ioutil.ReadFile(*keyFile)
			if err != nil {
				log.Fatalln(err)
			}
			utils.InitKey(data, nil)
		}
		err := utils.CheckPackage(*inFile)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("没有错误")
		return
	}

	if *outFile == "" {
		fmt.Println("输出为空")
		flag.Usage()
		return
	}

	if *archiveMode {
		if *encrypt {
			data, err := ioutil.ReadFile(*keyFile)
			if err != nil {
				log.Fatalln(err)
			}
			utils.InitKey(nil, data)
		}
		err := utils.Compress(*inFile, *outFile)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if *unArchiveMode {

		if *decrypt {
			data, err := ioutil.ReadFile(*keyFile)
			if err != nil {
				log.Fatalln(err)
			}
			utils.InitKey(data, nil)
		}
		err := utils.DeTar(*inFile, *outFile)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
	flag.Usage()
}
