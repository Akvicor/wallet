package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserExchangeRate 用户货币汇率
type UserExchangeRate struct {
	ID             int64           `gorm:"column:id;primaryKey" json:"id"`
	UID            int64           `gorm:"column:uid;uniqueIndex:uid_from_currency_id_to_currency_id;index" json:"uid"`                           // 用户ID
	FromCurrencyID int64           `gorm:"column:from_currency_id;uniqueIndex:uid_from_currency_id_to_currency_id;index" json:"from_currency_id"` // 货币种类ID
	ToCurrencyID   int64           `gorm:"column:to_currency_id;uniqueIndex:uid_from_currency_id_to_currency_id;index" json:"to_currency_id"`     // 货币种类ID
	Rate           decimal.Decimal `gorm:"column:rate" json:"rate"`                                                                               // 从From到To的汇率

	FromCurrency *Currency `gorm:"foreignKey:FromCurrencyID;references:ID" json:"from_currency"`
	ToCurrency   *Currency `gorm:"foreignKey:ToCurrencyID;references:ID" json:"to_currency"`
}

func (*UserExchangeRate) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*UserExchangeRate) TableName() string {
	return "user_exchange_rate"
}

func NewUserExchangeRate(uid, fromCurrencyID, toCurrencyID int64, rate decimal.Decimal) *UserExchangeRate {
	return &UserExchangeRate{
		ID:             0,
		UID:            uid,
		FromCurrencyID: fromCurrencyID,
		ToCurrencyID:   toCurrencyID,
		Rate:           rate,
	}
}

/**
Preloader
*/

type PreloaderUserExchangeRate struct {
	FromCurrencyPreload bool
	ToCurrencyPreload   bool
}

func NewPreloaderUserExchangeRate() *PreloaderUserExchangeRate {
	return &PreloaderUserExchangeRate{
		FromCurrencyPreload: false,
		ToCurrencyPreload:   false,
	}
}

func (p *PreloaderUserExchangeRate) Preload(tx *gorm.DB) *gorm.DB {
	if p.FromCurrencyPreload {
		tx = tx.Preload("FromCurrency")
	}
	if p.ToCurrencyPreload {
		tx = tx.Preload("ToCurrency")
	}
	return tx
}

func (p *PreloaderUserExchangeRate) All() *PreloaderUserExchangeRate {
	p.FromCurrencyPreload = true
	p.ToCurrencyPreload = true
	return p
}

func (p *PreloaderUserExchangeRate) FromCurrency() *PreloaderUserExchangeRate {
	p.FromCurrencyPreload = true
	return p
}

func (p *PreloaderUserExchangeRate) ToCurrency() *PreloaderUserExchangeRate {
	p.ToCurrencyPreload = true
	return p
}
