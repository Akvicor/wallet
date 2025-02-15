package model

import "gorm.io/gorm"

// Bank 银行信息
type Bank struct {
	ID          int64  `gorm:"column:id;primaryKey" json:"id"`
	Name        string `gorm:"column:name;uniqueIndex;not null" json:"name"`
	EnglishName string `gorm:"column:english_name;uniqueIndex;not null" json:"english_name"`
	Abbr        string `gorm:"column:abbr;not null" json:"abbr"`
	Phone       string `gorm:"column:phone;not null" json:"phone"`
}

func (*Bank) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*Bank) TableName() string {
	return "bank"
}

func NewBank(name, englishName, abbr, phone string) *Bank {
	return &Bank{
		ID:          0,
		Name:        name,
		EnglishName: englishName,
		Abbr:        abbr,
		Phone:       phone,
	}
}

/**
Preloader
*/

type PreloaderBank struct{}

func NewPreloaderBank() *PreloaderBank {
	return &PreloaderBank{}
}

func (p *PreloaderBank) Preload(tx *gorm.DB) *gorm.DB {
	return tx
}

func (p *PreloaderBank) All() *PreloaderBank {
	return p
}
