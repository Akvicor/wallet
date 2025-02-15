package repository

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
	"wallet/cmd/app/server/model"
)

var UserCardCurrency = new(userCardCurrencyRepository)

type userCardCurrencyRepository struct {
	base[*model.UserCardCurrency]
}

/**
查找
*/

func (l *userCardCurrencyRepository) FindByCardID(c context.Context, alive bool, preloader *model.PreloaderUserCardCurrency, userCardID int64) ([]*model.UserCardCurrency, error) {
	result := make([]*model.UserCardCurrency, 0)
	err := l.preload(c, alive, preloader).Where("user_card_id = ?", userCardID).Find(&result).Error
	return result, err
}

func (l *userCardCurrencyRepository) FindByID(c context.Context, alive bool, preloader *model.PreloaderUserCardCurrency, id int64) (*model.UserCardCurrency, error) {
	result := new(model.UserCardCurrency)
	err := l.WrapResultErr(l.preload(c, alive, preloader).Where("id = ?", id).First(result))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *userCardCurrencyRepository) GetCurrencyIDSByCardID(c context.Context, alive bool, preloader *model.PreloaderUserCardCurrency, userCardID int64) ([]int64, error) {
	list := make([]*model.UserCardCurrency, 0)
	err := l.preload(c, alive, preloader).Where("user_card_id = ?", userCardID).Find(&list).Error
	if err != nil {
		return nil, err
	}
	result := make([]int64, len(list))
	for k, v := range list {
		result[k] = v.CurrencyID
	}
	return result, nil
}

// DetectValidByCardID 判断UserCardCurrency是否属于某个银行卡
func (l *userCardCurrencyRepository) DetectValidByCardID(c context.Context, cardID int64, ids []int64) (err error) {
	count := int64(0)
	err = l.dbt(c).Where("user_card_id = ? AND id IN (?)", cardID, ids).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("inconsistent quantity")
	}
	return nil
}

// DetectCurrencyValidByCardID 判断货币是否属于某个银行卡
func (l *userCardCurrencyRepository) DetectCurrencyValidByCardID(c context.Context, cardID int64, currencyIDS []int64) (err error) {
	count := int64(0)
	err = l.dbt(c).Where("user_card_id = ? AND currency_id IN (?)", cardID, currencyIDS).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(currencyIDS)) {
		return errors.New("inconsistent quantity")
	}
	return nil
}

/**
创建
*/

// Create 创建记录
func (l *userCardCurrencyRepository) Create(c context.Context, currency *model.UserCardCurrency) error {
	return l.WrapResultErr(l.db(c).Create(currency))
}

// Creates 批量创建记录
func (l *userCardCurrencyRepository) Creates(c context.Context, currency []*model.UserCardCurrency) error {
	return l.WrapResultErr(l.db(c).Create(currency))
}

/**
更新
*/

// UpdateBalance 更新余额
func (l *userCardCurrencyRepository) UpdateBalance(c context.Context, alive bool, id int64, value decimal.Decimal) error {
	return l.WrapResultErr(l.alive(c, alive).Where("id = ?", id).UpdateColumn("balance", value))
}

// UpdateBalanceAdd 增加余额
func (l *userCardCurrencyRepository) UpdateBalanceAdd(c context.Context, alive bool, id int64, value decimal.Decimal) error {
	err := Common.Transaction(c, func(ctx context.Context) error {
		currency := new(model.UserCardCurrency)
		var err error
		err = l.alive(ctx, alive).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(currency).Error
		if err != nil {
			return err
		}
		currency.Balance = currency.Balance.Add(value)
		return l.alive(ctx, alive).Where("id = ?", id).UpdateColumn("balance", currency.Balance).Error
	})
	return err
}

// UpdateBalanceSub 减少余额
func (l *userCardCurrencyRepository) UpdateBalanceSub(c context.Context, alive bool, id int64, value decimal.Decimal) error {
	err := Common.Transaction(c, func(ctx context.Context) error {
		currency := new(model.UserCardCurrency)
		var err error
		err = l.alive(ctx, alive).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(currency).Error
		if err != nil {
			return err
		}
		currency.Balance = currency.Balance.Sub(value)
		return l.alive(ctx, alive).Where("id = ?", id).UpdateColumn("balance", currency.Balance).Error
	})
	return err
}

/**
删除
*/

// Delete 通过ID删除记录
func (l *userCardCurrencyRepository) Delete(c context.Context, id int64, permitNonZero bool) error {
	tx := l.db(c).Where("id = ?", id)
	if !permitNonZero {
		tx = tx.Where("balance IS NULL OR balance = '0'")
	}
	return l.WrapResultErr(tx.Delete(&model.UserCardCurrency{}))
}

// DeleteAllByCardID 通过银行卡ID删除记录
func (l *userCardCurrencyRepository) DeleteAllByCardID(c context.Context, userCardID int64, permitNonZero bool) error {
	tx := l.db(c).Where("user_card_id = ?", userCardID)
	if !permitNonZero {
		tx = tx.Where("balance IS NULL OR balance = '0'")
	}
	return l.WrapResultErr(tx.Delete(&model.UserCardCurrency{}))
}

// DeleteByCardID 通过银行卡ID删除记录
func (l *userCardCurrencyRepository) DeleteByCardID(c context.Context, id, userCardID int64, permitNonZero bool) error {
	tx := l.db(c).Where("id = ? AND user_card_id = ?", id, userCardID)
	if !permitNonZero {
		tx = tx.Where("balance IS NULL OR balance = '0'")
	}
	return l.WrapResultErr(tx.Delete(&model.UserCardCurrency{}))
}
