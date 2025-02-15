package dro

import (
	"fmt"
	"github.com/shopspring/decimal"
	"slices"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

func (u *UserPeriodPayViewResp) ParsePeriodPay(uid int64, total []*model.UserPeriodPay, search []*model.UserPeriodPay) {
	exchange := make(map[string]decimal.Decimal)
	exchangeGet := func(value decimal.Decimal, from, to int64) decimal.Decimal {
		if from == to {
			return value
		}
		exch, ok := exchange[fmt.Sprintf("%d_%d", from, to)]
		if !ok {
			rate, err := service.UserExchangeRate.FindByUID(true, nil, uid, from, to)
			if err != nil {
				rate, err = service.UserExchangeRate.FindByUID(true, nil, uid, to, from)
				if err != nil {
					return decimal.Zero
				}
				exch = decimal.NewFromInt(1).Div(rate.Rate)
			} else {
				exch = rate.Rate
			}
			exchange[fmt.Sprintf("%d_%d", from, to)] = exch
		}
		return value.Mul(exch)
	}
	u.Total = make([]*UserPeriodPayViewRespItem, 0, 2)
	{
		totalMap := make(map[int64]*UserPeriodPayViewRespItem)
		totalGet := func(currency *model.Currency) *UserPeriodPayViewRespItem {
			item, ok := totalMap[currency.ID]
			if !ok {
				item = &UserPeriodPayViewRespItem{
					Currency: currency,
					Day:      decimal.Zero,
					Month:    decimal.Zero,
					Year:     decimal.Zero,
				}
				totalMap[currency.ID] = item
				u.Total = append(u.Total, item)
			}
			return item
		}
		for _, tot := range total {
			item := totalGet(tot.Currency)
			// day
			day := tot.PeriodType.Day(tot.Value, tot.IntervalOfPeriod)
			item.Day = item.Day.Add(day)
			// month
			month := tot.PeriodType.Month(tot.Value, tot.IntervalOfPeriod)
			item.Month = item.Month.Add(month)
			// year
			year := tot.PeriodType.Year(tot.Value, tot.IntervalOfPeriod)
			item.Year = item.Year.Add(year)
		}
		for _, item := range u.Total {
			for _, from := range u.Total {
				item.DayMerge = item.DayMerge.Add(exchangeGet(from.Day, from.Currency.ID, item.Currency.ID))
				item.MonthMerge = item.MonthMerge.Add(exchangeGet(from.Month, from.Currency.ID, item.Currency.ID))
				item.YearMerge = item.YearMerge.Add(exchangeGet(from.Year, from.Currency.ID, item.Currency.ID))
			}
			item.Day = item.Day.Truncate(2)
			item.Month = item.Month.Truncate(2)
			item.Year = item.Year.Truncate(2)
			item.DayMerge = item.DayMerge.Truncate(2)
			item.MonthMerge = item.MonthMerge.Truncate(2)
			item.YearMerge = item.YearMerge.Truncate(2)
		}
		slices.SortFunc(u.Total, func(a, b *UserPeriodPayViewRespItem) int {
			return int(a.Currency.ID - b.Currency.ID)
		})
	}
	u.Search = make([]*UserPeriodPayViewRespItem, 0, 2)
	{
		searchMap := make(map[int64]*UserPeriodPayViewRespItem)
		searchGet := func(currency *model.Currency) *UserPeriodPayViewRespItem {
			item, ok := searchMap[currency.ID]
			if !ok {
				item = &UserPeriodPayViewRespItem{
					Currency: currency,
					Day:      decimal.Zero,
					Month:    decimal.Zero,
					Year:     decimal.Zero,
				}
				searchMap[currency.ID] = item
				u.Search = append(u.Search, item)
			}
			return item
		}
		for _, tot := range search {
			item := searchGet(tot.Currency)
			// day
			day := tot.PeriodType.Day(tot.Value, tot.IntervalOfPeriod)
			item.Day = item.Day.Add(day)
			// month
			month := tot.PeriodType.Month(tot.Value, tot.IntervalOfPeriod)
			item.Month = item.Month.Add(month)
			// year
			year := tot.PeriodType.Year(tot.Value, tot.IntervalOfPeriod)
			item.Year = item.Year.Add(year)
		}
		for _, item := range u.Search {
			for _, from := range u.Search {
				item.DayMerge = item.DayMerge.Add(exchangeGet(from.Day, from.Currency.ID, item.Currency.ID))
				item.MonthMerge = item.MonthMerge.Add(exchangeGet(from.Month, from.Currency.ID, item.Currency.ID))
				item.YearMerge = item.YearMerge.Add(exchangeGet(from.Year, from.Currency.ID, item.Currency.ID))
			}
			item.Day = item.Day.Truncate(2)
			item.Month = item.Month.Truncate(2)
			item.Year = item.Year.Truncate(2)
			item.DayMerge = item.DayMerge.Truncate(2)
			item.MonthMerge = item.MonthMerge.Truncate(2)
			item.YearMerge = item.YearMerge.Truncate(2)
		}
		slices.SortFunc(u.Search, func(a, b *UserPeriodPayViewRespItem) int {
			return int(a.Currency.ID - b.Currency.ID)
		})
	}
}
