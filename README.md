# Go-Gin Example with DNSPod and Aliyun DNS API

这是一个集成了DNSPod和阿里云DNS API功能的Go-Gin示例项目。

## 功能特性

- 原有的标签管理功能
- 新增的DNS API接口，支持：
  - DNSPod API接口
  - 阿里云DNS API接口
  - 支持多种DNS服务提供商
  - 数据库存储功能（域名列表、解析记录等）
  - DNS解析批量操作功能

## 配置

在 `conf/app.ini` 文件中配置DNS相关参数：

```ini
[dns]
DNSPOD_TOKEN = your_dns_pod_token_here  # DNSPod API Token，格式为 ID,Token
DOMAIN_NAME = example.com               # 默认域名

[aliyun_dns]
ALIYUN_ACCESS_KEY_ID = your_aliyun_access_key_id      # 阿里云AccessKey ID
ALIYUN_ACCESS_KEY_SECRET = your_aliyun_access_key_secret  # 阿里云AccessKey Secret
ALIYUN_REGION_ID = cn-hangzhou                        # 阿里云区域ID，默认为cn-hangzhou
```

## API接口

### 原有标签API接口
- `GET /api/v1/tags` - 获取标签列表
- `POST /api/v1/tags` - 创建标签
- `PUT /api/v1/tags/:id` - 更新标签
- `DELETE /api/v1/tags/:id` - 删除标签

### DNS API接口

#### 云服务提供商API接口
- `GET /api/v1/domains` - 获取域名列表
- `GET /api/v1/dns/records` - 获取DNS记录列表
- `POST /api/v1/dns/records` - 创建DNS记录
- `PUT /api/v1/dns/records/:id` - 更新DNS记录
- `DELETE /api/v1/dns/records/:id` - 删除DNS记录
- `PUT /api/v1/dns/records/:id/status` - 设置DNS记录状态

#### 数据库存储API接口
- `GET /api/v1/dns/domains` - 获取域名列表（数据库）
- `POST /api/v1/dns/domains` - 添加域名（数据库）
- `PUT /api/v1/dns/domains/:id` - 更新域名（数据库）
- `DELETE /api/v1/dns/domains/:id` - 删除域名（数据库）
- `GET /api/v1/dns/records_db` - 获取DNS解析记录列表（数据库）
- `POST /api/v1/dns/records_db` - 添加DNS解析记录（数据库）
- `PUT /api/v1/dns/records_db/:id` - 更新DNS解析记录（数据库）
- `DELETE /api/v1/dns/records_db/:id` - 删除DNS解析记录（数据库）

#### DNS解析批量操作API接口
- `POST /api/v1/dns/records/batch` - 批量创建DNS记录
- `PUT /api/v1/dns/records/batch` - 批量更新DNS记录
- `DELETE /api/v1/dns/records/batch` - 批量删除DNS记录
- `PUT /api/v1/dns/records/batch/status` - 批量更新DNS记录状态

### 参数说明

#### 云服务提供商API参数
- **通用参数**:
  - `provider` - DNS服务提供商 (dns_pod 或 aliyun)，默认为dns_pod

- **获取DNS记录列表**: 
  - `domain` - 域名 (可选，使用配置中的默认域名)
  - `sub_domain` - 子域名 (可选)
  - `provider` - DNS服务提供商

- **创建DNS记录**:
  - `domain_id` - 域名ID (DNSPod) 或域名名称 (阿里云)
  - `sub_domain` - 子域名
  - `record_type` - 记录类型 (A, CNAME, MX等)
  - `value` - 记录值
  - `record_line` - 线路 (DNSPod, 默认为"默认")
  - `ttl` - TTL值 (阿里云, 默认为600)
  - `provider` - DNS服务提供商

- **更新DNS记录**:
  - `id` - 记录ID (路径参数)
  - `domain_id` - 域名ID (DNSPod)
  - `sub_domain` - 子域名
  - `record_type` - 记录类型
  - `value` - 记录值
  - `record_line` - 线路 (DNSPod)
  - `ttl` - TTL值 (阿里云)
  - `provider` - DNS服务提供商

- **删除DNS记录**:
  - `id` - 记录ID (路径参数)
  - `domain_id` - 域名ID (DNSPod)
  - `provider` - DNS服务提供商

- **设置DNS记录状态**:
  - `id` - 记录ID (路径参数)
  - `domain_id` - 域名ID (DNSPod)
  - `status` - 状态 (enable/disable)
  - `provider` - DNS服务提供商

#### 数据库API参数
- **域名管理参数**:
  - `name` - 域名
  - `provider` - 服务提供商 (dns_pod, aliyun)
  - `domain_id` - 云服务商域名ID
  - `status` - 状态 (active, inactive)
  - `grade` - 域名等级
  - `owner` - 域名所有者
  - `remark` - 备注

- **DNS解析记录管理参数**:
  - `domain_id` - 关联域名ID（数据库中的ID）
  - `name` - 记录名称（如 www）
  - `type` - 记录类型（A, CNAME, MX等）
  - `value` - 记录值（如IP地址）
  - `status` - 状态 (enable, disable)
  - `line` - 线路
  - `ttl` - TTL值
  - `remark` - 备注
  - `provider` - 服务提供商
  - `remote_id` - 云服务商记录ID

#### 批量操作API参数
- **批量创建DNS记录** (`POST /api/v1/dns/records/batch`):
  - `provider` - DNS服务提供商 (dns_pod 或 aliyun)
  - **请求体**:
    ```json
    [
      {
        "domain_id": "域名ID",
        "name": "记录名称",
        "type": "记录类型",
        "value": "记录值",
        "line": "线路 (DNSPod)",
        "ttl": 600,
        "remark": "备注"
      }
    ]
    ```

