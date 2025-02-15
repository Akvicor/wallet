package repository

import (
	"context"
	"errors"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
)

var UserExchangeRate = new(userExchangeRateRepository)

type userExchangeRateRepository struct {
	base[*model.UserExchangeRate]
}

/**
查找
*/

// FindAllByUID 获取用户所有的记录, page为nil时不分页
func (u *userExchangeRateRepository) FindAllByUID(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserExchangeRate, uid int64) (userExchangeRate []*model.UserExchangeRate, err error) {
	return u.paging(page, u.preload(c, alive, preloader).Where("uid = ?", uid).Order("from_currency_id ASC").Order("to_currency_id ASC"))
}

// FindByUID 获取用户的货币兑换汇率记录
func (u *userExchangeRateRepository) FindByUID(c context.Context, alive bool, preloader *model.PreloaderUserExchangeRate, uid, from, to int64) (userExchangeRate *model.UserExchangeRate, err error) {
	userExchangeRate = new(model.UserExchangeRate)
	err = u.preload(c, alive, preloader).Where("uid = ? AND from_currency_id = ? AND to_currency_id = ?", uid, from, to).First(userExchangeRate).Error
	return
}

// DetectValidByUID 判断UserExchangeRate是否属于某个用户
func (u *userExchangeRateRepository) DetectValidByUID(c context.Context, uid int64, ids []int64) (err error) {
	count := int64(0)
	err = u.dbt(c).Where("uid = ? AND id IN (?)", uid, ids).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("inconsistent quantity")
	}
	return nil
}

/**
创建
*/

// Create 创建记录
func (u *userExchangeRateRepository) Create(c context.Context, currency *model.UserExchangeRate) error {
	return u.WrapResultErr(u.db(c).Create(currency))
}

/**
更新
*/

// Update 更新记录
func (u *userExchangeRateRepository) Update(c context.Context, currency *model.UserExchangeRate) error {
	return u.WrapResultErr(u.db(c).Where("id = ?", currency.ID).Select("*").Updates(currency))
}

/**
删除
*/

// DeleteWithUID 通过ID和UID删除记录
func (u *userExchangeRateRepository) DeleteWithUID(c context.Context, alive bool, uid, id int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, id).Delete(&model.UserExchangeRate{}))
}

// DeleteByUID 通过UID删除记录
func (u *userExchangeRateRepository) DeleteByUID(c context.Context, alive bool, uid int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ?", uid).Delete(&model.UserExchangeRate{}))
}

// DeleteByCurrencyID 通过CurrencyID删除记录
func (u *userExchangeRateRepository) DeleteByCurrencyID(c context.Context, alive bool, fromCurrencyID, toCurrencyID int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("from_currency_id = ? OR to_currency_id = ?", fromCurrencyID, toCurrencyID).Delete(&model.UserExchangeRate{}))
}
