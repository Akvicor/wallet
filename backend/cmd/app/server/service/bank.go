package service

import (
	"context"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var Bank = new(bankService)

type bankService struct {
	base
}

// FindAll 获取所有银行
func (u *bankService) FindAll(page *resp.PageModel, alive bool, preload *model.PreloaderBank) (banks []*model.Bank, err error) {
	return repository.Bank.FindAll(context.Background(), page, alive, preload)
}

// FindAllLike 获取所有银行
func (u *bankService) FindAllLike(page *resp.PageModel, alive bool, preload *model.PreloaderBank, like string) (banks []*model.Bank, err error) {
	return repository.Bank.FindAllLike(context.Background(), page, alive, preload, like)
}

// Create 创建银行
func (u *bankService) Create(name, englishName, abbr, phone string) (bank *model.Bank, err error) {
	bank = model.NewBank(name, englishName, abbr, phone)
	err = repository.Bank.Create(context.Background(), bank)
	if err != nil {
		return nil, err
	}
	return bank, nil
}

// Exist 判断银行是否存在
func (u *bankService) Exist(name, englishName string) (exist bool, err error) {
	return repository.Bank.ExistByName(context.Background(), name, englishName)
}

// Update 更新银行
func (u *bankService) Update(id int64, name, englishName, abbr, phone string) error {
	bank := model.NewBank(name, englishName, abbr, phone)
	bank.ID = id
	return repository.Bank.Update(context.Background(), bank)
}

// Delete 删除银行
func (u *bankService) Delete(id int64) error {
	return repository.Bank.Delete(context.Background(), id)
}
