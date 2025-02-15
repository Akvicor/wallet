package repository

import (
	"context"
	"wallet/cmd/app/server/model"
)

var UserCardRelTypes = new(userCardRelTypesRepository)

type userCardRelTypesRepository struct {
	base[*model.UserCardRelTypes]
}

/**
查找
*/

/**
创建
*/

// Create 创建记录
func (l *userCardRelTypesRepository) Create(c context.Context, types *model.UserCardRelTypes) error {
	return l.WrapResultErr(l.db(c).Create(types))
}

// Creates 批量创建记录
func (l *userCardRelTypesRepository) Creates(c context.Context, types []*model.UserCardRelTypes) error {
	return l.WrapResultErr(l.db(c).Create(types))
}

/**
更新
*/

/**
删除
*/

// DeleteByCardID 通过银行卡ID删除记录
func (l *userCardRelTypesRepository) DeleteByCardID(c context.Context, cid int64) error {
	return l.WrapResultErr(l.db(c).Where("user_card_id = ?", cid).Delete(&model.UserCardRelTypes{}))
}
