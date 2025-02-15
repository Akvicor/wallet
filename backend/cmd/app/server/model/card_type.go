package model

import "gorm.io/gorm"

// CardType 银行卡种类(Visa/Mastercard)
type CardType struct {
	ID   int64  `gorm:"column:id;primaryKey;uniqueIndex" json:"id"`
	Name string `gorm:"column:name;uniqueIndex;not null" json:"name"`
}

func (*CardType) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*CardType) TableName() string {
	return "card_type"
}

func NewCardType(name string) *CardType {
	return &CardType{
		ID:   0,
		Name: name,
	}
}

/**
Preloader
*/

type PreloaderCardType struct{}

func NewPreloaderCardType() *PreloaderCardType {
	return &PreloaderCardType{}
}

func (p *PreloaderCardType) Preload(tx *gorm.DB) *gorm.DB {
	return tx
}

func (p *PreloaderCardType) All() *PreloaderCardType {
	return p
}
