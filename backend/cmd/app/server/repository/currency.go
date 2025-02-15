package repository

import (
	"context"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
)

var Currency = new(currencyRepository)

type currencyRepository struct {
	base[*model.Currency]
}

/**
查找
*/

// FindAll 获取所有的记录, page为nil时不分页
func (u *currencyRepository) FindAll(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderCurrency) (currency []*model.Currency, err error) {
	return u.paging(page, u.preload(c, alive, preloader))
}

// FindAllLike 模糊搜索, page为nil时不分页
func (u *currencyRepository) FindAllLike(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderCurrency, like string) (currency []*model.Currency, err error) {
	tx := u.preload(c, alive, preloader)

	if len(like) > 0 {
		like = "%" + like + "%"
		tx = tx.Where("name LIKE ? OR english_name LIKE ? OR code LIKE ?", like, like, like)
	}

	return u.paging(page, tx)
}

// ExistById 通过ID判断记录是否存在
func (u *currencyRepository) ExistById(c context.Context, alive bool, id []int64) (exist bool, err error) {
	var count uint64
	err = u.preload(c, alive, nil).Where("id IN (?)", id).Select("count(*)").Find(&count).Error
	if err != nil {
		return false, err
	}
	return count == uint64(len(id)), nil
}

// ExistByName 通过名字判断记录是否存在
func (u *currencyRepository) ExistByName(c context.Context, name, englishName string) (exist bool, err error) {
	var count uint64
	err = u.dbm(c).Where("name = ? AND english_name = ?", name, englishName).Select("count(*)").Find(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

/**
创建
*/

// Create 创建记录
func (u *currencyRepository) Create(c context.Context, currency *model.Currency) error {
	return u.WrapResultErr(u.db(c).Create(currency))
}

/**
更新
*/

// Update 更新记录
func (u *currencyRepository) Update(c context.Context, currency *model.Currency) error {
	return u.WrapResultErr(u.db(c).Where("id = ?", currency.ID).Select("*").Updates(currency))
}

/**
删除
*/

// Delete 删除记录
func (u *currencyRepository) Delete(c context.Context, id int64) error {
	return u.WrapResultErr(u.db(c).Where("id = ?", id).Delete(&model.Currency{}))
}
