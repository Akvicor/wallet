package repository

import (
	"context"
	"errors"
	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"time"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
)

var UserPeriodPay = new(userPeriodPayRepository)

type userPeriodPayRepository struct {
	base[*model.UserPeriodPay]
}

/**
查找
*/

// FindAll 获取所有的周期性付款, page为nil时不分页
func (u *userPeriodPayRepository) FindAll(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserPeriodPay) (userPeriodPay []*model.UserPeriodPay, err error) {
	return u.paging(page, u.preload(c, alive, preloader).Order("next_of_period ASC").Order("name ASC"))
}

// FindAllByUID 获取所有的周期性付款, page为nil时不分页
func (u *userPeriodPayRepository) FindAllByUID(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserPeriodPay, uid int64) (userPeriodPay []*model.UserPeriodPay, err error) {
	return u.paging(page, u.preload(c, alive, preloader).Where("uid = ?", uid).Order("next_of_period ASC").Order("name ASC"))
}

// FindAllByUIDLike 模糊搜索, page为nil时不分页
func (u *userPeriodPayRepository) FindAllByUIDLike(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserPeriodPay, uid int64, like string) (userPeriodPay []*model.UserPeriodPay, err error) {
	tx := u.preload(c, alive, preloader).Where("uid = ?", uid)

	if len(like) > 0 {
		like = "%" + like + "%"
		tx = tx.Where("name LIKE ? OR description LIKE ?", like, like)
	}

	return u.paging(page, tx.Order("next_of_period ASC").Order("name ASC"))
}

// FindByUID 获取所有的记录, page为nil时不分页
func (u *userPeriodPayRepository) FindByUID(c context.Context, alive bool, preloader *model.PreloaderUserPeriodPay, uid, id int64) (userPeriodPay *model.UserPeriodPay, err error) {
	userPeriodPay = new(model.UserPeriodPay)
	err = u.WrapResultErr(u.preload(c, alive, preloader).Where("uid = ? AND id = ?", uid, id).First(userPeriodPay))
	return
}

// DetectValidByUID 判断PeriodPay是否属于某个用户
func (u *userPeriodPayRepository) DetectValidByUID(c context.Context, uid int64, ids []int64) (err error) {
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

// FindAllByPeriodByDayLimit 获取Limit天内所有的周期性付款
func (u *userPeriodPayRepository) FindAllByPeriodByDayLimit(c context.Context, alive bool, preloader *model.PreloaderUserPeriodPay, dayLimit int) (userPeriodPay []*model.UserPeriodPay, err error) {
	nextAt := now.With(time.Now().AddDate(0, 0, dayLimit)).BeginningOfDay()
	next := nextAt.Unix()
	userPeriodPay = make([]*model.UserPeriodPay, 0)
	err = u.preload(c, alive, preloader).Where("start_at <= ? AND expiration_date >= ? AND next_of_period <= ?", next, next, next).Find(&userPeriodPay).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return userPeriodPay, err
	}
	return userPeriodPay, nil
}

/**
创建
*/

// Create 创建记录
func (u *userPeriodPayRepository) Create(c context.Context, userPeriodPay *model.UserPeriodPay) error {
	return u.WrapResultErr(u.db(c).Create(userPeriodPay))
}

/**
更新
*/

// UpdateByUID 更新记录
func (u *userPeriodPayRepository) UpdateByUID(c context.Context, alive bool, uid int64, userPeriodPay *model.UserPeriodPay) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, userPeriodPay.ID).Select("*").Updates(userPeriodPay))
}

// UpdateNextPeriod 更新记录
func (u *userPeriodPayRepository) UpdateNextPeriod(c context.Context, pay *model.UserPeriodPay) error {
	return u.WrapResultErr(u.dbm(c).Where("id = ?", pay.ID).Updates(map[string]any{
		"next_of_period":   pay.NextOfPeriod,
		"expiration_times": pay.ExpirationTimes,
	}))
}

/**
删除
*/

// DeleteByUID 删除记录
func (u *userPeriodPayRepository) DeleteByUID(c context.Context, uid, id int64) error {
	return u.WrapResultErr(u.db(c).Where("uid = ? AND id = ?", uid, id).Delete(&model.UserPeriodPay{}))
}
