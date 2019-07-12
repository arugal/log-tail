# log tail

[![Build Status](https://travis-ci.org/arugal/log-tail.svg?branch=master)](https://travis-ci.org/arugal/log-tail)
[![Coverage](https://codecov.io/gh/arugal/log-tail/branch/master/graph/badge.svg)](https://codecov.io/gh/arugal/log-tail)


log-tail 是一个支持浏览器实时查看日志的小工具


## 目录

* [开发状态](#开发状态)
* [快速启动](#快速启动)
    * [直接启动](#直接启动)
    * [部署为系统服务](#部署为系统服务)
* [配置文件说明](#配置文件说明)
* [开发计划](#开发计划)


## 开发状态
log-tail 仍然处于开发阶段，未经充分测试与验证，不推荐用于生产环境。

## 快速启动

根据对应的操作系统以及架构，从[Release](https://github.com/arugal/log-tail/releases)页面下载最新的版本的程序

### 直接启动
  ./log-tail -c log_tail.ini
### 部署为系统服务
   ...
## 配置文件说明
```
[common] # 服务配置
bind_addr = 127.0.0.1 # 服务绑定地址
bind_port = 3000 # 服务监听端口
# minute
conn_max_time = 10 # 单个日志连接连接最长时间，单位:分钟
# second
heart_interval = 10 # 日志查看连接心跳间隔，单位:秒
# console or real logFile path like ./log_tail.log
log_file = ./log_tail.log # 日志输入目录
# trace, debug, warn, error
log_level = info # 日志登记
log_max_days = 3 # 日志保存时间
# set web address for control log-tail action by http api such as reload
user = admin # 登陆页面账号，设置为空字符串既为无账号密码
pwd = admin # 登陆页面密码
# ignore file, the scope is global
ignore_suffix = .jar,.war,.html,.js,.css,.java,.class,.gz,.tar,.zip,.rar,.jpg,.png,.xls,.xlxs,.pdf # 文件后缀过滤
ignore_regexp = # 文件正则过滤
# start reading position [size - offset:size]
last_read_offset = 1000 # 日志开始读取位置
assets_dir = "" # 指定前端代码地址

[catalog1] # 浏览目录配置
path = /var/application/logs # 日志目录路径
# ignore file, the scope is this
ignore_suffix = .txt # 文件后缀过滤
ignore_regexp = # 文件正则过滤
```

## 开发计划

- [x] travis、Codecov集成
- [x] 项目包路径切换
- [ ] web自适应优化
- [x] 项目文档编写

- [ ] 配置文件替换成Yaml文件
- [ ] 热加载配置文件
- [ ] 微内核+可拓展模式的探索