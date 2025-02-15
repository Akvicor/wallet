package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserCardCurrency 用户银行卡各类型货币余额
type UserCardCurrency struct {
	ID         int64           `gorm:"column:id;primaryKey" json:"id"`
	UserCardID int64           `gorm:"column:user_card_id;uniqueIndex:user_card_id_currency_id;index" json:"user_card_id"`
	CurrencyID int64           `gorm:"column:currency_id;uniqueIndex:user_card_id_currency_id;index" json:"currency_id"`
	Balance    decimal.Decimal `gorm:"column:balance" json:"balance"`

	Currency *Currency `gorm:"foreignKey:CurrencyID;references:ID" json:"currency"`
	UserCard *UserCard `gorm:"foreignKey:UserCardID;references:ID" json:"user_card"`
}

func (*UserCardCurrency) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*UserCardCurrency) TableName() string {
	return "user_card_currency"
}

func NewUserCardCurrency(userCardID, currencyID int64, balance decimal.Decimal) *UserCardCurrency {
	return &UserCardCurrency{
		ID:         0,
		UserCardID: userCardID,
		CurrencyID: currencyID,
		Balance:    balance,
	}
}

/**
Preloader
*/

type PreloaderUserCardCurrency struct {
	CurrencyPreload bool
	UserCardPreload bool
}

func NewPreloaderUserCardCurrency() *PreloaderUserCardCurrency {
	return &PreloaderUserCardCurrency{
		CurrencyPreload: false,
		UserCardPreload: false,
	}
}

func (p *PreloaderUserCardCurrency) Preload(tx *gorm.DB) *gorm.DB {
	if p.CurrencyPreload {
		tx = tx.Preload("Currency")
	}
	if p.UserCardPreload {
		tx = tx.Preload("UserCard")
	}
	return tx
}

func (p *PreloaderUserCardCurrency) All() *PreloaderUserCardCurrency {
	p.CurrencyPreload = true
	p.UserCardPreload = true
	return p
}

func (p *PreloaderUserCardCurrency) Currency() *PreloaderUserCardCurrency {
	p.CurrencyPreload = true
	return p
}

func (p *PreloaderUserCardCurrency) UserCard() *PreloaderUserCardCurrency {
	p.UserCardPreload = true
	return p
}
