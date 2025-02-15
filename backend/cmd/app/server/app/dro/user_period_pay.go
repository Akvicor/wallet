package dro

import (
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/model"
)

type UserPeriodPayViewResp struct {
	Total []*UserPeriodPayViewRespItem `json:"total"` // 总计
	//TotalMerge  []*UserPeriodPayViewRespItem `json:"total_merge"`  // 总计, 合并货币
	Search []*UserPeriodPayViewRespItem `json:"search"` // 过滤
	//SearchMerge []*UserPeriodPayViewRespItem `json:"search_merge"` // 过滤, 合并货币
}

type UserPeriodPayViewRespItem struct {
	Currency   *model.Currency `json:"currency"`
	Day        decimal.Decimal `json:"day"`
	Month      decimal.Decimal `json:"month"`
	Year       decimal.Decimal `json:"year"`
	DayMerge   decimal.Decimal `json:"day_merge"`
	MonthMerge decimal.Decimal `json:"month_merge"`
	YearMerge  decimal.Decimal `json:"year_merge"`
}
