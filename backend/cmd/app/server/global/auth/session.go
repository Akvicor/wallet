package auth

import "wallet/cmd/app/server/common/sign"

type Session struct {
	CardCVVValid                 string // Empty: 未设置 返回数据不包含CVV, sign.?: 验证成功, other: 正确的验证码
	CardCVVValidLastRequest      int64  // 上次请求的时间戳, 防止短时间重复发送
	CardPasswordValid            string // Empty: 未设置 返回数据不包含密码, sign.?: 验证成功, other: 正确的验证码
	CardPasswordValidLastRequest int64  // 上次请求的时间戳, 防止短时间重复发送
}

func newSession() *Session {
	return &Session{
		CardCVVValid:                 sign.SessionEmpty,
		CardCVVValidLastRequest:      0,
		CardPasswordValid:            sign.SessionEmpty,
		CardPasswordValidLastRequest: 0,
	}
}
