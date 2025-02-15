package model

import "gorm.io/gorm"

// Currency 货币类型
type Currency struct {
	ID          int64  `gorm:"column:id;primaryKey" json:"id"`
	Name        string `gorm:"column:name;uniqueIndex;not null" json:"name"`
	EnglishName string `gorm:"column:english_name;uniqueIndex;not null" json:"english_name"`
	Code        string `gorm:"column:code;uniqueIndex;not null" json:"code"`
	Symbol      string `gorm:"column:symbol;not null" json:"symbol"`
}

func (*Currency) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*Currency) TableName() string {
	return "currency"
}

func NewCurrency(name, englishName, code, symbol string) *Currency {
	return &Currency{
		ID:          0,
		Name:        name,
		EnglishName: englishName,
		Code:        code,
		Symbol:      symbol,
	}
}

/**
Preloader
*/

type PreloaderCurrency struct{}

func NewPreloaderCurrency() *PreloaderCurrency {
	return &PreloaderCurrency{}
}

func (p *PreloaderCurrency) Preload(tx *gorm.DB) *gorm.DB {
	return tx
}

func (p *PreloaderCurrency) All() *PreloaderCurrency {
	return p
}