- **批量更新DNS记录** (`PUT /api/v1/dns/records/batch`):
  - `provider` - DNS服务提供商 (dns_pod 或 aliyun)
  - **请求体**:
    ```json
    [
      {
        "id": "记录ID",
        "domain_id": "域名ID",
        "name": "记录名称",
        "type": "记录类型",
        "value": "记录值",
        "line": "线路 (DNSPod)",
        "ttl": 600,
        "remark": "备注"
      }
    ]
    ```

- **批量删除DNS记录** (`DELETE /api/v1/dns/records/batch`):
  - `provider` - DNS服务提供商 (dns_pod 或 aliyun)
  - **请求体**:
    ```json
    [
      {
        "id": "记录ID",
        "domain_id": "域名ID"
      }
    ]
    ```

- **批量更新DNS记录状态** (`PUT /api/v1/dns/records/batch/status`):
  - `provider` - DNS服务提供商 (dns_pod 或 aliyun)
  - `status` - 状态 (enable/disable)
  - **请求体**:
    ```json
    [
      {
        "id": "记录ID",
        "domain_id": "域名ID"
      }
    ]
    ```

## 使用示例

### 云服务提供商API使用示例

#### 获取DNSPod记录列表
```bash
curl -X GET "http://localhost:8000/api/v1/dns/records?domain=example.com&provider=dns_pod"
```

#### 获取阿里云DNS记录列表
```bash
curl -X GET "http://localhost:8000/api/v1/dns/records?domain=example.com&provider=aliyun"
```

#### 创建DNSPod记录
```bash
curl -X POST "http://localhost:8000/api/v1/dns/records?domain_id=123456&sub_domain=www&record_type=A&value=1.2.3.4&provider=dns_pod"
```

#### 创建阿里云DNS记录
```bash
curl -X POST "http://localhost:8000/api/v1/dns/records?domain_id=example.com&sub_domain=www&record_type=A&value=1.2.3.4&ttl=600&provider=aliyun"
```

### 数据库API使用示例

#### 获取域名列表（数据库）
```bash
curl -X GET "http://localhost:8000/api/v1/dns/domains"
```

#### 添加域名（数据库）
```bash
curl -X POST "http://localhost:8000/api/v1/dns/domains?name=example.com&provider=dns_pod&domain_id=123456&status=active"
```

#### 获取DNS解析记录列表（数据库）
```bash
curl -X GET "http://localhost:8000/api/v1/dns/records_db?domain_id=1"
```

#### 添加DNS解析记录（数据库）
```bash
curl -X POST "http://localhost:8000/api/v1/dns/records_db?domain_id=1&name=www&type=A&value=1.2.3.4&provider=dns_pod&status=enable"
```

### 批量操作API使用示例

#### 批量创建DNS记录
```bash
curl -X POST "http://localhost:8000/api/v1/dns/records/batch?provider=dns_pod" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "domain_id": "123456",
      "name": "www",
      "type": "A",
      "value": "1.2.3.4",
      "line": "默认",
      "ttl": 600
    },
    {
      "domain_id": "123456",
      "name": "mail",
      "type": "MX",
      "value": "mail.example.com",
      "line": "默认",
      "ttl": 600
    }
  ]'
```

#### 批量更新DNS记录
```bash
curl -X PUT "http://localhost:8000/api/v1/dns/records/batch?provider=dns_pod" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "id": "789012",
      "domain_id": "123456",
      "name": "www",
      "type": "A",
      "value": "2.3.4.5",
      "line": "默认",
      "ttl": 1200
    },
    {
      "id": "789013",
      "domain_id": "123456",
      "name": "mail",
      "type": "MX",
      "value": "newmail.example.com",
      "line": "默认",
      "ttl": 600
    }
  ]'
```

#### 批量删除DNS记录
```bash
curl -X DELETE "http://localhost:8000/api/v1/dns/records/batch?provider=dns_pod" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "id": "789012",
      "domain_id": "123456"
    },
    {
      "id": "789013",
      "domain_id": "123456"
    }
  ]'
```

#### 批量更新DNS记录状态
```bash
curl -X PUT "http://localhost:8000/api/v1/dns/records/batch/status?provider=dns_pod&status=disable" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "id": "789012",
      "domain_id": "123456"
    },
    {
      "id": "789013",
      "domain_id": "123456"
    }
  ]'
```

## 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8000` 启动。

## 数据库初始化

项目使用GORM进行数据库操作，支持MySQL数据库。确保在 `conf/app.ini` 中正确配置数据库连接信息。

## DNSPod Token获取

1. 登录DNSPod控制台（现在是腾讯云的一部分）
2. 进入"用户中心" -> "安全设置" -> "API密钥" -> "API Token管理"
3. 创建新的API Token
4. 格式为 `ID,Token`，如 `12345,abc12345def`

**重要提示**: 请确保在 `conf/app.ini` 中正确配置 `DNSPOD_TOKEN`，格式为 `Token ID,Token`。

## 阿里云AccessKey获取

1. 登录阿里云控制台
2. 进入"访问控制(RAM)" -> "用户管理" -> "用户" -> "安全凭证"
3. 创建AccessKey ID和AccessKey Secret
4. 配置在 `conf/app.ini` 中的 `ALIYUN_ACCESS_KEY_ID` 和 `ALIYUN_ACCESS_KEY_SECRET`

**重要提示**: 
- 请妥善保管AccessKey，不要泄露
- 建议创建RAM用户并分配最小权限
- 阿里云DNS API需要相应的权限策略

## 错误码

- `200` - 成功
- `500` - 服务器错误
- `400` - 参数错误