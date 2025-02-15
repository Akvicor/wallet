package server

import (
	"github.com/urfave/cli/v2"
	"wallet/cmd/app/server/app"
	"wallet/cmd/app/server/common/cache"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/app/server/global/db"
	"wallet/cmd/config"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "config",
		Usage:   "config file path",
		Value:   "./data/config.toml",
		Aliases: []string{"c"},
	},
}

func Action(ctx *cli.Context) (err error) {
	// 加载配置文件
	config.Load(ctx.String("config"))
	// 加载数据库
	db.Load()
	// 配置缓存
	cache.SetTokenManagerOnEvicted(auth.OnTokenEvicted)
	// 运行服务
	return app.Run()
}
