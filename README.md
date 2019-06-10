log tail in web

[![Build Status](https://travis-ci.org/Arugal/log-tail.svg?branch=master)](https://travis-ci.org/Arugal/log-tail)



# go 安装

+ go version >= 1.12
+ [windowns安装](https://www.runoob.com/go/go-environment.html)
+ [linux安装](https://www.runoob.com/go/go-environment.html)


# 代码编译
+ linux 如果有 make 直接执行 make all 接口,编译后的二进制文件在bin/目录，配置文件在conf/目录, 通过 -c 参数指定配置文件路径 例： ./log-tail -c conf/log_tail.ini
+ windows 目前没有脚本，直接执行 go build -o bin/log-tail ./cmd 编译

# 注意
+ go只有在编译的时候需要go的环境,在各个os上需要重新编译，这点与java不一样