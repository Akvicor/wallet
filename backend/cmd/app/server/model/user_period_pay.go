package model

import (
	"github.com/jinzhu/now"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
	"wallet/cmd/app/server/common/types/period"
)

// UserPeriodPay 周期付费
type UserPeriodPay struct {
	ID               int64           `gorm:"column:id;primaryKey" json:"id"`
	UID              int64           `gorm:"column:uid;index;not null" json:"uid"`
	Name             string          `gorm:"column:name;not null" json:"name"`
	Description      string          `gorm:"column:description" json:"description"`
	CurrencyID       int64           `gorm:"column:currency_id" json:"currency_id"`               // 货币类型
	Value            decimal.Decimal `gorm:"column:value" json:"value"`                           // 金额
	PeriodType       period.Type     `gorm:"column:period_type" json:"period_type"`               // 周期类型
	StartAt          int64           `gorm:"column:start_at" json:"start_at"`                     // 开始时间
	NextOfPeriod     int64           `gorm:"column:next_of_period" json:"next_of_period"`         // 下一次到期时间
	IntervalOfPeriod int64           `gorm:"column:interval_of_period" json:"interval_of_period"` // 间隔
	ExpirationDate   int64           `gorm:"column:expiration_date" json:"expiration_date"`       // 截止日期, -1表示禁用
	ExpirationTimes  int64           `gorm:"column:expiration_times" json:"expiration_times"`     // 截止次数, -1表示禁用

	User     *User     `gorm:"foreignKey:UID;references:ID" json:"user"`
	Currency *Currency `gorm:"foreignKey:CurrencyID;references:ID" json:"currency"`
}

func (*UserPeriodPay) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*UserPeriodPay) TableName() string {
	return "user_period_pay"
}

func NewUserPeriodPay(uid int64, name, description string, currencyID int64, value decimal.Decimal, periodType period.Type, startAt, nextOfPeriod, intervalOfPeriod, expirationDate, expirationTimes int64) *UserPeriodPay {
	return &UserPeriodPay{
		ID:               0,
		UID:              uid,
		Name:             name,
		Description:      description,
		CurrencyID:       currencyID,
		Value:            value,
		PeriodType:       periodType,
		StartAt:          startAt,
		NextOfPeriod:     nextOfPeriod,
		IntervalOfPeriod: intervalOfPeriod,
		ExpirationDate:   expirationDate,
		ExpirationTimes:  expirationTimes,
		Currency:         nil,
	}
}

/**
Preloader
*/

type PreloaderUserPeriodPay struct {
	UserPreload     bool
	CurrencyPreload bool
}

func NewPreloaderUserPeriodPay() *PreloaderUserPeriodPay {
	return &PreloaderUserPeriodPay{
		UserPreload:     false,
		CurrencyPreload: false,
	}
}

func (p *PreloaderUserPeriodPay) Preload(tx *gorm.DB) *gorm.DB {
	if p.UserPreload {
		tx = tx.Preload("User")
	}
	if p.CurrencyPreload {
		tx = tx.Preload("Currency")
	}
	return tx
}

func (p *PreloaderUserPeriodPay) All() *PreloaderUserPeriodPay {
	p.UserPreload = true
	p.CurrencyPreload = true
	return p
}

func (p *PreloaderUserPeriodPay) User() *PreloaderUserPeriodPay {
	p.UserPreload = true
	return p
}

func (p *PreloaderUserPeriodPay) Currency() *PreloaderUserPeriodPay {
	p.CurrencyPreload = true
	return p
}

/*
Function
*/

func (p *UserPeriodPay) NextPeriod() {
	currentTime := now.With(time.Now()).BeginningOfDay()
	current := currentTime.Unix()

	// 减少通知次数
	if p.ExpirationTimes > 0 {
		p.ExpirationTimes -= 1
	}

	if p.PeriodType == period.TypeDaily {
		p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 0, 1).Unix()
		for p.NextOfPeriod < current {
			p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 0, 1).Unix()
		}
		return
	}

	if p.PeriodType == period.TypeWeekly {
		p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 0, 7).Unix()
		for p.NextOfPeriod < current {
			p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 0, 7).Unix()
		}
		return
	}

	if p.PeriodType == period.TypeMonthly {
		p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 1, 0).Unix()
		for p.NextOfPeriod < current {
			p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 1, 0).Unix()
		}
		return
	}

	if p.PeriodType == period.TypeQuarterly {
		p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 3, 0).Unix()
		for p.NextOfPeriod < current {
			p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 3, 0).Unix()
		}
		return
	}

	if p.PeriodType == period.TypeYearly {
		p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(1, 0, 0).Unix()
		for p.NextOfPeriod < current {
			p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(1, 0, 0).Unix()
		}
		return
	}

	if p.PeriodType == period.TypeDayInterval {
		p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 0, int(p.IntervalOfPeriod)).Unix()
		for p.NextOfPeriod < current {
			p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, 0, int(p.IntervalOfPeriod)).Unix()
		}
	}

	if p.PeriodType == period.TypeMonthInterval {
		p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, int(p.IntervalOfPeriod), 0).Unix()
		for p.NextOfPeriod < current {
			p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(0, int(p.IntervalOfPeriod), 0).Unix()
		}
	}

	if p.PeriodType == period.TypeYearInterval {
		p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(int(p.IntervalOfPeriod), 0, 0).Unix()
		for p.NextOfPeriod < current {
			p.NextOfPeriod = time.Unix(p.NextOfPeriod, 0).AddDate(int(p.IntervalOfPeriod), 0, 0).Unix()
		}
	}
}
