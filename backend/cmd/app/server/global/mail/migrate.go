package mail

import (
	"fmt"
	"wallet/cmd/config"
)

func Migrate(to, nickname, username string) {
	subject := fmt.Sprintf("%s 初始化通知", config.Global.AppName)
	text := fmt.Sprintf(`您好，%s。
	初始化成功, 默认管理员账户信息如下。
	账号：%s
`, nickname, username)
	_ = Send(to, subject, Plain, text)
}
