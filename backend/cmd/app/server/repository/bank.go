package repository

import (
	"context"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
)

var Bank = new(bankRepository)

type bankRepository struct {
	base[*model.Bank]
}

/**
查找
*/

// FindAll 获取所有的记录, page为nil时不分页
func (u *bankRepository) FindAll(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderBank) (banks []*model.Bank, err error) {
	return u.paging(page, u.preload(c, alive, preloader))
}

// FindAllLike 模糊搜索, page为nil时不分页
func (u *bankRepository) FindAllLike(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderBank, like string) (banks []*model.Bank, err error) {
	tx := u.preload(c, alive, preloader)

	if len(like) > 0 {
		like = "%" + like + "%"
		tx = tx.Where("name LIKE ? OR english_name LIKE ? OR abbr LIKE ?", like, like, like)
	}

	return u.paging(page, tx)
}

// ExistById 通过ID判断记录是否存在
func (u *bankRepository) ExistById(c context.Context, alive bool, id []int64) (exist bool, err error) {
	var count uint64
	err = u.preload(c, alive, nil).Where("id IN (?)", id).Select("count(*)").Find(&count).Error
	if err != nil {
		return false, err
	}
	return count == uint64(len(id)), nil
}

// ExistByName 通过名字判断记录是否存在
func (u *bankRepository) ExistByName(c context.Context, name, englishName string) (exist bool, err error) {
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
func (u *bankRepository) Create(c context.Context, bank *model.Bank) error {
	return u.WrapResultErr(u.db(c).Create(bank))
}

/**
更新
*/

// Update 更新记录
func (u *bankRepository) Update(c context.Context, bank *model.Bank) error {
	return u.WrapResultErr(u.db(c).Where("id = ?", bank.ID).Select("*").Updates(bank))
}

/**
删除
*/

// Delete 删除记录
func (u *bankRepository) Delete(c context.Context, id int64) error {
	return u.WrapResultErr(u.db(c).Where("id = ?", id).Delete(&model.Bank{}))
}
