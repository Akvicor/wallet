package dto

import (
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/common/resp"
)

type UserExchangeRateCreate struct {
	FromCurrencyID int64           `json:"from_currency_id" form:"from_currency_id" query:"from_currency_id"` // 原货币ID
	ToCurrencyID   int64           `json:"to_currency_id" form:"to_currency_id" query:"to_currency_id"`       // 目标货币ID
	Rate           decimal.Decimal `json:"rate" form:"rate" query:"rate"`                                     // 汇率
}

type UserExchangeRateFind struct {
	resp.PageModel
}

type UserExchangeRateUpdate struct {
	Id
	UserExchangeRateCreate
}

type UserExchangeRateDelete struct {
	Id
}
