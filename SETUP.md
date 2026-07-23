# Stock Info 服务 - Kratos初始化完成

## 项目结构

```
stock_info/
├── api/              # API定义（protobuf）
├── cmd/              # 应用程序入口
├── configs/          # 配置文件
├── internal/         # 内部代码
│   ├── biz/          # 业务逻辑层
│   ├── conf/         # 配置结构
│   ├── data/         # 数据访问层
│   ├── server/       # HTTP/gRPC服务器
│   └── service/      # 服务实现层
└── third_party/      # 第三方proto文件
```

## 如何运行

### 1. 编译项目
```bash
go build -o ./bin/ ./...
```

### 2. 运行服务
```bash
./bin/stock_info -conf ./configs
```

服务将监听：
- HTTP: http://0.0.0.0:8000
- gRPC: 0.0.0.0:9000

### 3. 测试API
```bash
# 测试 HTTP 接口
curl http://localhost:8000/helloworld/kratos

# 测试 gRPC 接口（需要安装 grpcurl）
grpcurl -plaintext localhost:9000 list
```

## 开发指南

### 添加新的API
```bash
# 1. 创建 proto 文件
kratos proto add api/stock_info/v1/stock.proto

# 2. 生成代码
kratos proto client api/stock_info/v1/stock.proto

# 3. 生成服务实现
kratos proto server api/stock_info/v1/stock.proto -t internal/service
```

### 使用 Makefile
```bash
make init   # 下载依赖
make api    # 生成API代码
make build  # 编译项目
make all    # 生成所有文件并编译
```

## 配置说明

配置文件位于 `configs/config.yaml`，包含：
- HTTP/gRPC 服务器配置
- 数据库配置（MySQL）
- Redis 配置

## 技术栈

- Go 1.18
- Kratos v2.7.3
- gRPC
- Protobuf
- Wire (依赖注入)

## 数据库迁移 (Goose)

### 安装 Goose
Goose v3.15.0 已安装在 `~/go/bin/goose`

### 配置数据库连接
```bash
# 复制配置文件
cp .goose.env .goose.env.local

# 编辑配置文件，修改数据库连接信息
vim .goose.env.local
```

### 常用命令

#### 使用 Makefile（推荐）
```bash
# 查看迁移状态
make migrate-status

# 执行所有待执行的迁移
make migrate-up

# 回滚最后一个迁移
make migrate-down

# 创建新的迁移文件
make migrate-create name=add_users_table

# 查看当前版本
make migrate-version
```

#### 直接使用 goose 命令
```bash
# 查看状态
goose -dir migrations mysql "root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true" status

# 执行迁移
goose -dir migrations mysql "root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true" up

# 回滚迁移
goose -dir migrations mysql "root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true" down
```

### 已有的迁移文件
1. `20260723160518_create_stocks_table.sql` - 股票信息表
2. `20260723160545_create_stock_prices_table.sql` - 股票价格历史表

### 数据库表结构

#### stocks 表（股票信息）
- 股票代码、名称、市场
- 实时价格信息（当前价、开盘价、最高价、最低价）
- 成交量、市值
- 状态标识

#### stock_prices 表（价格历史）
- 历史K线数据
- 每日开盘、收盘、最高、最低价
- 成交量、成交额
- 涨跌幅

详细说明见：`migrations/README.md`

