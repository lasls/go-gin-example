# Go-Gin-Example

这是一个基于Go语言和Gin框架的Web应用程序示例。

## 功能特性

- DNS管理：支持DNSPod、阿里云DNS、火山引擎DNS
- 域名验证：支持DNS验证、CNAME验证、文件验证
- 数据库存储：支持将DNS记录存储到本地数据库
- 批量操作：支持批量创建、更新、删除DNS记录

## 安装

```bash
go mod tidy
go build -o server .
```

## 配置

在 `conf/app.ini` 文件中配置相关参数：

```ini
[server]
# 服务端口
PORT = 8080

[database]
TYPE = mongodb
CONNECTION_STRING = mongodb://localhost:2707
DATABASE = blog
USERNAME = 
PASSWORD = 

[app]
# 默认域名
DOMAIN_NAME = example.com

# DNSPod Token
DNSPOD_TOKEN = your_dnspod_token

# 阿里云AccessKey
ALIYUN_ACCESS_KEY_ID = your_aliyun_access_key_id
ALIYUN_ACCESS_KEY_SECRET = your_aliyun_access_key_secret
ALIYUN_REGION_ID = cn-hangzhou

# 火山引擎AccessKey
VOLCENGINE_ACCESS_KEY_ID = your_volcengine_access_key_id
VOLCENGINE_ACCESS_KEY_SECRET = your_volcengine_access_key_secret
VOLCENGINE_REGION_ID = cn-beijing

[dns_servers]
# DNS服务器列表，支持多个DNS服务器，以逗号分隔
DNS_SERVERS = 8.8.8.8:53,8.8.4.4:53,223.5.5.5:53,1.1.1.1:53,114.114.114.114:53
```

## API接口

### DNS服务器配置

域名验证功能支持多DNS服务器配置，系统会自动轮询配置的DNS服务器进行查询，提高查询的可靠性和可用性。

在 `conf/app.ini` 中配置多个DNS服务器：

```ini
[dns_servers]
DNS_SERVERS = 8.8.8.8:53,8.8.4.4:53,223.5.5.5:53,1.1.1.1:53,114.114.114.114:53
```

系统支持以下常用的公共DNS服务器：
- Google DNS: 8.8.8.8:53, 8.8.4.4:53
- 阿里云DNS: 223.5.5.5:53
- Cloudflare DNS: 1.1.1.1:53
- 114 DNS: 114.114.114.114:53

当某个DNS服务器无法访问时，系统会自动尝试下一个DNS服务器，确保域名验证功能的可靠性。

### 域名验证接口

域名验证功能用于验证您对域名的管理权限，支持三种验证方式：

#### 1. 生成验证挑战

**POST** `/api/v1/dns/verification/challenge`

请求参数：
```json
{
  "domain": "example.com",
  "type": "dns"
}
```

参数说明：
- `domain`: 待验证的域名
- `type`: 验证类型，可选值：`dns`、`cname`、`file`

**DNS验证** (`type: "dns"`)
- 响应示例：
```json
{
  "code": 1,
  "msg": "生成验证挑战成功",
  "data": {
    "success": true,
    "message": "请在您的DNS设置中添加以下TXT记录",
    "challenge": "a1b2c3d4e5f67890",
    "record_type": "TXT",
    "record_name": "_cdnauth.example.com"
  }
}
```
- 用户需在DNS中添加 `_cdnauth.example.com` 的TXT记录，值为 `a1b2c3d4e5f67890`

**CNAME验证** (`type: "cname"`)
- 响应示例：
```json
{
  "code": 1,
  "msg": "生成验证挑战成功",
  "data": {
    "success": true,
    "message": "请在您的DNS设置中添加以下CNAME记录",
    "challenge": "checkpoint-a1b2c3d4.xldns.com",
    "record_type": "CNAME",
    "record_name": "_cdnauth.example.com"
  }
}
```
- 用户需在DNS中添加 `_cdnauth.example.com` 的CNAME记录，值为 `checkpoint-a1b2c3d4.xldns.com`

