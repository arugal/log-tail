log tail in web

[![Build Status](https://travis-ci.org/Arugal/log-tail.svg?branch=master)](https://travis-ci.org/Arugal/log-tail)



# go 安装

+ go version >= 1.12
+ [windowns安装](https://www.runoob.com/go/go-environment.html)
+ [linux安装](https://www.runoob.com/go/go-environment.html)


# 代码编译
+ ```linux``` 如果有 ```make``` 直接执行 ```make all``` 接口,编译后的二进制文件在```bin/```目录，配置文件在```conf/```目录,
 通过 ```-c``` 参数指定配置文件路径 例： ```./log-tail -c conf/log_tail.ini```
+ ```windows``` 目前没有脚本，直接执行 ```go build -o bin/log-tail ./cmd``` 编译
+ 出现 ```cannot find package``` 异常时，将环境变量 ```GO111MODULE``` 设置为 ```on```

# 注意
+ ```go```只有在编译的时候需要```go```的环境,在各个```os```上需要重新编译，这点与```java```不一样
+ 出现下载依赖超时异常添加环境变量```GOPROXY=https://goproxy.io```

- TODO
- [ ] travis、Codecov集成
- [x] 项目内引用路径切换
- [ ] web自适应
- [ ] 项目文档编写

- Next TODO
- [ ] 配置文件替换成Yaml文件
- [ ] 查看目录配置热加载
- [ ] 微内核+可拓展模式的探索