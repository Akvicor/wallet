package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/types/wallet_partition"
)

// UserWalletPartition 用户钱包资金划分
type UserWalletPartition struct {
	ID          int64                        `gorm:"column:id;primaryKey" json:"id"`
	WalletID    int64                        `gorm:"column:wallet_id;uniqueIndex:wallet_id_card_id_currency_id_name_disabled;index" json:"wallet_id"`     // 钱包ID
	CardID      int64                        `gorm:"column:card_id;uniqueIndex:wallet_id_card_id_currency_id_name_disabled;index" json:"card_id"`         // 货币ID
	CurrencyID  int64                        `gorm:"column:currency_id;uniqueIndex:wallet_id_card_id_currency_id_name_disabled;index" json:"currency_id"` // 货币ID
	Name        string                       `gorm:"column:name;uniqueIndex:wallet_id_card_id_currency_id_name_disabled;index" json:"name"`               // 划分名字
	Description string                       `gorm:"column:description" json:"description"`                                                               // 划分描述
	Limit       decimal.Decimal              `gorm:"column:limit" json:"limit"`                                                                           // 划分余额上限
	Balance     decimal.Decimal              `gorm:"column:balance" json:"balance"`                                                                       // 划分余额
	Average     wallet_partition.AverageType `gorm:"column:average" json:"average"`                                                                       // 均分余额
	Sequence    int64                        `gorm:"column:sequence;index" json:"sequence"`                                                               // 划分序号
	Disabled    int64                        `gorm:"column:disabled;uniqueIndex:wallet_id_card_id_currency_id_name_disabled;index" json:"disabled"`

	Wallet   *UserWallet `gorm:"foreignKey:WalletID;references:ID" json:"wallet"`
	Card     *UserCard   `gorm:"foreignKey:CardID;references:ID" json:"card"`
	Currency *Currency   `gorm:"foreignKey:CurrencyID;references:ID" json:"currency"`
}

func (*UserWalletPartition) Alive(tx *gorm.DB) *gorm.DB {
	return tx.Where("disabled = ?", 0)
}

func (u *UserWalletPartition) Pure() *UserWalletPartition {
	u.Currency = nil
	return u
}

func (*UserWalletPartition) TableName() string {
	return "user_wallet_partition"
}

func NewUserWalletPartition(walletID, cardID, currencyID int64, name, description string, limit, balance decimal.Decimal, average wallet_partition.AverageType) *UserWalletPartition {
	return &UserWalletPartition{
		ID:          0,
		WalletID:    walletID,
		CardID:      cardID,
		CurrencyID:  currencyID,
		Name:        name,
		Description: description,
		Limit:       limit,
		Balance:     balance,
		Average:     average,
		Sequence:    0,
		Disabled:    0,
		Currency:    nil,
	}
}

/**
Preloader
*/

type PreloaderUserWalletPartition struct {
	WalletPreload       bool
	CardPreload         bool
	CardCurrencyPreload bool
	CurrencyPreload     bool
}

func NewPreloaderUserWalletPartition() *PreloaderUserWalletPartition {
	return &PreloaderUserWalletPartition{
		WalletPreload:       false,
		CardPreload:         false,
		CardCurrencyPreload: false,
		CurrencyPreload:     false,
	}
}

func (p *PreloaderUserWalletPartition) Preload(tx *gorm.DB) *gorm.DB {
	if p.WalletPreload {
		tx = tx.Preload("Wallet")
	}
	if p.CardPreload {
		tx = tx.Preload("Card")
	}
	if p.CardCurrencyPreload {
		tx = tx.Preload("Card.Currency")
	}
	if p.CurrencyPreload {
		tx = tx.Preload("Currency")
	}
	return tx
}

func (p *PreloaderUserWalletPartition) All() *PreloaderUserWalletPartition {
	p.WalletPreload = true
	p.CardPreload = true
	p.CardCurrencyPreload = true
	p.CurrencyPreload = true
	return p
}

func (p *PreloaderUserWalletPartition) Wallet() *PreloaderUserWalletPartition {
	p.WalletPreload = true
	return p
}

func (p *PreloaderUserWalletPartition) Card() *PreloaderUserWalletPartition {
	p.CardPreload = true
	return p
}

func (p *PreloaderUserWalletPartition) CardCurrency() *PreloaderUserWalletPartition {
	p.CardCurrencyPreload = true
	return p.Card()
}

func (p *PreloaderUserWalletPartition) Currency() *PreloaderUserWalletPartition {
	p.CurrencyPreload = true
	return p
}
