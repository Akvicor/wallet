package service

import (
	"context"
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var UserExchangeRate = new(userExchangeRateService)

type userExchangeRateService struct {
	base
}

// FindAllByUID 获取用户的所有货币汇率
func (u *userExchangeRateService) FindAllByUID(page *resp.PageModel, alive bool, preload *model.PreloaderUserExchangeRate, uid int64) (userExchangeRate []*model.UserExchangeRate, err error) {
	return repository.UserExchangeRate.FindAllByUID(context.Background(), page, alive, preload, uid)
}

// FindByUID 获取用户的所有货币汇率
func (u *userExchangeRateService) FindByUID(alive bool, preload *model.PreloaderUserExchangeRate, uid, from, to int64) (userExchangeRate *model.UserExchangeRate, err error) {
	return repository.UserExchangeRate.FindByUID(context.Background(), alive, preload, uid, from, to)
}

// Create 创建用户的货币汇率
func (u *userExchangeRateService) Create(uid, fromCurrencyID, toCurrencyID int64, rate decimal.Decimal) (userExchangeRate *model.UserExchangeRate, err error) {
	userExchangeRate = model.NewUserExchangeRate(uid, fromCurrencyID, toCurrencyID, rate)
	err = repository.UserExchangeRate.Create(context.Background(), userExchangeRate)
	if err != nil {
		return nil, err
	}
	return userExchangeRate, nil
}

// Update 更新用户的货币汇率
func (u *userExchangeRateService) Update(uid, id, fromCurrencyID, toCurrencyID int64, rate decimal.Decimal) (err error) {
	userExchangeRate := model.NewUserExchangeRate(uid, fromCurrencyID, toCurrencyID, rate)
	userExchangeRate.ID = id
	err = repository.UserExchangeRate.Update(context.Background(), userExchangeRate)
	if err != nil {
		return err
	}
	return nil
}

// DeleteWithUID 通过UID删除货币汇率
func (u *userExchangeRateService) DeleteWithUID(alive bool, uid, id int64) error {
	return repository.UserExchangeRate.DeleteWithUID(context.Background(), alive, uid, id)
}

// DeleteByUID 通过UID删除货币汇率
func (u *userExchangeRateService) DeleteByUID(alive bool, uid int64) error {
	return repository.UserExchangeRate.DeleteByUID(context.Background(), alive, uid)
}

// DeleteByCurrencyID 通过CurrencyID删除货币汇率
func (u *userExchangeRateService) DeleteByCurrencyID(alive bool, fromCurrencyID, toCurrencyID int64) error {
	return repository.UserExchangeRate.DeleteByCurrencyID(context.Background(), alive, fromCurrencyID, toCurrencyID)
}
