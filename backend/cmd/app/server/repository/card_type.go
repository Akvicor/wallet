package repository

import (
	"context"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
)

var CardType = new(cardTypeRepository)

type cardTypeRepository struct {
	base[*model.CardType]
}

/**
查找
*/

// FindAll 获取所有的记录, page为nil时不分页
func (u *cardTypeRepository) FindAll(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderCardType) (cardTypes []*model.CardType, err error) {
	return u.paging(page, u.preload(c, alive, preloader))
}

// FindAllLike 模糊搜索, page为nil时不分页
func (u *cardTypeRepository) FindAllLike(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderCardType, like string) (cardTypes []*model.CardType, err error) {
	tx := u.preload(c, alive, preloader)

	if len(like) > 0 {
		like = "%" + like + "%"
		tx = tx.Where("name LIKE ?", like)
	}

	return u.paging(page, tx)
}

// ExistById 通过ID判断记录是否存在
func (u *cardTypeRepository) ExistById(c context.Context, alive bool, id []int64) (exist bool, err error) {
	var count uint64
	err = u.preload(c, alive, nil).Where("id IN (?)", id).Select("count(*)").Find(&count).Error
	if err != nil {
		return false, err
	}
	return count == uint64(len(id)), nil
}

// ExistByName 通过名字判断记录是否存在
func (u *cardTypeRepository) ExistByName(c context.Context, name string) (exist bool, err error) {
	var count uint64
	err = u.dbm(c).Where("name = ?", name).Select("count(*)").Find(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

/**
创建
*/

// Create 创建记录
func (u *cardTypeRepository) Create(c context.Context, cardType *model.CardType) error {
	return u.WrapResultErr(u.db(c).Create(cardType))
}

/**
更新
*/

// Update 更新记录
func (u *cardTypeRepository) Update(c context.Context, cardType *model.CardType) error {
	return u.WrapResultErr(u.db(c).Where("id = ?", cardType.ID).Select("*").Updates(cardType))
}

/**
删除
*/

// Delete 删除记录
func (u *cardTypeRepository) Delete(c context.Context, id int64) error {
	return u.WrapResultErr(u.db(c).Where("id = ?", id).Delete(&model.CardType{}))
}
