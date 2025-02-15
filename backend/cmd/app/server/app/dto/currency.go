package dto

import (
	"wallet/cmd/app/server/common/resp"
)

type CurrencyCreate struct {
	Name        string `json:"name" form:"name" query:"name"`                         // 货币中文名称
	EnglishName string `json:"english_name" form:"english_name" query:"english_name"` // 货币中文名称
	Code        string `json:"code" form:"code" query:"code"`                         // 货币代码
	Symbol      string `json:"symbol" form:"symbol" query:"symbol"`                   // 货币符号
}

type CurrencyFind struct {
	resp.PageModel
	Search string `json:"search" form:"search" query:"search"`
}

type CurrencyUpdate struct {
	Id
	CurrencyCreate
}

type CurrencyDelete struct {
	Id
}
