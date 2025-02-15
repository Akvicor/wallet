package service

import (
	"context"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var Currency = new(currencyService)

type currencyService struct {
	base
}

// FindAll 获取所有货币类型
func (u *currencyService) FindAll(page *resp.PageModel, alive bool, preload *model.PreloaderCurrency) (currency []*model.Currency, err error) {
	return repository.Currency.FindAll(context.Background(), page, alive, preload)
}

// FindAllLike 获取所有货币类型
func (u *currencyService) FindAllLike(page *resp.PageModel, alive bool, preload *model.PreloaderCurrency, like string) (currency []*model.Currency, err error) {
	return repository.Currency.FindAllLike(context.Background(), page, alive, preload, like)
}

// Create 创建货币类型
func (u *currencyService) Create(name, englishName, code, symbol string) (currency *model.Currency, err error) {
	currency = model.NewCurrency(name, englishName, code, symbol)
	err = repository.Currency.Create(context.Background(), currency)
	if err != nil {
		return nil, err
	}
	return currency, nil
}

// Exist 判断货币类型是否存在
func (u *currencyService) Exist(name, englishName string) (exist bool, err error) {
	return repository.Currency.ExistByName(context.Background(), name, englishName)
}

// Update 更新货币类型
func (u *currencyService) Update(id int64, name, englishName, code, symbol string) error {
	currency := model.NewCurrency(name, englishName, code, symbol)
	currency.ID = id
	return repository.Currency.Update(context.Background(), currency)
}

// Delete 删除货币类型
func (u *currencyService) Delete(id int64) error {
	return repository.Currency.Delete(context.Background(), id)
}
