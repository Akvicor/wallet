package sms

import (
	"fmt"
	"wallet/cmd/config"
)

func UserCreate(to, nickname, username string) {
	msg := fmt.Sprintf(`%s 注册通知
您好，%s。
管理员为你开通了账户。
账号：%s`, config.Global.AppName, nickname, username)
	Send(to, msg)
}
