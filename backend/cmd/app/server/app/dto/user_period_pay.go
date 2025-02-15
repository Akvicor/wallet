package dto

import (
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/period"
)

type UserPeriodPayCreate struct {
	Name             string          `json:"name" form:"name" query:"name"`                                           // 名称
	Description      string          `json:"description" form:"description" query:"description"`                      // 描述
	CurrencyID       int64           `json:"currency_id" form:"currency_id" query:"currency_id"`                      // 货币类型
	Value            decimal.Decimal `json:"value" form:"value" query:"value"`                                        // 金额
	PeriodType       period.Type     `json:"period_type" form:"period_type" query:"period_type"`                      // 周期类型
	StartAt          int64           `json:"start_at" form:"start_at" query:"start_at"`                               // 开始时间
	NextOfPeriod     int64           `json:"next_of_period" form:"next_of_period" query:"next_of_period"`             // 下一次到期时间
	IntervalOfPeriod int64           `json:"interval_of_period" form:"interval_of_period" query:"interval_of_period"` // 间隔
	ExpirationDate   int64           `json:"expiration_date" form:"expiration_date" query:"expiration_date"`          // 过期时间
	ExpirationTimes  int64           `json:"expiration_times" form:"expiration_times" query:"expiration_times"`       // 次数
}

type UserPeriodPayFind struct {
	resp.PageModel
	Id
	All    bool   `json:"all" form:"all" query:"all"`
	Search string `json:"search" form:"search" query:"search"`
}

type UserPeriodPaySummary struct {
	UserPeriodPayFind
}

type UserPeriodPayUpdate struct {
	Id
	UserPeriodPayCreate
}

type UserPeriodPayUpdateNextPeriod struct {
	Id
}

type UserPeriodPayDisable struct {
	Id
}

type UserPeriodPayEnable struct {
	Id
}
