package dto

import "github.com/shopspring/decimal"

type UserCardCurrencyCreate struct {
	UserCardID int64 `json:"user_card_id" form:"user_card_id" query:"user_card_id"` // 银行卡ID
	CurrencyID int64 `json:"currency_id" form:"currency_id" query:"currency_id"`    // 货币ID
}

type UserCardCurrencyUpdateBalance struct {
	Id
	Balance decimal.Decimal `json:"balance" form:"balance" query:"balance"` // 划分资金余额
}

type UserCardCurrencyDelete struct {
	Id
}
