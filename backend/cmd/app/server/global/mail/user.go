package mail

import (
	"fmt"
	"wallet/cmd/config"
)

func UserCreate(to, nickname, username string) {
	subject := fmt.Sprintf("%s 注册通知", config.Global.AppName)
	text := fmt.Sprintf(`您好，%s。
	管理员为你开通了账户。
	账号：%s
`, nickname, username)
	_ = Send(to, subject, Plain, text)
}
