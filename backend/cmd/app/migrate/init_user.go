package migrate

import (
	"github.com/Akvicor/glog"
	"github.com/urfave/cli/v2"
	"wallet/cmd/app/server/common/types/role"
	"wallet/cmd/app/server/service"
)

// 创建默认管理员用户
func initUser(ctx *cli.Context) (err error) {
	username := ctx.String("username")
	password := ctx.String("password")
	nickname := ctx.String("nickname")
	avatar := ctx.String("avatar")
	email := ctx.String("mail")
	phone := ctx.String("phone")

	_, err = service.User.FindByUsername(false, nil, username)
	if err == nil {
		return nil
	}

	// 初始化用户
	_, err = service.User.Create(username, password, nickname, avatar, email, phone, role.TypeAdmin, 1)
	if err != nil {
		glog.Fatal("初始化用户数据异常: %v", err)
		return err
	}
	glog.Info("初始化用户成功, 用户名: %s, 密码: %s", username, password)
	return nil
}
