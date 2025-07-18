package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource struct {
		Read  string
		Write string
	}
	MigrationPath string
	Cache cache.CacheConf
}
