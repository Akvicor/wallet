package model

import (
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/types/wallet"
)

// UserWallet 钱包
type UserWallet struct {
	ID          int64       `gorm:"column:id;primaryKey" json:"id"`
	UID         int64       `gorm:"column:uid;index;not null" json:"uid"`
	Name        string      `gorm:"column:name;index;not null" json:"name"`
	Description string      `gorm:"column:description" json:"description"`
	WalletType  wallet.Type `gorm:"column:wallet_type;index;not null" json:"wallet_type"`
	Sequence    int64       `gorm:"column:sequence;index" json:"sequence"`
	Disabled    int64       `gorm:"column:disabled" json:"disabled"`

	Partition []*UserWalletPartition `gorm:"foreignKey:WalletID;references:ID" json:"partition"`
}

func (*UserWallet) Alive(tx *gorm.DB) *gorm.DB {
	return tx.Where("disabled = ?", 0)
}

func (*UserWallet) TableName() string {
	return "user_wallet"
}

func NewUserWallet(uid int64, name, description string, walletType wallet.Type) *UserWallet {
	return &UserWallet{
		ID:          0,
		UID:         uid,
		Name:        name,
		Description: description,
		WalletType:  walletType,
		Sequence:    0,
		Disabled:    0,
		Partition:   nil,
	}
}

/**
Preloader
*/

type PreloaderUserWallet struct {
	PartitionPreload         bool
	PartitionPreloadAlive    bool
	PartitionCardPreload     bool
	PartitionCurrencyPreload bool
}

func NewPreloaderUserWallet() *PreloaderUserWallet {
	return &PreloaderUserWallet{
		PartitionPreload:         false,
		PartitionPreloadAlive:    false,
		PartitionCardPreload:     false,
		PartitionCurrencyPreload: false,
	}
}

func (p *PreloaderUserWallet) Preload(tx *gorm.DB) *gorm.DB {
	if p.PartitionPreload {
		if p.PartitionPreloadAlive {
			var part *UserWalletPartition
			tx = tx.Preload("Partition", func(tx *gorm.DB) *gorm.DB { return part.Alive(tx).Order("disabled ASC").Order("sequence ASC") })
		} else {
			tx = tx.Preload("Partition", func(tx *gorm.DB) *gorm.DB { return tx.Order("disabled ASC").Order("sequence ASC") })
		}
		if p.PartitionCardPreload {
			tx = tx.Preload("Partition.Card")
		}
		if p.PartitionCurrencyPreload {
			tx = tx.Preload("Partition.Currency")
		}
	}
	return tx
}

func (p *PreloaderUserWallet) All(partitionAlive bool) *PreloaderUserWallet {
	p.PartitionPreload = true
	p.PartitionPreloadAlive = partitionAlive
	p.PartitionCardPreload = true
	p.PartitionCurrencyPreload = true
	return p
}

func (p *PreloaderUserWallet) Partition(alive bool) *PreloaderUserWallet {
	p.PartitionPreload = true
	p.PartitionPreloadAlive = alive
	return p
}

func (p *PreloaderUserWallet) PartitionCard(partitionAlive bool) *PreloaderUserWallet {
	p.PartitionCardPreload = true
	return p.Partition(partitionAlive)
}

func (p *PreloaderUserWallet) PartitionCurrency(partitionAlive bool) *PreloaderUserWallet {
	p.PartitionCurrencyPreload = true
	return p.Partition(partitionAlive)
}
