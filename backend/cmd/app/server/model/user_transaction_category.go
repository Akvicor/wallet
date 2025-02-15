package model

import (
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/types/transaction"
)

// UserTransactionCategory 交易类型
type UserTransactionCategory struct {
	ID       int64            `gorm:"column:id;primaryKey" json:"id"`
	UID      int64            `gorm:"column:uid;index;not null" json:"uid"`
	Type     transaction.Type `gorm:"column:type;index;not null" json:"type"`
	Name     string           `gorm:"column:name;not null" json:"name"`
	Colour   string           `gorm:"column:colour;not null" json:"colour"`
	Sequence int64            `gorm:"column:sequence;index" json:"sequence"`
}

func (*UserTransactionCategory) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*UserTransactionCategory) TableName() string {
	return "user_transaction_category"
}

func NewUserTransactionCategory(uid int64, transactionType transaction.Type, name, colour string) *UserTransactionCategory {
	return &UserTransactionCategory{
		ID:       0,
		UID:      uid,
		Type:     transactionType,
		Name:     name,
		Sequence: 0,
		Colour:   colour,
	}
}

/**
Preloader
*/

type PreloaderUserTransactionCategory struct{}

func NewPreloaderUserTransactionCategory() *PreloaderUserTransactionCategory {
	return &PreloaderUserTransactionCategory{}
}

func (p *PreloaderUserTransactionCategory) Preload(tx *gorm.DB) *gorm.DB {
	return tx
}

func (p *PreloaderUserTransactionCategory) All() *PreloaderUserTransactionCategory {
	return p
}
