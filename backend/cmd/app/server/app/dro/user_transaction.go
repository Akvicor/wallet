package dro

import (
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/model"
)

/*
View
*/

type UserTransactionViewResp struct {
	IncomeMerge  []*UserTransactionViewRespCategorySum `json:"income_merge"`
	ExpenseMerge []*UserTransactionViewRespCategorySum `json:"expense_merge"`
	Income       []*UserTransactionViewRespCategory    `json:"income"`  // 总收入
	Expense      []*UserTransactionViewRespCategory    `json:"expense"` // 总支出
	List         []*model.UserTransaction              `json:"list"`
}

type UserTransactionViewRespCategory struct {
	Category *model.UserTransactionCategory        `json:"category"`
	Sum      []*UserTransactionViewRespCategorySum `json:"sum"`
}

type UserTransactionViewRespCategorySum struct {
	Currency *model.Currency `json:"currency"`
	Value    decimal.Decimal `json:"value"`
}

/*
View Day
*/

type UserTransactionViewDayResp struct {
	Merge *UserTransactionViewResp            `json:"merge"`
	Days  map[string]*UserTransactionViewResp `json:"days"`
}

/*
View Month
*/

type UserTransactionViewMonthResp struct {
	Merge  *UserTransactionViewResp               `json:"merge"`
	Months map[string]*UserTransactionViewDayResp `json:"months"`
}

/*
View Year
*/

type UserTransactionViewYearResp struct {
	Merge *UserTransactionViewResp                 `json:"merge"`
	Years map[string]*UserTransactionViewMonthResp `json:"years"`
}

/*
Chart
*/

type UserTransactionChartResp []*UserTransactionChartRespItem

type UserTransactionChartRespItem struct {
	Date     string          `json:"date"` // xxxx-xx-xx
	Currency string          `json:"currency"`
	TValue   decimal.Decimal `json:"-"`
	Value    float64         `json:"value"`
}

/*
Pie Chart
*/

type UserTransactionPieChartResp []*UserTransactionPieChartRespItem

type UserTransactionPieChartRespItem struct {
	Type   string          `json:"type"`
	Colour string          `json:"colour"`
	TValue decimal.Decimal `json:"-"`
	Value  float64         `json:"value"`
}
