// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorm

import (
	"time"

	"github.com/douyu/jupiter/pkg/ecode"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
)

// StdConfig 标准配置，规范配置文件头
func StdConfig(name string) Config {
	return RawConfig("jupiter.mysql." + name)
}

// RawConfig 传入mapstructure格式的配置
// example: RawConfig("jupiter.mysql.stt_config")
func RawConfig(key string) Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		panic(err)
	}

	return config
}

// config options
type Config struct {
	// DSN地址: mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN string `json:"dsn" toml:"dsn"`
	// Debug开关
	Debug bool `json:"debug" toml:"debug"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns" toml:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns" toml:"maxOpenConns"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" toml:"connMaxLifetime"`
	// 创建连接的错误级别，=panic时，如果创建失败，立即panic
	OnDialError string `json:"level" toml:"level"`
	// 慢日志阈值
	SlowThreshold time.Duration `json:"slowThreshold" toml:"slowThreshold"`
	// 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	raw         interface{}
	*xlog.Logger
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		DSN:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: xtime.Duration("300s"),
		OnDialError:     "panic",
		SlowThreshold:   xtime.Duration("500ms"),
		DialTimeout:     xtime.Duration("1s"),
		raw:             nil,
	}
}

// Build ...
func (config Config) Build() *DB {
	db, err := Open("mysql", &config)
	if err != nil {
		xlog.Panic(ecode.MsgClientMysqlOpenPanic, xlog.FieldMod(ecode.ModClientMySQL), xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldValueAny(config))
	}

	if err := db.DB().Ping(); err != nil {
		xlog.Panic(ecode.MsgClientMysqlPingPanic, xlog.FieldMod(ecode.ModClientMySQL), xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldValueAny(config))
	}

	// 上面已经判断过dsn了，这里err可以暂时不判断
	// TODO 将addr，传过来是最好的，先打印数据
	// TODO 上面的value里面有密码，应该用下面解析过数据，过滤掉密码
	d, err := ParseDSN(config.DSN)
	if err == nil {
		xlog.Info(ecode.MsgClientMysqlOpenStart, xlog.FieldMod(ecode.ModClientMySQL), xlog.FieldAddr(d.Addr), xlog.FieldName(d.DBName))
	}

	return db
}
