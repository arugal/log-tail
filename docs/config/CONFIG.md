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