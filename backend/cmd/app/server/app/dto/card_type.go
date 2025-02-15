package dto

import (
	"wallet/cmd/app/server/common/resp"
)

type CardTypeCreate struct {
	Name string `json:"name" form:"name" query:"name"` // 银行卡类型
}

type CardTypeFind struct {
	resp.PageModel
	Search string `json:"search" form:"search" query:"search"`
}

type CardTypeUpdate struct {
	Id
	CardTypeCreate
}

type CardTypeDelete struct {
	Id
}
