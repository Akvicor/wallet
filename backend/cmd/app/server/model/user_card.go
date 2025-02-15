package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserCard 用户银行卡
type UserCard struct {
	ID                  int64               `gorm:"column:id;primaryKey" json:"id"`
	UID                 int64               `gorm:"column:uid;index;not null" json:"uid"`
	BankID              int64               `gorm:"column:bank_id;index;not null" json:"bank_id"`
	Types               []*CardType         `gorm:"many2many:user_card_rel_types;foreignKey:ID;joinForeignKey:user_card_id;references:ID;joinReferences:card_type_id" json:"types"`
	Name                string              `gorm:"column:name;index;not null" json:"name"`
	Description         string              `gorm:"column:description" json:"description"`
	Number              string              `gorm:"column:number;index" json:"number"`
	ExpDate             int64               `gorm:"column:exp_date;index" json:"exp_date"`
	CVV                 string              `gorm:"column:cvv" json:"cvv"`
	StatementClosingDay int64               `gorm:"column:statement_closing_day" json:"statement_closing_day"`
	PaymentDueDay       int64               `gorm:"column:payment_due_day" json:"payment_due_day"`
	Password            string              `gorm:"column:password" json:"password"`
	MasterCurrencyID    int64               `gorm:"column:master_currency_id" json:"master_currency_id"`
	Currency            []*UserCardCurrency `gorm:"foreignKey:UserCardID;references:ID" json:"currency"`
	Limit               decimal.Decimal     `gorm:"column:limit" json:"limit"`
	Fee                 string              `gorm:"column:fee" json:"fee"`
	HideBalance         bool                `gorm:"column:hide_balance" json:"hide_balance"`
	Sequence            int64               `gorm:"column:sequence;index" json:"sequence"`
	Disabled            int64               `gorm:"column:disabled;index" json:"disabled"`

	User           *User     `gorm:"foreignKey:UID;references:ID" json:"user"`
	MasterCurrency *Currency `gorm:"foreignKey:MasterCurrencyID;references:ID" json:"master_currency"`
	Bank           *Bank     `gorm:"foreignKey:BankID;references:ID" json:"bank"`
}

func (*UserCard) Alive(tx *gorm.DB) *gorm.DB {
	return tx.Where("disabled = ?", 0)
}

func (*UserCard) TableName() string {
	return "user_card"
}

func NewUserCard(uid, bankID int64, name, description, number string, expDate int64, cvv string, statementClosingDay, paymentDueDay int64, password string, masterCurrencyID int64, limit decimal.Decimal, fee string, hideBalance bool) *UserCard {
	return &UserCard{
		ID:                  0,
		UID:                 uid,
		BankID:              bankID,
		Types:               nil,
		Name:                name,
		Description:         description,
		Number:              number,
		ExpDate:             expDate,
		CVV:                 cvv,
		StatementClosingDay: statementClosingDay,
		PaymentDueDay:       paymentDueDay,
		Password:            password,
		MasterCurrencyID:    masterCurrencyID,
		Currency:            nil,
		Limit:               limit,
		Fee:                 fee,
		HideBalance:         hideBalance,
		Sequence:            0,
		Disabled:            0,
		Bank:                nil,
	}
}

/**
Preloader
*/

type PreloaderUserCard struct {
	UserPreload             bool
	BankPreload             bool
	TypesPreload            bool
	MasterCurrencyPreload   bool
	CurrencyPreload         bool
	CurrencyCurrencyPreload bool
}

func NewPreloaderUserCard() *PreloaderUserCard {
	return &PreloaderUserCard{
		UserPreload:             false,
		BankPreload:             false,
		TypesPreload:            false,
		MasterCurrencyPreload:   false,
		CurrencyPreload:         false,
		CurrencyCurrencyPreload: false,
	}
}

func (p *PreloaderUserCard) Preload(tx *gorm.DB) *gorm.DB {
	if p.UserPreload {
		tx = tx.Preload("User")
	}
	if p.BankPreload {
		tx = tx.Preload("Bank")
	}
	if p.TypesPreload {
		tx = tx.Preload("Types")
	}
	if p.MasterCurrencyPreload {
		tx = tx.Preload("MasterCurrency")
	}
	if p.CurrencyPreload {
		tx = tx.Preload("Currency")
	}
	if p.CurrencyCurrencyPreload {
		tx = tx.Preload("Currency.Currency")
	}
	return tx
}

func (p *PreloaderUserCard) All() *PreloaderUserCard {
	p.UserPreload = true
	p.BankPreload = true
	p.TypesPreload = true
	p.MasterCurrencyPreload = true
	p.CurrencyPreload = true
	p.CurrencyCurrencyPreload = true
	return p
}

func (p *PreloaderUserCard) User() *PreloaderUserCard {
	p.UserPreload = true
	return p
}

func (p *PreloaderUserCard) Bank() *PreloaderUserCard {
	p.BankPreload = true
	return p
}

func (p *PreloaderUserCard) Types() *PreloaderUserCard {
	p.TypesPreload = true
	return p
}

func (p *PreloaderUserCard) MasterCurrency() *PreloaderUserCard {
	p.MasterCurrencyPreload = true
	return p
}

func (p *PreloaderUserCard) Currency() *PreloaderUserCard {
	p.CurrencyPreload = true
	return p
}

func (p *PreloaderUserCard) CurrencyCurrency() *PreloaderUserCard {
	p.CurrencyCurrencyPreload = true
	return p.Currency()
}
