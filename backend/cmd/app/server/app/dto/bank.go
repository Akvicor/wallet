package dto

import (
	"wallet/cmd/app/server/common/resp"
)

type BankCreate struct {
	Name        string `json:"name" form:"name" query:"name"`                         // 中文名称
	EnglishName string `json:"english_name" form:"english_name" query:"english_name"` // 英文名称
	Abbr        string `json:"abbr" form:"abbr" query:"abbr"`                         // 缩写
	Phone       string `json:"phone" form:"phone" query:"phone"`                      // 银行电话
}

type BankFind struct {
	resp.PageModel
	Search string `json:"search" form:"search" query:"search"`
}

type BankUpdate struct {
	Id
	BankCreate
}

type BankDelete struct {
	Id
}
