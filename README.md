# log tail

[![Build Status](https://travis-ci.org/arugal/log-tail.svg?branch=master)](https://travis-ci.org/arugal/log-tail)
[![Coverage](https://codecov.io/gh/arugal/log-tail/branch/master/graph/badge.svg)](https://codecov.io/gh/arugal/log-tail)


log-tail 是一个支持浏览器实时查看日志的小工具


## 目录

* [开发状态](#开发状态)
* [快速启动](#快速启动)
    * [直接启动](#直接启动)
    * [部署为系统服务](#部署为系统服务)
* [命令说明](#命令说明)
* [配置文件说明](#配置文件说明)
* [开发计划](#开发计划)


## 开发状态
log-tail 仍然处于开发阶段，未经充分测试与验证，不推荐用于生产环境。

## 快速启动

根据对应的操作系统以及架构，从[Release](https://github.com/arugal/log-tail/releases)页面下载最新的版本的程序

### 直接启动
  ./log-tail -c ./config.yaml
### 部署为系统服务
   ...
## 命令说明
```
  -c, --config string     config file of log-tail (default "./config.yaml")
  -h, --help              help for log-tail
  -H, --host string       host of log-tail (default "-")
  -l, --loglevel string   log level of log-tail (default "-")
  -p, --port int          port of log-tail (default -1)
  -v, --version           version of log-tail

```
## 配置文件说明
```
server:
  host: 127.0.0.1 # 服务绑定地址
  port: 3000 # 服务监听端口
  secure:
    user: admin # 登陆页面账号，设置为空字符串即为无账号密码
    pwd: admin # 登陆页面密码
common:
  last_read_offset: 1000 # 日志开始读取位置
  conn_max_time: 10 # 单个日志查看连接最长时间（到时自动关闭连接），单位:分钟
  heart_interval: 10 # 日志查看窗口连接心跳间隔（超时自动关闭连接），单位:秒
  log:
    file: console # 日志输出目录，设置为console时输出在控制台
    level: info # 日志等级：trace, debug, warn, error
    max_days: 1 # 日志保存时间
  ignore: # ignore file, the scope is global
    suffix: # 文件后缀过滤
      - .jar
      - .war
    regexp: # 文件正则过滤
      - "*.log.*"

catalogs: # 浏览目录配置，可配置多个目录
  - name: app1
    path: /tmp/app1/logs # 日志目录路径
    ignore:
      suffix: # 文件后缀过滤
        - .txt
      regexp: # 文件正则过滤
        - "*.out.*"
```

## 开发计划

- [x] travis、Codecov集成
- [x] 项目包路径切换
- [ ] web自适应优化
- [x] 项目文档编写

- [x] 配置文件替换成Yaml文件
- [ ] 热加载配置文件
- [ ] 微内核+可拓展模式的探索
