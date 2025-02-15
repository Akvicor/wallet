package service

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/period"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var UserPeriodPay = new(userPeriodPayService)

type userPeriodPayService struct {
	base
}

// FindAll 获取所有的周期性付款, page为nil时不分页
func (w *userPeriodPayService) FindAll(page *resp.PageModel, alive bool, preloader *model.PreloaderUserPeriodPay) ([]*model.UserPeriodPay, error) {
	return repository.UserPeriodPay.FindAll(context.Background(), page, alive, preloader)
}

// FindAllByUID 获取所有的周期性付款, page为nil时不分页
func (w *userPeriodPayService) FindAllByUID(page *resp.PageModel, alive bool, preloader *model.PreloaderUserPeriodPay, uid int64) ([]*model.UserPeriodPay, error) {
	return repository.UserPeriodPay.FindAllByUID(context.Background(), page, alive, preloader, uid)
}

// FindAllByUIDLike 获取所有的周期性付款, page为nil时不分页
func (w *userPeriodPayService) FindAllByUIDLike(page *resp.PageModel, alive bool, preloader *model.PreloaderUserPeriodPay, uid int64, like string) ([]*model.UserPeriodPay, error) {
	return repository.UserPeriodPay.FindAllByUIDLike(context.Background(), page, alive, preloader, uid, like)
}

// FindByUID 获取周期性付款
func (w *userPeriodPayService) FindByUID(alive bool, preloader *model.PreloaderUserPeriodPay, uid, id int64) (*model.UserPeriodPay, error) {
	return repository.UserPeriodPay.FindByUID(context.Background(), alive, preloader, uid, id)
}

// FindAllByPeriodByDayLimit 获取DayLimit天内周期性付款
func (w *userPeriodPayService) FindAllByPeriodByDayLimit(alive bool, preloader *model.PreloaderUserPeriodPay, dayLimit int) ([]*model.UserPeriodPay, error) {
	return repository.UserPeriodPay.FindAllByPeriodByDayLimit(context.Background(), alive, preloader, dayLimit)
}

// Create 创建周期性付款
func (w *userPeriodPayService) Create(uid int64, name, description string, currencyID int64, value decimal.Decimal, periodType period.Type, startAt, nextOfPeriod, intervalOfPeriod, expirationDate, expirationTimes int64) (userPeriodPay *model.UserPeriodPay, err error) {
	// 判断周期类型是否合法
	if !periodType.Valid() {
		return nil, errors.New("period type is invalid")
	}
	startAt, nextOfPeriod, intervalOfPeriod = periodType.Calculate(startAt, nextOfPeriod, intervalOfPeriod)

	userPeriodPay = model.NewUserPeriodPay(uid, name, description, currencyID, value, periodType, startAt, nextOfPeriod, intervalOfPeriod, expirationDate, expirationTimes)

	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断货币是否存在
		exist, e = repository.Currency.ExistById(ctx, true, []int64{currencyID})
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("货币种类不存在")
		}
		err = repository.UserPeriodPay.Create(ctx, userPeriodPay)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return userPeriodPay, nil
}

// Update 更新周期性付款
func (w *userPeriodPayService) Update(uid, id int64, name, description string, currencyID int64, value decimal.Decimal, periodType period.Type, startAt, nextOfPeriod, intervalOfPeriod, expirationDate, expirationTimes int64) (err error) {
	// 判断周期类型是否合法
	if !periodType.Valid() {
		return errors.New("period type is invalid")
	}
	startAt, nextOfPeriod, intervalOfPeriod = periodType.Calculate(startAt, nextOfPeriod, intervalOfPeriod)

	userPeriodPay := model.NewUserPeriodPay(uid, name, description, currencyID, value, periodType, startAt, nextOfPeriod, intervalOfPeriod, expirationDate, expirationTimes)
	userPeriodPay.ID = id

	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断货币是否存在
		exist, e = repository.Currency.ExistById(ctx, true, []int64{currencyID})
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("货币种类不存在")
		}
		// 更新周期性付款
		err = repository.UserPeriodPay.UpdateByUID(ctx, false, uid, userPeriodPay)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdatesNextPeriodByID 更新周期性付款的过期次数
func (w *userPeriodPayService) UpdatesNextPeriodByID(uid, id int64) (err error) {
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		pay, err := repository.UserPeriodPay.FindByUID(ctx, false, nil, uid, id)
		if err != nil {
			return err
		}
		pay.NextPeriod()
		err = repository.UserPeriodPay.UpdateNextPeriod(ctx, pay)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdatesNextPeriod 更新周期性付款的过期次数
func (w *userPeriodPayService) UpdatesNextPeriod(pays []*model.UserPeriodPay) (err error) {
	for _, pay := range pays {
		pay.NextPeriod()
	}
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		for _, pay := range pays {
			err = repository.UserPeriodPay.UpdateNextPeriod(ctx, pay)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// DeleteByUID 删除周期性付款
func (w *userPeriodPayService) DeleteByUID(uid, id int64) error {
	return repository.UserPeriodPay.DeleteByUID(context.Background(), uid, id)
}
