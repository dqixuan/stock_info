# Database Migrations

此目录包含使用 [goose](https://github.com/pressly/goose) 管理的数据库迁移文件。

## 快速开始

### 1. 配置数据库连接

复制配置文件并修改数据库连接信息：
```bash
cp .goose.env .goose.env.local
# 编辑 .goose.env.local 文件，配置你的数据库连接
```

### 2. 运行迁移

```bash
# 查看当前状态
goose -dir migrations mysql "root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true" status

# 执行所有待执行的迁移（升级）
goose -dir migrations mysql "root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true" up

# 只执行下一个迁移
goose -dir migrations mysql "root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true" up-by-one

# 回滚最后一个迁移
goose -dir migrations mysql "root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true" down

# 重置数据库（回滚所有迁移）
goose -dir migrations mysql "root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true" reset
```

### 3. 创建新的迁移

```bash
# 创建SQL迁移文件
goose -dir migrations create add_user_table sql

# 创建Go迁移文件（用于复杂的迁移逻辑）
goose -dir migrations create seed_initial_data go
```

## 使用Makefile（推荐）

在项目根目录的Makefile中添加以下命令：

```makefile
# Database connection (override with environment variables)
DB_DRIVER ?= mysql
DB_STRING ?= root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true

# Migration commands
.PHONY: migrate-status
migrate-status:
	goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" status

.PHONY: migrate-up
migrate-up:
	goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" up

.PHONY: migrate-down
migrate-down:
	goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" down

.PHONY: migrate-reset
migrate-reset:
	goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" reset

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql
```

然后可以使用：
```bash
make migrate-status
make migrate-up
make migrate-down
make migrate-create
```

## 迁移文件命名规范

迁移文件格式：`YYYYMMDDHHMMSS_description.sql`

例如：
- `20260723160518_create_stocks_table.sql`
- `20260723160545_create_stock_prices_table.sql`

## 迁移文件结构

SQL迁移文件必须包含 Up 和 Down 部分：

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
```

## 已有的迁移

1. `20260723160518_create_stocks_table.sql` - 创建股票信息表
2. `20260723160545_create_stock_prices_table.sql` - 创建股票价格历史表

## 最佳实践

1. **总是先测试迁移**：在开发环境测试迁移后再应用到生产环境
2. **可回滚性**：确保每个Up迁移都有对应的Down迁移
3. **原子性**：每个迁移文件应该只做一件事
4. **版本控制**：将迁移文件纳入版本控制
5. **不要修改已应用的迁移**：如果需要更改，创建新的迁移文件
6. **数据迁移**：对于包含数据的迁移，考虑使用事务

## 故障排除

### 错误：migration failed

检查：
1. 数据库连接是否正确
2. SQL语法是否正确
3. 是否有权限执行操作
4. 是否有表名冲突

### 查看已应用的迁移

```bash
# 查看goose_db_version表
mysql -u root -p -e "SELECT * FROM stock_info.goose_db_version ORDER BY id;"
```

## 参考文档

- [Goose官方文档](https://github.com/pressly/goose)
- [MySQL数据类型](https://dev.mysql.com/doc/refman/8.0/en/data-types.html)
