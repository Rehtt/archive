使用RSA 4096 bit 非对称加密、解密的tar.xz压缩、解压缩的简单实现

# 命令

## 生成公钥与私钥
```shell
./archive -generateRsaKey
```
将生成`public.pem（加密）`与`private.pem（解密）`证书

> `private.pem`必须妥善保管。若有遗失，则加密的文件就再也无法打开

## 普通tar.xz打包压缩
```shell
./archive -a -in test -out test.tar.xz
```

## 加密压缩
```shell
./archive -a -e -inKey public.pem -in test -out test.tar.xz
```

## 普通tar.xz解压
```shell
./archive -ua -in test.tar.xz -out test
```

## 解密解压
```shell
./archive -ua -inKey public.pem -in test.tar.xz -out test
```

# 命令树
```shell
  -a    压缩
  -check
        检测压缩包的错误
  -d    开启解密
  -e    开启加密
  -generateRsaKey
        生成公钥和私钥
  -in string
        输入文件
  -inKey string
        解密指定私钥，加密指定公钥
  -out string
        输出文件
  -ua
        解压

```
