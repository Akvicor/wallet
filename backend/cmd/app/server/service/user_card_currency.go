package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var UserCardCurrency = new(userCardCurrencyService)

type userCardCurrencyService struct {
	base
}

// Create 创建银行卡绑定货币
func (w *userCardCurrencyService) Create(uid, cardID, currencyID int64) (userCardCurrency *model.UserCardCurrency, err error) {
	userCardCurrency = model.NewUserCardCurrency(cardID, currencyID, decimal.NewFromInt(0))
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断银行卡是否属于用户
		err = repository.UserCard.DetectValidByUID(ctx, uid, []int64{cardID})
		if err != nil {
			return err
		}
		// 创建钱包
		err = repository.UserCardCurrency.Create(ctx, userCardCurrency)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return userCardCurrency, nil
}

// UpdateBalance 更新银行卡绑定货币余额
func (w *userCardCurrencyService) UpdateBalance(uid, id int64, balance decimal.Decimal) (err error) {
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 获取银行卡货币余额
		userCurrencyBalance, err := repository.UserCardCurrency.FindByID(ctx, true, model.NewPreloaderUserCardCurrency().UserCard(), id)
		if err != nil {
			return err
		}
		if userCurrencyBalance.UserCard == nil {
			return errors.New("find user card failed")
		}
		if userCurrencyBalance.UserCard.UID != uid {
			return errors.New("invalid id")
		}
		// 更新余额
		err = repository.UserCardCurrency.UpdateBalance(ctx, true, id, balance)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByID 删除银行卡绑定货币
func (w *userCardCurrencyService) DeleteByID(uid, id int64) (err error) {
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 获取银行卡货币余额
		userCurrencyBalance, err := repository.UserCardCurrency.FindByID(ctx, true, model.NewPreloaderUserCardCurrency().UserCard(), id)
		if err != nil {
			return fmt.Errorf("find user card currency failed: %v", err)
		}
		if userCurrencyBalance.UserCard == nil {
			return errors.New("find user card failed")
		}
		if userCurrencyBalance.UserCard.UID != uid {
			return errors.New("invalid id")
		}
		if userCurrencyBalance.CurrencyID == userCurrencyBalance.UserCard.MasterCurrencyID {
			return errors.New("can not delete master currency")
		}
		err = repository.UserCardCurrency.Delete(ctx, id, false)
		return err
	})
	return err
}
