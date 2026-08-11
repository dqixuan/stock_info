GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;

wire:
	@echo "Generating wire..."
	@wire ./...
	wire ./...

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# Database migration settings
DB_DRIVER ?= mysql
DB_STRING ?= root:root@tcp(127.0.0.1:3306)/stock_info?parseTime=true

.PHONY: migrate-status
# show migration status
migrate-status:
	@echo "Checking migration status..."
	@goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" status || echo "Run 'make migrate-up' to apply migrations"

.PHONY: migrate-up
# apply all pending migrations
migrate-up:
	@echo "Applying migrations..."
	@goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" up

.PHONY: migrate-up-one
# apply next migration
migrate-up-one:
	@echo "Applying next migration..."
	@goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" up-by-one

.PHONY: migrate-down
# rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" down

.PHONY: migrate-reset
# rollback all migrations
migrate-reset:
	@echo "Resetting all migrations..."
	@goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" reset

.PHONY: migrate-create
# create new migration file (usage: make migrate-create name=add_users_table)
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migrate-create name=add_users_table"; \
		exit 1; \
	fi
	@goose -dir migrations create $(name) sql
	@echo "Created new migration: $(name)"

.PHONY: migrate-version
# show current migration version
migrate-version:
	@goose -dir migrations $(DB_DRIVER) "$(DB_STRING)" version

