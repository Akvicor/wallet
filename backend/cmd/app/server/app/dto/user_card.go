package dto

import (
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/common/resp"
)

type UserCardCreate struct {
	BankID              int64           `json:"bank_id" form:"bank_id" query:"bank_id"`                                           // 银行卡所属钱包
	TypeID              []int64         `json:"type_id" form:"type_id" query:"type_id"`                                           // 银行卡类型
	Name                string          `json:"name" form:"name" query:"name"`                                                    // 银行卡名称
	Description         string          `json:"description" form:"description" query:"description"`                               // 描述
	Number              string          `json:"number" form:"number" query:"number"`                                              // 银行卡号
	ExpDate             int64           `json:"exp_date" form:"exp_date" query:"exp_date"`                                        // 过期时间
	CVV                 string          `json:"cvv" form:"cvv" query:"cvv"`                                                       // 安全码
	StatementClosingDay int64           `json:"statement_closing_day" form:"statement_closing_day" query:"statement_closing_day"` // 账单日
	PaymentDueDay       int64           `json:"payment_due_day" form:"payment_due_day" query:"payment_due_day"`                   // 还款日
	Password            string          `json:"password" form:"password" query:"password"`                                        // 密码
	MasterCurrencyID    int64           `json:"master_currency_id" form:"master_currency_id" query:"master_currency_id"`          // 密码
	CurrencyID          []int64         `json:"currency_id" form:"currency_id" query:"currency_id"`                               // 银行卡支持的货币
	Limit               decimal.Decimal `json:"limit" form:"limit" query:"limit"`                                                 // 限额
	Fee                 string          `json:"fee" form:"fee" query:"fee"`                                                       // 管理费/年费
	HideBalance         bool            `json:"hide_balance" form:"hide_balance" query:"hide_balance"`                            // 隐藏余额
}

type UserCardFind struct {
	resp.PageModel
	Id
	All    bool   `json:"all" form:"all" query:"all"`
	Search string `json:"search" form:"search" query:"search"`
}

type UserCardValidRequest struct {
	Key    string `json:"key" form:"key" query:"key"`
	Method string `json:"method" form:"method" query:"method"`
}

type UserCardValidInput struct {
	Key        string `json:"key" form:"key" query:"key"`
	VerifyCode string `json:"verify_code" form:"verify_code" query:"verify_code"`
}

type UserCardValidCancel struct {
	Key string `json:"key" form:"key" query:"key"`
}

type UserCardUpdate struct {
	Id
	UserCardCreate
}

type UserCardUpdateSequence struct {
	Id
	Target int64 `json:"target" form:"target" query:"target"` // 目标序号
}

type UserCardDisable struct {
	Id
}

type UserCardEnable struct {
	Id
}
