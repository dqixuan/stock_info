package configs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Http  HttpServer `json:"http" mapstructure:"http"`
	Grpc  GrpcServer `json:"grpc" mapstructure:"grpc"`
	Mysql Mysql      `json:"mysql" mapstructure:"mysql"`
}

type HttpServer struct {
	Addr    string        `json:"addr" mapstructure:"addr"`
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`
}

type GrpcServer struct {
	Addr    string        `json:"addr" mapstructure:"addr"`
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`
}

type Mysql struct {
	Root          string        `json:"root" mapstructure:"root"`
	ReadTimeout   time.Duration `json:"read_timeout" mapstructure:"read_timeout"`
	WriterTimeout time.Duration `json:"write_timeout" mapstructure:"write_timeout"`
	MaxIdleConns  int           `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns  int           `json:"max_open_conns" mapstructure:"max_open_conns"`
	LogLevel      int           `json:"log_level" mapstructure:"log_level"`
}

var GlobalConfig *Config

func NewConfig(path string) error {
	v := viper.New()

	if path == "" {
		path = "./configs/config.yaml"
	}
	v.SetConfigFile(path)

	// 4. 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read file failed: %v", err)
	}

	// 5. 将配置解析到结构体中
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("parse config failed: %v", err)
	}
	fmt.Printf("get config success: %+v\n", cfg)
	GlobalConfig = &cfg
	return nil
}

func NewHttpCfg() *HttpServer {
	return &GlobalConfig.Http
}

func NewGrpcCfg() *GrpcServer {
	return &GlobalConfig.Grpc
}

func NewMysqlCfg() *Mysql {
	return &GlobalConfig.Mysql
}
