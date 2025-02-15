package service

import (
	"context"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var CardType = new(cardTypeService)

type cardTypeService struct {
	base
}

// FindAll 获取所有银行卡种类
func (u *cardTypeService) FindAll(page *resp.PageModel, alive bool, preload *model.PreloaderCardType) (cardTypes []*model.CardType, err error) {
	return repository.CardType.FindAll(context.Background(), page, alive, preload)
}

// FindAllLike 获取所有银行卡种类
func (u *cardTypeService) FindAllLike(page *resp.PageModel, alive bool, preload *model.PreloaderCardType, like string) (cardTypes []*model.CardType, err error) {
	return repository.CardType.FindAllLike(context.Background(), page, alive, preload, like)
}

// Create 创建银行卡种类
func (u *cardTypeService) Create(name string) (cardType *model.CardType, err error) {
	cardType = model.NewCardType(name)
	err = repository.CardType.Create(context.Background(), cardType)
	if err != nil {
		return nil, err
	}
	return cardType, nil
}

// Exist 判断卡类型是否存在
func (u *cardTypeService) Exist(name string) (exist bool, err error) {
	return repository.CardType.ExistByName(context.Background(), name)
}

// Update 更新银行卡种类
func (u *cardTypeService) Update(id int64, name string) error {
	cardType := model.NewCardType(name)
	cardType.ID = id
	return repository.CardType.Update(context.Background(), cardType)
}

// Delete 删除银行卡种类
func (u *cardTypeService) Delete(id int64) error {
	return repository.CardType.Delete(context.Background(), id)
}
