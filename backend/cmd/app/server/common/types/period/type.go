package period

import (
	"github.com/jinzhu/now"
	"github.com/shopspring/decimal"
	"time"
)

type Type int64

const (
	TypeDaily         Type = 1 // 每天
	TypeWeekly        Type = 2 // 每周
	TypeMonthly       Type = 3 // 每月
	TypeQuarterly     Type = 4 // 每季度
	TypeYearly        Type = 5 // 每年
	TypeDayInterval   Type = 6 // 间隔天 //
	TypeMonthInterval Type = 7 // 间隔月 //
	TypeYearInterval  Type = 8 // 间隔年 //
)

func (t *Type) Valid() bool {
	switch *t {
	case TypeDaily, TypeWeekly, TypeMonthly, TypeQuarterly, TypeYearly, TypeDayInterval, TypeMonthInterval, TypeYearInterval:
		return true
	default:
		return false
	}
}

func (t Type) String() string {
	switch t {
	case TypeDaily:
		return "每天"
	case TypeWeekly:
		return "每周"
	case TypeMonthly:
		return "每月"
	case TypeQuarterly:
		return "每季度"
	case TypeYearly:
		return "每年"
	case TypeDayInterval:
		return "间隔天"
	case TypeMonthInterval:
		return "间隔月"
	case TypeYearInterval:
		return "间隔年"
	default:
		return "未知"
	}
}

// Day 计算每Type的Value在1天中的值
func (t Type) Day(value decimal.Decimal, interval int64) decimal.Decimal {
	if interval == 0 {
		interval = 1
	}
	switch t {
	case TypeDaily:
		return value
	case TypeWeekly:
		return value.Div(decimal.NewFromInt(7))
	case TypeMonthly:
		return value.Div(decimal.NewFromInt(30))
	case TypeQuarterly:
		return value.Div(decimal.NewFromInt(90))
	case TypeYearly:
		return value.Div(decimal.NewFromInt(365))
	case TypeDayInterval:
		return value.Div(decimal.NewFromInt(interval))
	case TypeMonthInterval:
		return value.Div(decimal.NewFromInt(30)).Div(decimal.NewFromInt(interval))
	case TypeYearInterval:
		return value.Div(decimal.NewFromInt(365)).Div(decimal.NewFromInt(interval))
	default:
		return decimal.Zero
	}
}

// Month 计算每Type的Value在1个月中的值
func (t Type) Month(value decimal.Decimal, interval int64) decimal.Decimal {
	if interval == 0 {
		interval = 1
	}
	switch t {
	case TypeDaily:
		return value.Mul(decimal.NewFromInt(30))
	case TypeWeekly:
		return value.Mul(decimal.NewFromInt(30)).Div(decimal.NewFromInt(7))
	case TypeMonthly:
		return value
	case TypeQuarterly:
		return value.Div(decimal.NewFromInt(3))
	case TypeYearly:
		return value.Div(decimal.NewFromInt(12))
	case TypeDayInterval:
		return value.Mul(decimal.NewFromInt(30)).Div(decimal.NewFromInt(interval))
	case TypeMonthInterval:
		return value.Div(decimal.NewFromInt(interval))
	case TypeYearInterval:
		return value.Div(decimal.NewFromInt(12)).Div(decimal.NewFromInt(interval))
	default:
		return decimal.Zero
	}
}

// Year 计算每Type的Value在1年中的值
func (t Type) Year(value decimal.Decimal, interval int64) decimal.Decimal {
	if interval == 0 {
		interval = 1
	}
	switch t {
	case TypeDaily:
		return value.Mul(decimal.NewFromInt(365))
	case TypeWeekly:
		return value.Mul(decimal.NewFromInt(365)).Div(decimal.NewFromInt(7))
	case TypeMonthly:
		return value.Mul(decimal.NewFromInt(12))
	case TypeQuarterly:
		return value.Mul(decimal.NewFromInt(4))
	case TypeYearly:
		return value
	case TypeDayInterval:
		return value.Mul(decimal.NewFromInt(365)).Div(decimal.NewFromInt(interval))
	case TypeMonthInterval:
		return value.Mul(decimal.NewFromInt(12)).Div(decimal.NewFromInt(interval))
	case TypeYearInterval:
		return value.Div(decimal.NewFromInt(interval))
	default:
		return decimal.Zero
	}
}

func (t Type) StringEnglish() string {
	switch t {
	case TypeDaily:
		return "Daily"
	case TypeWeekly:
		return "Weekly"
	case TypeMonthly:
		return "Monthly"
	case TypeQuarterly:
		return "Quarterly"
	case TypeYearly:
		return "Yearly"
	case TypeDayInterval:
		return "DayInterval"
	case TypeMonthInterval:
		return "MonthInterval"
	case TypeYearInterval:
		return "YearInterval"
	default:
		return "Unknown"
	}
}

func (t *Type) Calculate(startAt, nextOfPeriod, intervalOfPeriod int64) (start, next, interval int64) {
	startAtTime := now.With(time.Unix(startAt, 0)).BeginningOfDay()
	start = startAtTime.Unix()

	switch *t {
	case TypeDaily:
		next = start
		interval = 0
	case TypeWeekly:
		next = start
		interval = 0
	case TypeMonthly:
		next = start
		interval = 0
	case TypeQuarterly:
		next = start
		interval = 0
	case TypeYearly:
		next = start
		interval = 0
	case TypeDayInterval:
		next = start
		interval = intervalOfPeriod
	case TypeMonthInterval:
		next = start
		interval = intervalOfPeriod
	case TypeYearInterval:
		next = start
		interval = intervalOfPeriod
	}
	if nextOfPeriod > next {
		next = nextOfPeriod
	}
	return start, next, interval
}

type TypeDetail struct {
	Type        Type   `json:"type"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
}

var AllType = []*TypeDetail{
	{TypeDaily, TypeDaily.String(), TypeDaily.StringEnglish()},
	{TypeWeekly, TypeWeekly.String(), TypeWeekly.StringEnglish()},
	{TypeMonthly, TypeMonthly.String(), TypeMonthly.StringEnglish()},
	{TypeQuarterly, TypeQuarterly.String(), TypeQuarterly.StringEnglish()},
	{TypeYearly, TypeYearly.String(), TypeYearly.StringEnglish()},
	{TypeDayInterval, TypeDayInterval.String(), TypeDayInterval.StringEnglish()},
	{TypeMonthInterval, TypeMonthInterval.String(), TypeMonthInterval.StringEnglish()},
	{TypeYearInterval, TypeYearInterval.String(), TypeYearInterval.StringEnglish()},
}
