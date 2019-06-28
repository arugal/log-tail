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