package model

import (
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/types/role"
)

// User 用户信息
type User struct {
	ID               int64     `gorm:"column:id;primaryKey" json:"id"`
	Username         string    `gorm:"column:username;uniqueIndex;not null" json:"username"`
	Password         string    `gorm:"column:password;not null" json:"-"`
	Nickname         string    `gorm:"column:nickname;index" json:"nickname"`
	Avatar           string    `gorm:"column:avatar;index" json:"avatar"`
	Mail             string    `gorm:"column:mail;index" json:"mail"`
	Phone            string    `gorm:"column:phone;index" json:"phone"`
	MasterCurrencyID int64     `gorm:"column:master_currency_id;default:1" json:"master_currency_id"` // 用户主要货币类型
	Role             role.Type `gorm:"column:role;not null" json:"role"`                              // 用户身份（管理员, 普通用户, 浏览者
	Disabled         int64     `gorm:"column:disabled" json:"disabled"`

	MasterCurrency *Currency          `gorm:"foreignKey:MasterCurrencyID;references:ID" json:"master_currency"`
	AccessToken    []*UserAccessToken `gorm:"foreignKey:UID;references:ID" json:"access_token"` // 访问密钥,用户可以通过密钥访问API
}

func (*User) Alive(tx *gorm.DB) *gorm.DB {
	return tx.Where("disabled = ?", 0)
}

func (*User) TableName() string {
	return "user"
}

func NewUser(username, password, nickname, avatar, mail, phone string, rol role.Type, masterCurrencyID int64) *User {
	return &User{
		ID:               0,
		Username:         username,
		Password:         password,
		Nickname:         nickname,
		Avatar:           avatar,
		Mail:             mail,
		Phone:            phone,
		Role:             rol,
		MasterCurrencyID: masterCurrencyID,
		Disabled:         0,
		MasterCurrency:   nil,
		AccessToken:      nil,
	}
}

/**
Preloader
*/

type PreloaderUser struct {
	MasterCurrencyPreload bool
	AccessTokenPreload    bool
}

func NewPreloaderUser() *PreloaderUser {
	return &PreloaderUser{
		MasterCurrencyPreload: false,
		AccessTokenPreload:    false,
	}
}

func (p *PreloaderUser) Preload(tx *gorm.DB) *gorm.DB {
	if p.MasterCurrencyPreload {
		tx = tx.Preload("MasterCurrency")
	}
	if p.AccessTokenPreload {
		tx = tx.Preload("AccessToken")
	}
	return tx
}

func (p *PreloaderUser) All() *PreloaderUser {
	p.MasterCurrencyPreload = true
	p.AccessTokenPreload = true
	return p
}

func (p *PreloaderUser) MasterCurrency() *PreloaderUser {
	p.MasterCurrencyPreload = true
	return p
}

func (p *PreloaderUser) AccessToken() *PreloaderUser {
	p.AccessTokenPreload = true
	return p
}