**文件验证** (`type: "file"`)
- 响应示例：
```json
{
  "code": 1,
  "msg": "生成验证挑战成功",
  "data": {
    "success": true,
    "message": "请将以下内容保存为文件并上传到网站根目录",
    "file_name": "cdn-auth-a1b2c3d4e5f6.txt",
    "file_content": "a1b2c3d4e5f67890"
  }
}
```
- 用户需将内容 `a1b2c3d4e5f67890` 保存为 `cdn-auth-a1b2c3d4e5f6.txt` 文件，并上传到网站根目录下的 `.well-known/` 目录中

#### 2. 执行域名验证

**POST** `/api/v1/dns/verification/verify`

请求参数：
```json
{
  "domain": "example.com",
  "type": "dns",
  "challenge": "a1b2c3d4e5f67890"
}
```

参数说明：
- `domain`: 待验证的域名
- `type`: 验证类型，可选值：`dns`、`cname`、`file`
- `challenge`: 生成验证挑战时返回的challenge值

响应示例：
```json
{
  "code": 1,
  "msg": "DNS验证成功",
  "data": {
    "success": true,
    "message": "DNS验证成功"
  }
}
```

### DNS管理接口

#### 获取域名列表

**GET** `/api/v1/domains?provider=dnspod`

参数：
- `provider`: DNS提供商，可选值：`dnspod`、`aliyun`、`volcengine`

#### 获取DNS记录列表

**GET** `/api/v1/dns/records?provider=dnspod&domain=example.com&sub_domain=www`

参数：
- `provider`: DNS提供商
- `domain`: 域名
- `sub_domain`: 子域名

#### 创建DNS记录

**POST** `/api/v1/dns/records?provider=dnspod&domain_id=12345&sub_domain=www&record_type=A&value=1.1.1.1&record_line=默认`

参数：
- `provider`: DNS提供商
- `domain_id`: 域名ID
- `sub_domain`: 子域名
- `record_type`: 记录类型
- `value`: 记录值
- `record_line`: 线路
- `ttl`: TTL值（可选，默认600）

#### 更新DNS记录

**PUT** `/api/v1/dns/records/{id}?provider=dnspod&domain_id=12345&sub_domain=www&record_type=A&value=1.1.1.1&record_line=默认`

参数：
- `id`: 记录ID
- 其他参数同创建DNS记录

#### 删除DNS记录

**DELETE** `/api/v1/dns/records/{id}?provider=dnspod&domain_id=12345`

参数：
- `id`: 记录ID
- `provider`: DNS提供商
- `domain_id`: 域名ID

#### 设置DNS记录状态

**PUT** `/api/v1/dns/records/{id}/status?provider=dnspod&domain_id=12345&status=enable`

参数：
- `id`: 记录ID
- `provider`: DNS提供商
- `domain_id`: 域名ID
- `status`: 状态，可选值：`enable`、`disable`

### DNS数据库接口

提供将DNS记录存储到本地数据库的接口，路径与DNS管理接口类似，但路径中包含`_db`：

- `GET /api/v1/dns/domains` - 获取数据库中的域名列表
- `POST /api/v1/dns/domains` - 添加域名到数据库
- `GET /api/v1/dns/records_db` - 获取数据库中的DNS记录
- `POST /api/v1/dns/records_db` - 添加DNS记录到数据库
- 等等

### DNS批量操作接口

提供批量操作DNS记录的接口：

- `POST /api/v1/dns/records/batch` - 批量创建DNS记录
- `PUT /api/v1/dns/records/batch` - 批量更新DNS记录
- `DELETE /api/v1/dns/records/batch` - 批量删除DNS记录
- `PUT /api/v1/dns/records/batch/status` - 批量更新DNS记录状态

## 使用示例

启动服务：
```bash
./server
```

服务将运行在配置的端口上（默认8080）。

## 依赖库

- `github.com/gin-gonic/gin`: Web框架
- `github.com/jinzhu/gorm`: ORM库
- `github.com/go-ini/ini`: 配置文件解析
- `github.com/miekg/dns`: DNS查询库
- `github.com/tencentcloud/tencentcloud-sdk-go`: 腾讯云SDK
- `github.com/aliyun/alibaba-cloud-sdk-go`: 阿里云SDK
- `github.com/volcengine/volcengine-go-sdk`: 火山引擎SDK

## 许可证

MIT
