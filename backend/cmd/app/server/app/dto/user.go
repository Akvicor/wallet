package dto

import (
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/role"
)

type UserCreate struct {
	Username         string    `json:"username" form:"username" query:"username"`                               // 用户名
	Password         string    `json:"password" form:"password" query:"password"`                               // 密码
	Nickname         string    `json:"nickname" form:"nickname" query:"nickname"`                               // 昵称
	Avatar           string    `json:"avatar" form:"avatar" query:"avatar"`                                     // 头像
	Mail             string    `json:"mail" form:"mail" query:"mail"`                                           // 邮箱
	Phone            string    `json:"phone" form:"phone" query:"phone"`                                        // 手机
	Role             role.Type `json:"role" form:"role" query:"role"`                                           // 角色
	MasterCurrencyID int64     `json:"master_currency_id" form:"master_currency_id" query:"master_currency_id"` // 主要货币
}

type UserFind struct {
	resp.PageModel
	Id
	Search string `json:"search" form:"search" query:"search"`
}

type UserLoginLogFind struct {
	resp.PageModel
	Id
	Search string `json:"search" form:"search" query:"search"`
}

type UserUpdate struct {
	Id
	UserCreate
}

type UserLogin struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
	Remember bool   `json:"remember" form:"remember" query:"remember"`
}

type UserCreateAccessToken struct {
	Name string `json:"name" form:"name" query:"name"`
}

type UserDisableEnable struct {
	Id
	Username string `json:"username" form:"username" query:"username"`
}

type UserInfoUpdate struct {
	Username         string `json:"username" form:"username" query:"username"`                               // 用户名
	Password         string `json:"password" form:"password" query:"password"`                               // 密码
	Nickname         string `json:"nickname" form:"nickname" query:"nickname"`                               // 昵称
	Avatar           string `json:"avatar" form:"avatar" query:"avatar"`                                     // 头像
	Mail             string `json:"mail" form:"mail" query:"mail"`                                           // 邮箱
	Phone            string `json:"phone" form:"phone" query:"phone"`                                        // 手机
	MasterCurrencyID int64  `json:"master_currency_id" form:"master_currency_id" query:"master_currency_id"` // 主要货币
}
