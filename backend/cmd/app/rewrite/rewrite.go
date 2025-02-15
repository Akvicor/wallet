package rewrite

import (
	"github.com/Akvicor/glog"
	"github.com/urfave/cli/v2"
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

	glog.Info("Rewrite Transaction")
	err = rewriteTransaction()
	if err != nil {
		return err
	}
	return nil
}
