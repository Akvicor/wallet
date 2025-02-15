package model

import "gorm.io/gorm"

// UserCardRelTypes 用户银行卡与银行卡种类的多对多关系
type UserCardRelTypes struct {
	UserCardID int64 `gorm:"column:user_card_id;primaryKey" json:"user_card_id"`
	CardTypeID int64 `gorm:"column:card_type_id;primaryKey" json:"card_type_id"`
}

func (*UserCardRelTypes) Alive(tx *gorm.DB) *gorm.DB {
	return tx
}

func (*UserCardRelTypes) TableName() string {
	return "user_card_rel_types"
}

func NewUserCardRelTypes(cid, tid int64) *UserCardRelTypes {
	return &UserCardRelTypes{
		UserCardID: cid,
		CardTypeID: tid,
	}
}

/**
Preloader
*/

type PreloaderUserCardRelTypes struct {
}

func NewPreloaderUserCardRelTypes() *PreloaderUserCardRelTypes {
	return &PreloaderUserCardRelTypes{}
}

func (p *PreloaderUserCardRelTypes) Preload(tx *gorm.DB) *gorm.DB {
	return tx
}

func (p *PreloaderUserCardRelTypes) All() *PreloaderUserCardRelTypes {
	return p
}
