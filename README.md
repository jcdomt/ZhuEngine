# ZhuEngine 流量均衡网关

## 配置手册

### ze.ini 系统主配置文件:

#### zhu_engine 配置节点
项目 | 类型 | 简介
-- | -- | --
port | int | http 监听端口
host | string | 监听的主域名
```ini
[zhu_engine]
port=80
host=test.wzjer.com
```

#### https 配置节点
项目 | 类型 | 简介
-- | -- | --
port | bool | 是否启用 https
crt | string | 公钥路径
key | string | 私钥路径
port | int | https 监听端口
force | bool | 是否强制将 http 转为 https
```ini
[https]
enable=true
crt=./cert/fullchain.pem
key=./cert/privkey.key
port=443
force=true
```

### web.ini 站点配置文件:

```ini
[站点名称]
type=解析类型： pattern / domain
url=解析子域名或路径
server=服务器配置
autorun=是否自动运行
exec=自动运行的脚本
schedule=负载均衡调度器
cgi=CGI名称
```

负载均衡 | server | 简介
-- | -- | --
round | ip1,ip2 | 轮询
random | ip1?weight1,ip2?weight2 | 随机均衡

cgi 和 schedule 不能同时生效
如果同时配置，系统会优先代理 CGI

若配置 CGI，则 server 配置为脚本根目录

### cgi.ini CGI配置文件:
```ini
[CGI名称]
cgi=CGI路径
```