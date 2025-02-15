package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/types/transaction"
)

// UserTransaction 交易
type UserTransaction struct {
	ID              int64            `gorm:"column:id;primaryKey" json:"id"`
	UID             int64            `gorm:"column:uid;index" json:"uid"`                             // 用户ID
	FromPartitionID int64            `gorm:"column:from_partition_id;index" json:"from_partition_id"` // 源划分
	FromCurrencyID  int64            `gorm:"column:from_currency_id;index" json:"from_currency_id"`   // 源货币
	ToPartitionID   int64            `gorm:"column:to_partition_id;index" json:"to_partition_id"`     // 目标划分
	ToCurrencyID    int64            `gorm:"column:to_currency_id;index" json:"to_currency_id"`       // 目标货币
	CategoryID      int64            `gorm:"column:category_id;index" json:"category_id"`             // 分类ID
	Type            transaction.Type `gorm:"column:type;index" json:"type"`                           // 交易类型用户ID
	Description     string           `gorm:"column:description;index" json:"description"`             // 交易描述用户ID
	FromValue       decimal.Decimal  `gorm:"column:from_value" json:"from_value"`                     // 源金额
	ToValue         decimal.Decimal  `gorm:"column:to_value" json:"to_value"`                         // 目标金额
	Fee             decimal.Decimal  `gorm:"column:fee" json:"fee"`                                   // 手续费(源货币类型
	Created         int64            `gorm:"column:created;index" json:"created"`                     // 交易发生时间
	Checked         int64            `gorm:"column:checked;index" json:"checked"`                     // 交易确认时间

	FromPartition *UserWalletPartition     `gorm:"foreignKey:FromPartitionID;references:ID" json:"from_partition"`
	FromCurrency  *Currency                `gorm:"foreignKey:FromCurrencyID;references:ID" json:"from_currency"`
	ToPartition   *UserWalletPartition     `gorm:"foreignKey:ToPartitionID;references:ID" json:"to_partition"`
	ToCurrency    *Currency                `gorm:"foreignKey:ToCurrencyID;references:ID" json:"to_currency"`
	Category      *UserTransactionCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
}

func (*UserTransaction) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*UserTransaction) TableName() string {
	return "user_transaction"
}

func NewUserTransaction(uid, fromPartitionID, fromCurrencyID, toPartitionID, toCurrencyID, categoryID int64, transactionType transaction.Type, description string, fromValue, toValue, fee decimal.Decimal, created int64) *UserTransaction {
	return &UserTransaction{
		ID:              0,
		UID:             uid,
		FromPartitionID: fromPartitionID,
		FromCurrencyID:  fromCurrencyID,
		ToPartitionID:   toPartitionID,
		ToCurrencyID:    toCurrencyID,
		CategoryID:      categoryID,
		Type:            transactionType,
		Description:     description,
		FromValue:       fromValue,
		ToValue:         toValue,
		Fee:             fee,
		Created:         created,
		Checked:         0,
		FromPartition:   nil,
		ToPartition:     nil,
		Category:        nil,
	}
}

/**
Preloader
*/

type PreloaderUserTransaction struct {
	FromPartitionPreload         bool
	FromPartitionWalletPreload   bool
	FromPartitionCardPreload     bool
	FromPartitionCurrencyPreload bool
	FromCurrencyPreload          bool
	ToPartitionPreload           bool
	ToPartitionWalletPreload     bool
	ToPartitionCardPreload       bool
	ToPartitionCurrencyPreload   bool
	ToCurrencyPreload            bool
	CategoryPreload              bool
}

func NewPreloaderUserTransaction() *PreloaderUserTransaction {
	return &PreloaderUserTransaction{
		FromPartitionPreload:         false,
		FromPartitionWalletPreload:   false,
		FromPartitionCardPreload:     false,
		FromPartitionCurrencyPreload: false,
		FromCurrencyPreload:          false,
		ToPartitionPreload:           false,
		ToPartitionWalletPreload:     false,
		ToPartitionCardPreload:       false,
		ToPartitionCurrencyPreload:   false,
		ToCurrencyPreload:            false,
		CategoryPreload:              false,
	}
}

func (p *PreloaderUserTransaction) Preload(tx *gorm.DB) *gorm.DB {
	if p.FromPartitionPreload {
		tx = tx.Preload("FromPartition")
	}
	if p.FromPartitionWalletPreload {
		tx = tx.Preload("FromPartition.Wallet")
	}
	if p.FromPartitionCardPreload {
		tx = tx.Preload("FromPartition.Card")
	}
	if p.FromPartitionCurrencyPreload {
		tx = tx.Preload("FromPartition.Currency")
	}
	if p.FromCurrencyPreload {
		tx = tx.Preload("FromCurrency")
	}
	if p.ToPartitionPreload {
		tx = tx.Preload("ToPartition")
	}
	if p.ToPartitionWalletPreload {
		tx = tx.Preload("ToPartition.Wallet")
	}
	if p.ToPartitionCardPreload {
		tx = tx.Preload("ToPartition.Card")
	}
	if p.ToPartitionCurrencyPreload {
		tx = tx.Preload("ToPartition.Currency")
	}
	if p.ToCurrencyPreload {
		tx = tx.Preload("ToCurrency")
	}
	if p.CategoryPreload {
		tx = tx.Preload("Category")
	}
	return tx
}

func (p *PreloaderUserTransaction) All() *PreloaderUserTransaction {
	p.FromPartitionPreload = true
	p.FromPartitionWalletPreload = true
	p.FromPartitionCardPreload = true
	p.FromPartitionCurrencyPreload = true
	p.FromCurrencyPreload = true
	p.ToPartitionPreload = true
	p.ToPartitionWalletPreload = true
	p.ToPartitionCardPreload = true
	p.ToPartitionCurrencyPreload = true
	p.ToCurrencyPreload = true
	p.CategoryPreload = true
	return p
}

func (p *PreloaderUserTransaction) FromPartition() *PreloaderUserTransaction {
	p.FromPartitionPreload = true
	return p
}

func (p *PreloaderUserTransaction) FromPartitionWallet() *PreloaderUserTransaction {
	p.FromPartitionWalletPreload = true
	return p.FromPartition()
}

func (p *PreloaderUserTransaction) FromPartitionCard() *PreloaderUserTransaction {
	p.FromPartitionCardPreload = true
	return p.FromPartition()
}

func (p *PreloaderUserTransaction) FromPartitionCurrency() *PreloaderUserTransaction {
	p.FromPartitionCurrencyPreload = true
	return p.FromPartition()
}

func (p *PreloaderUserTransaction) FromCurrency() *PreloaderUserTransaction {
	p.FromCurrencyPreload = true
	return p
}

func (p *PreloaderUserTransaction) ToPartition() *PreloaderUserTransaction {
	p.ToPartitionPreload = true
	return p
}

func (p *PreloaderUserTransaction) ToPartitionWallet() *PreloaderUserTransaction {
	p.ToPartitionWalletPreload = true
	return p.ToPartition()
}

func (p *PreloaderUserTransaction) ToPartitionCard() *PreloaderUserTransaction {
	p.ToPartitionCardPreload = true
	return p.ToPartition()
}

func (p *PreloaderUserTransaction) ToPartitionCurrency() *PreloaderUserTransaction {
	p.ToPartitionCurrencyPreload = true
	return p.ToPartition()
}

func (p *PreloaderUserTransaction) ToCurrency() *PreloaderUserTransaction {
	p.ToCurrencyPreload = true
	return p
}

func (p *PreloaderUserTransaction) Category() *PreloaderUserTransaction {
	p.CategoryPreload = true
	return p
}
