package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var UserTransaction = new(userTransactionService)

type userTransactionService struct {
	base
}

// FindAllByUID 通过UID获取所有交易
func (u *userTransactionService) FindAllByUID(page *resp.PageModel, alive bool, preload *model.PreloaderUserTransaction, uid int64, like string, fromPartitionID, toPartitionID, fromCurrencyID, toCurrencyID, categoryID []int64, transactionType []transaction.Type) (transactions []*model.UserTransaction, err error) {
	return repository.UserTransaction.FindAllByUID(context.Background(), page, alive, preload, uid, like, fromPartitionID, toPartitionID, fromCurrencyID, toCurrencyID, categoryID, transactionType)
}

// FindAllByUIDCreated 通过UID获取所有交易
func (u *userTransactionService) FindAllByUIDCreated(page *resp.PageModel, alive bool, preload *model.PreloaderUserTransaction, uid, start, end int64) (transactions []*model.UserTransaction, err error) {
	return repository.UserTransaction.FindAllByUIDCreated(context.Background(), page, alive, preload, uid, start, end)
}

// FindAllByUIDCreatedFromPartition 通过UID获取所有交易
func (u *userTransactionService) FindAllByUIDCreatedFromPartition(page *resp.PageModel, alive bool, preload *model.PreloaderUserTransaction, uid int64, fromPartitionID []int64, start, end int64) (transactions []*model.UserTransaction, err error) {
	return repository.UserTransaction.FindAllByUIDCreatedFromPartition(context.Background(), page, alive, preload, uid, fromPartitionID, start, end)
}

// FindAllByUIDUnchecked 通过UID获取所有未确认交易
func (u *userTransactionService) FindAllByUIDUnchecked(page *resp.PageModel, alive bool, preload *model.PreloaderUserTransaction, uid int64) (transactions []*model.UserTransaction, err error) {
	return repository.UserTransaction.FindAllByUIDUnchecked(context.Background(), page, alive, preload, uid)
}

// FindByUID 获取所有的交易, page为nil时不分页
func (u *userTransactionService) FindByUID(alive bool, preloader *model.PreloaderUserTransaction, uid, id int64) (transactions *model.UserTransaction, err error) {
	return repository.UserTransaction.FindByUID(context.Background(), alive, preloader, uid, id)
}

// CreateIncome 创建收入交易
func (u *userTransactionService) CreateIncome(transactionType transaction.Type, uid, toPartitionID, categoryID int64, description string, value decimal.Decimal, created int64) (*model.UserTransaction, error) {
	transactions := model.NewUserTransaction(uid, 0, 0, toPartitionID, 0, categoryID, transactionType, description, value, value, decimal.NewFromInt(0), created)
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断目标钱包和划分
		var toPartition *model.UserWalletPartition
		toPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), toPartitionID)
		if e != nil {
			return fmt.Errorf("find to partition failed: %v", e)
		}
		if toPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if toPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if toPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if toPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 补充交易信息
		transactions.ToCurrencyID = toPartition.CurrencyID
		// 判断交易分类是否属于用户
		e = repository.UserTransactionCategory.DetectValidByUIDType(ctx, uid, transactions.Type, []int64{categoryID})
		if e != nil {
			return fmt.Errorf("find transaction category failed: %v", e)
		}
		// 增加收款划分余额
		//toPartition.Balance = toPartition.Balance.Add(transactions.ToValue)
		e = repository.UserWalletPartition.UpdateBalanceAdd(ctx, true, toPartition.Wallet.ID, toPartition.ID, transactions.ToValue)
		if e != nil {
			return fmt.Errorf("update to partition balance failed: %v", e)
		}
		// 增加收款银行卡余额
		changed := false
		for _, currency := range toPartition.Card.Currency {
			if currency.CurrencyID == toPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Add(transactions.ToValue)
				e = repository.UserCardCurrency.UpdateBalanceAdd(ctx, true, currency.ID, transactions.ToValue)
				if e != nil {
					return fmt.Errorf("update to partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 创建用户交易
		e = repository.UserTransaction.Create(ctx, transactions)
		if e != nil {
			return fmt.Errorf("create transaction failed: %v", e)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// CreateExpense 创建支出交易
func (u *userTransactionService) CreateExpense(transactionType transaction.Type, uid, fromPartitionID, currencyID, categoryID int64, description string, value decimal.Decimal, created int64) (*model.UserTransaction, error) {
	transactions := model.NewUserTransaction(uid, fromPartitionID, 0, 0, 0, categoryID, transactionType, description, decimal.NewFromInt(0), value, decimal.NewFromInt(0), created)
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断货币是否存在
		exist, e = repository.Currency.ExistById(ctx, true, []int64{currencyID})
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("货币已被停用")
		}
		// 判断目标钱包和划分
		var fromPartition *model.UserWalletPartition
		fromPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), fromPartitionID)
		if e != nil {
			return fmt.Errorf("find from partition failed: %v", e)
		}
		if fromPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if fromPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if fromPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if fromPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		rate := decimal.NewFromInt(1)
		if currencyID != fromPartition.CurrencyID {
			// 查询汇率
			var exchange *model.UserExchangeRate
			exchange, e = repository.UserExchangeRate.FindByUID(ctx, true, nil, uid, currencyID, fromPartition.CurrencyID)
			if e != nil {
				return fmt.Errorf("find exchange rate failed: %v", e)
			}
			rate = exchange.Rate
			transactions.ToCurrencyID = currencyID
		}
		// 补充交易信息
		transactions.FromCurrencyID = fromPartition.CurrencyID
		transactions.FromValue = transactions.ToValue.Mul(rate)
		// 判断交易分类是否属于用户
		e = repository.UserTransactionCategory.DetectValidByUIDType(ctx, uid, transactions.Type, []int64{categoryID})
		if e != nil {
			return fmt.Errorf("find transaction category failed: %v", e)
		}
		// 减少支出划分余额
		//fromPartition.Balance = fromPartition.Balance.Sub(transactions.FromValue)
		e = repository.UserWalletPartition.UpdateBalanceSub(ctx, true, fromPartition.Wallet.ID, fromPartition.ID, transactions.FromValue)
		if e != nil {
			return fmt.Errorf("update from partition balance failed: %v", e)
		}
		// 减少支出银行卡余额
		changed := false
		for _, currency := range fromPartition.Card.Currency {
			if currency.CurrencyID == fromPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Sub(transactions.FromValue)
				e = repository.UserCardCurrency.UpdateBalanceSub(ctx, true, currency.ID, transactions.FromValue)
				if e != nil {
					return fmt.Errorf("update from partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 创建用户交易
		e = repository.UserTransaction.Create(ctx, transactions)
		if e != nil {
			return fmt.Errorf("create transaction failed: %v", e)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// CreateTransfer 创建转账交易
func (u *userTransactionService) CreateTransfer(transactionType transaction.Type, uid, fromPartitionID, toPartitionID, categoryID int64, description string, fromValue, toValue decimal.Decimal, created int64) (*model.UserTransaction, error) {
	if fromValue.LessThan(toValue) {
		return nil, errors.New("from value must >= to value")
	}
	transactions := model.NewUserTransaction(uid, fromPartitionID, 0, toPartitionID, 0, categoryID, transactionType, description, fromValue, toValue, fromValue.Sub(toValue), created)
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断目标钱包和划分
		var fromPartition *model.UserWalletPartition
		fromPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), fromPartitionID)
		if e != nil {
			return fmt.Errorf("find from partition failed: %v", e)
		}
		if fromPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if fromPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if fromPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if fromPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 判断目标钱包和划分
		var toPartition *model.UserWalletPartition
		toPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), toPartitionID)
		if e != nil {
			return fmt.Errorf("find to partition failed: %v", e)
		}
		if toPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if toPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if toPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if toPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 判断货币类型是否一致
		if fromPartition.CurrencyID != toPartition.CurrencyID {
			return errors.New("can't transfer to different currency, please use exchange")
		}
		// 补充交易信息
		transactions.FromCurrencyID = fromPartition.CurrencyID
		transactions.ToCurrencyID = toPartition.CurrencyID
		// 判断交易分类是否属于用户
		e = repository.UserTransactionCategory.DetectValidByUIDType(ctx, uid, transactions.Type, []int64{categoryID})
		if e != nil {
			return fmt.Errorf("find transaction category failed: %v", e)
		}
		// 减少支出划分余额
		//fromPartition.Balance = fromPartition.Balance.Sub(transactions.FromValue)
		e = repository.UserWalletPartition.UpdateBalanceSub(ctx, true, fromPartition.Wallet.ID, fromPartition.ID, transactions.FromValue)
		if e != nil {
			return fmt.Errorf("update from partition balance failed: %v", e)
		}
		// 减少支出银行卡余额
		changed := false
		for _, currency := range fromPartition.Card.Currency {
			if currency.CurrencyID == fromPartition.CurrencyID {
				changed = true
				currency.Balance = currency.Balance.Sub(transactions.FromValue)
				e = repository.UserCardCurrency.UpdateBalanceSub(ctx, true, currency.ID, transactions.FromValue)
				if e != nil {
					return fmt.Errorf("update from partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 增加收款划分余额
		toPartition.Balance = toPartition.Balance.Add(transactions.ToValue)
		e = repository.UserWalletPartition.UpdateBalanceAdd(ctx, true, toPartition.Wallet.ID, toPartition.ID, transactions.ToValue)
		if e != nil {
			return fmt.Errorf("update to partition balance failed: %v", e)
		}
		// 增加收款银行卡余额
		changed = false
		for _, currency := range toPartition.Card.Currency {
			if currency.CurrencyID == toPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Add(transactions.ToValue)
				e = repository.UserCardCurrency.UpdateBalanceAdd(ctx, true, currency.ID, transactions.ToValue)
				if e != nil {
					return fmt.Errorf("update to partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 创建用户交易
		e = repository.UserTransaction.Create(ctx, transactions)
		if e != nil {
			return fmt.Errorf("create transaction failed: %v", e)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// CreateAutoTransfer 创建自动转账交易
func (u *userTransactionService) CreateAutoTransfer(transactionType transaction.Type, uid, fromPartitionID, toPartitionID, currencyID, categoryID int64, description string, fromValue decimal.Decimal, created int64) (*model.UserTransaction, error) {
	transactions := model.NewUserTransaction(uid, fromPartitionID, 0, toPartitionID, 0, categoryID, transactionType, description, decimal.NewFromInt(0), decimal.NewFromInt(0), decimal.NewFromInt(0), created)
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断目标钱包和划分
		var fromPartition *model.UserWalletPartition
		fromPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), fromPartitionID)
		if e != nil {
			return fmt.Errorf("find from partition failed: %v", e)
		}
		if fromPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if fromPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if fromPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if fromPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 判断目标钱包和划分
		var toPartition *model.UserWalletPartition
		toPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), toPartitionID)
		if e != nil {
			return fmt.Errorf("find to partition failed: %v", e)
		}
		if toPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if toPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if toPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if toPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 判断货币类型是否一致
		if fromPartition.CurrencyID != toPartition.CurrencyID {
			return errors.New("can't auto-transfer to different currency, please use auto-exchange")
		}
		// 补充交易信息
		transactions.FromCurrencyID = fromPartition.CurrencyID
		transactions.ToCurrencyID = toPartition.CurrencyID
		// 判断交易分类是否属于用户
		e = repository.UserTransactionCategory.DetectValidByUIDType(ctx, uid, transactions.Type, []int64{categoryID})
		if e != nil {
			return fmt.Errorf("find transaction category failed: %v", e)
		}
		// 查找汇率
		var exchange *model.UserExchangeRate
		exchange, e = repository.UserExchangeRate.FindByUID(ctx, true, nil, uid, currencyID, transactions.FromCurrencyID)
		if e != nil {
			return fmt.Errorf("find exchange failed: %v", e)
		}
		transactions.FromValue = fromValue.Mul(exchange.Rate)
		transactions.ToValue = transactions.FromValue
		// 减少支出划分余额
		//fromPartition.Balance = fromPartition.Balance.Sub(transactions.FromValue)
		e = repository.UserWalletPartition.UpdateBalanceSub(ctx, true, fromPartition.Wallet.ID, fromPartition.ID, transactions.FromValue)
		if e != nil {
			return fmt.Errorf("update from partition balance failed: %v", e)
		}
		// 减少支出银行卡余额
		changed := false
		for _, currency := range fromPartition.Card.Currency {
			if currency.CurrencyID == fromPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Sub(transactions.FromValue)
				e = repository.UserCardCurrency.UpdateBalanceSub(ctx, true, currency.ID, transactions.FromValue)
				if e != nil {
					return fmt.Errorf("update from partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 增加收款划分余额
		toPartition.Balance = toPartition.Balance.Add(transactions.ToValue)
		e = repository.UserWalletPartition.UpdateBalanceAdd(ctx, true, toPartition.Wallet.ID, toPartition.ID, transactions.ToValue)
		if e != nil {
			return fmt.Errorf("update to partition balance failed: %v", e)
		}
		// 增加收款银行卡余额
		changed = false
		for _, currency := range toPartition.Card.Currency {
			if currency.CurrencyID == toPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Add(transactions.ToValue)
				e = repository.UserCardCurrency.UpdateBalanceAdd(ctx, true, currency.ID, transactions.ToValue)
				if e != nil {
					return fmt.Errorf("update to partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 创建用户交易
		e = repository.UserTransaction.Create(ctx, transactions)
		if e != nil {
			return fmt.Errorf("create transaction failed: %v", e)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// CreateExchange 创建兑换交易
func (u *userTransactionService) CreateExchange(transactionType transaction.Type, uid, fromPartitionID, toPartitionID, categoryID int64, description string, fromValue, toValue decimal.Decimal, created int64) (*model.UserTransaction, error) {
	transactions := model.NewUserTransaction(uid, fromPartitionID, 0, toPartitionID, 0, categoryID, transactionType, description, fromValue, toValue, decimal.NewFromInt(0), created)
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断目标钱包和划分
		var fromPartition *model.UserWalletPartition
		fromPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), fromPartitionID)
		if e != nil {
			return fmt.Errorf("find from partition failed: %v", e)
		}
		if fromPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if fromPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if fromPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if fromPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 判断目标钱包和划分
		var toPartition *model.UserWalletPartition
		toPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), toPartitionID)
		if e != nil {
			return fmt.Errorf("find to partition failed: %v", e)
		}
		if toPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if toPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if toPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if toPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 判断货币类型是否一致
		if fromPartition.CurrencyID == toPartition.CurrencyID {
			return errors.New("can't exchange to same currency, please use transfer")
		}
		// 补充交易信息
		transactions.FromCurrencyID = fromPartition.CurrencyID
		transactions.ToCurrencyID = toPartition.CurrencyID
		// 判断交易分类是否属于用户
		e = repository.UserTransactionCategory.DetectValidByUIDType(ctx, uid, transactions.Type, []int64{categoryID})
		if e != nil {
			return fmt.Errorf("find transaction category failed: %v", e)
		}
		// 减少支出划分余额
		//fromPartition.Balance = fromPartition.Balance.Sub(transactions.FromValue)
		e = repository.UserWalletPartition.UpdateBalanceSub(ctx, true, fromPartition.Wallet.ID, fromPartition.ID, transactions.FromValue)
		if e != nil {
			return fmt.Errorf("update from partition balance failed: %v", e)
		}
		// 减少支出银行卡余额
		changed := false
		for _, currency := range fromPartition.Card.Currency {
			if currency.CurrencyID == fromPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Sub(transactions.FromValue)
				e = repository.UserCardCurrency.UpdateBalanceSub(ctx, true, currency.ID, transactions.FromValue)
				if e != nil {
					return fmt.Errorf("update from partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 增加收款划分余额
		toPartition.Balance = toPartition.Balance.Add(transactions.ToValue)
		e = repository.UserWalletPartition.UpdateBalanceAdd(ctx, true, toPartition.Wallet.ID, toPartition.ID, transactions.ToValue)
		if e != nil {
			return fmt.Errorf("update to partition balance failed: %v", e)
		}
		// 增加收款银行卡余额
		changed = false
		for _, currency := range toPartition.Card.Currency {
			if currency.CurrencyID == toPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Add(transactions.ToValue)
				e = repository.UserCardCurrency.UpdateBalanceAdd(ctx, true, currency.ID, transactions.ToValue)
				if e != nil {
					return fmt.Errorf("update to partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 创建用户交易
		e = repository.UserTransaction.Create(ctx, transactions)
		if e != nil {
			return fmt.Errorf("create transaction failed: %v", e)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// CreateAutoExchange 创建自动兑换交易
func (u *userTransactionService) CreateAutoExchange(transactionType transaction.Type, uid, fromPartitionID, toPartitionID, categoryID int64, description string, toValue decimal.Decimal, created int64) (*model.UserTransaction, error) {
	transactions := model.NewUserTransaction(uid, fromPartitionID, 0, toPartitionID, 0, categoryID, transactionType, description, decimal.NewFromInt(0), toValue, decimal.NewFromInt(0), created)
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断目标钱包和划分
		var fromPartition *model.UserWalletPartition
		fromPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), fromPartitionID)
		if e != nil {
			return fmt.Errorf("find from partition failed: %v", e)
		}
		if fromPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if fromPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if fromPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if fromPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 判断目标钱包和划分
		var toPartition *model.UserWalletPartition
		toPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), toPartitionID)
		if e != nil {
			return fmt.Errorf("find to partition failed: %v", e)
		}
		if toPartition.Wallet == nil {
			return errors.New("cannot find wallet for partition")
		}
		if toPartition.Wallet.UID != uid {
			return errors.New("partition not found")
		}
		if toPartition.Card == nil {
			return errors.New("cannot find wallet for partition card")
		}
		if toPartition.Card.Currency == nil {
			return errors.New("cannot find wallet for partition card currency")
		}
		// 判断货币类型是否一致
		if fromPartition.CurrencyID == toPartition.CurrencyID {
			return errors.New("can't auto-exchange to same currency, please use auto-transfer")
		}
		// 补充交易信息
		transactions.FromCurrencyID = fromPartition.CurrencyID
		transactions.ToCurrencyID = toPartition.CurrencyID
		// 判断交易分类是否属于用户
		e = repository.UserTransactionCategory.DetectValidByUIDType(ctx, uid, transactions.Type, []int64{categoryID})
		if e != nil {
			return fmt.Errorf("find transaction category failed: %v", e)
		}
		// 查找汇率
		var exchange *model.UserExchangeRate
		exchange, e = repository.UserExchangeRate.FindByUID(ctx, true, nil, uid, transactions.ToCurrencyID, transactions.FromCurrencyID)
		if e != nil {
			return fmt.Errorf("find exchange failed: %v", e)
		}
		transactions.FromValue = transactions.ToValue.Mul(exchange.Rate)
		// 减少支出划分余额
		//fromPartition.Balance = fromPartition.Balance.Sub(transactions.FromValue)
		e = repository.UserWalletPartition.UpdateBalanceSub(ctx, true, fromPartition.Wallet.ID, fromPartition.ID, transactions.FromValue)
		if e != nil {
			return fmt.Errorf("update from partition balance failed: %v", e)
		}
		// 减少支出银行卡余额
		changed := false
		for _, currency := range fromPartition.Card.Currency {
			if currency.CurrencyID == fromPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Sub(transactions.FromValue)
				e = repository.UserCardCurrency.UpdateBalanceSub(ctx, true, currency.ID, transactions.FromValue)
				if e != nil {
					return fmt.Errorf("update from partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 增加收款划分余额
		//toPartition.Balance = toPartition.Balance.Add(transactions.ToValue)
		e = repository.UserWalletPartition.UpdateBalanceAdd(ctx, true, toPartition.Wallet.ID, toPartition.ID, transactions.ToValue)
		if e != nil {
			return fmt.Errorf("update to partition balance failed: %v", e)
		}
		// 增加收款银行卡余额
		changed = false
		for _, currency := range toPartition.Card.Currency {
			if currency.CurrencyID == toPartition.CurrencyID {
				changed = true
				//currency.Balance = currency.Balance.Add(transactions.ToValue)
				e = repository.UserCardCurrency.UpdateBalanceAdd(ctx, true, currency.ID, transactions.ToValue)
				if e != nil {
					return fmt.Errorf("update to partition card currency balance failed: %v", e)
				}
				break
			}
		}
		if !changed {
			return errors.New("user card not changed")
		}
		// 创建用户交易
		e = repository.UserTransaction.Create(ctx, transactions)
		if e != nil {
			return fmt.Errorf("create transaction failed: %v", e)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// DeleteByUID 删除用户交易
func (u *userTransactionService) DeleteByUID(uid, id int64) error {
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		// 获取旧交易信息
		oldTransaction, e := repository.UserTransaction.FindByUID(ctx, true, model.NewPreloaderUserTransaction().FromPartition().ToPartition(), uid, id)
		if e != nil {
			return e
		}
		if oldTransaction.Type == transaction.TypeIncome {
			var toPartition *model.UserWalletPartition
			toPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), oldTransaction.ToPartitionID)
			if e != nil {
				return fmt.Errorf("find to partition failed: %v", e)
			}
			if toPartition.Wallet == nil {
				return errors.New("cannot find wallet for partition")
			}
			if toPartition.Wallet.UID != uid {
				return errors.New("partition not found")
			}
			if toPartition.Card == nil {
				return errors.New("cannot find wallet for partition card")
			}
			if toPartition.Card.Currency == nil {
				return errors.New("cannot find wallet for partition card currency")
			}
			// 回滚收款划分余额
			//toPartition.Balance = toPartition.Balance.Sub(oldTransaction.ToValue)
			e = repository.UserWalletPartition.UpdateBalanceSub(ctx, true, toPartition.Wallet.ID, toPartition.ID, oldTransaction.ToValue)
			if e != nil {
				return fmt.Errorf("update to partition balance failed: %v", e)
			}
			// 回滚收款银行卡余额
			changed := false
			for _, currency := range toPartition.Card.Currency {
				if currency.CurrencyID == toPartition.CurrencyID {
					changed = true
					//currency.Balance = currency.Balance.Sub(oldTransaction.ToValue)
					e = repository.UserCardCurrency.UpdateBalanceSub(ctx, true, currency.ID, oldTransaction.ToValue)
					if e != nil {
						return fmt.Errorf("update to partition card currency balance failed: %v", e)
					}
					break
				}
			}
			if !changed {
				return errors.New("user card not changed")
			}
		} else if oldTransaction.Type == transaction.TypeExpense {
			// 判断目标钱包和划分
			var fromPartition *model.UserWalletPartition
			fromPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), oldTransaction.FromPartitionID)
			if e != nil {
				return fmt.Errorf("find from partition failed: %v", e)
			}
			if fromPartition.Wallet == nil {
				return errors.New("cannot find wallet for partition")
			}
			if fromPartition.Wallet.UID != uid {
				return errors.New("partition not found")
			}
			if fromPartition.Card == nil {
				return errors.New("cannot find wallet for partition card")
			}
			if fromPartition.Card.Currency == nil {
				return errors.New("cannot find wallet for partition card currency")
			}
			// 回滚支出划分余额
			//fromPartition.Balance = fromPartition.Balance.Add(oldTransaction.FromValue)
			e = repository.UserWalletPartition.UpdateBalanceAdd(ctx, true, fromPartition.Wallet.ID, fromPartition.ID, oldTransaction.FromValue)
			if e != nil {
				return fmt.Errorf("update from partition balance failed: %v", e)
			}
			// 回滚支出银行卡余额
			changed := false
			for _, currency := range fromPartition.Card.Currency {
				if currency.CurrencyID == fromPartition.CurrencyID {
					changed = true
					//currency.Balance = currency.Balance.Add(oldTransaction.FromValue)
					e = repository.UserCardCurrency.UpdateBalanceAdd(ctx, true, currency.ID, oldTransaction.FromValue)
					if e != nil {
						return fmt.Errorf("update from partition card currency balance failed: %v", e)
					}
					break
				}
			}
			if !changed {
				return errors.New("user card not changed")
			}
		} else if oldTransaction.Type == transaction.TypeTransfer ||
			oldTransaction.Type == transaction.TypeAutoTransfer ||
			oldTransaction.Type == transaction.TypeExchange ||
			oldTransaction.Type == transaction.TypeAutoExchange {
			// 判断目标钱包和划分
			var fromPartition *model.UserWalletPartition
			fromPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), oldTransaction.FromPartitionID)
			if e != nil {
				return fmt.Errorf("find from partition failed: %v", e)
			}
			if fromPartition.Wallet == nil {
				return errors.New("cannot find wallet for partition")
			}
			if fromPartition.Wallet.UID != uid {
				return errors.New("partition not found")
			}
			if fromPartition.Card == nil {
				return errors.New("cannot find wallet for partition card")
			}
			if fromPartition.Card.Currency == nil {
				return errors.New("cannot find wallet for partition card currency")
			}
			// 判断目标钱包和划分
			var toPartition *model.UserWalletPartition
			toPartition, e = repository.UserWalletPartition.FindByID(ctx, true, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), oldTransaction.ToPartitionID)
			if e != nil {
				return fmt.Errorf("find to partition failed: %v", e)
			}
			if toPartition.Wallet == nil {
				return errors.New("cannot find wallet for partition")
			}
			if toPartition.Wallet.UID != uid {
				return errors.New("partition not found")
			}
			if toPartition.Card == nil {
				return errors.New("cannot find wallet for partition card")
			}
			if toPartition.Card.Currency == nil {
				return errors.New("cannot find wallet for partition card currency")
			}
			// 增加支出划分余额
			//fromPartition.Balance = fromPartition.Balance.Add(oldTransaction.FromValue)
			e = repository.UserWalletPartition.UpdateBalanceAdd(ctx, true, fromPartition.Wallet.ID, fromPartition.ID, oldTransaction.FromValue)
			if e != nil {
				return fmt.Errorf("update from partition balance failed: %v", e)
			}
			// 增加支出银行卡余额
			changed := false
			for _, currency := range fromPartition.Card.Currency {
				if currency.CurrencyID == fromPartition.CurrencyID {
					changed = true
					//currency.Balance = currency.Balance.Add(oldTransaction.FromValue)
					e = repository.UserCardCurrency.UpdateBalanceAdd(ctx, true, currency.ID, oldTransaction.FromValue)
					if e != nil {
						return fmt.Errorf("update from partition card currency balance failed: %v", e)
					}
					break
				}
			}
			if !changed {
				return errors.New("user card not changed")
			}
			// 减少收款划分余额
			//toPartition.Balance = toPartition.Balance.Sub(oldTransaction.ToValue)
			e = repository.UserWalletPartition.UpdateBalanceSub(ctx, true, toPartition.Wallet.ID, toPartition.ID, oldTransaction.ToValue)
			if e != nil {
				return fmt.Errorf("update to partition balance failed: %v", e)
			}
			// 减少收款银行卡余额
			changed = false
			for _, currency := range toPartition.Card.Currency {
				if currency.CurrencyID == toPartition.CurrencyID {
					changed = true
					//currency.Balance = currency.Balance.Sub(oldTransaction.ToValue)
					e = repository.UserCardCurrency.UpdateBalanceSub(ctx, true, currency.ID, oldTransaction.ToValue)
					if e != nil {
						return fmt.Errorf("update to partition card currency balance failed: %v", e)
					}
					break
				}
			}
			if !changed {
				return errors.New("user card not changed")
			}
		} else {
			return errors.New("unknown transaction type")
		}
		e = repository.UserTransaction.DeleteByUID(ctx, uid, id)
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

// UpdateByUID 更新
func (u *userTransactionService) UpdateByUID(uid, id, categoryID int64, description string, created int64) error {
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		// 查找交易
		tx, err := repository.UserTransaction.FindByUID(ctx, false, nil, uid, id)
		if err != nil {
			return fmt.Errorf("find user transaction failed: %v", err)
		}
		// 判断交易分类是否属于用户
		err = repository.UserTransactionCategory.DetectValidByUIDType(ctx, uid, tx.Type, []int64{categoryID})
		if err != nil {
			return fmt.Errorf("find transaction category failed: %v", err)
		}
		return repository.UserTransaction.UpdateByUID(ctx, uid, id, categoryID, description, created)
	})
	return err
}

// UpdateCheckedByUID 更新Checked为当前时间
func (u *userTransactionService) UpdateCheckedByUID(uid, id int64) error {
	return repository.UserTransaction.UpdateCheckedByUID(context.Background(), uid, id, time.Now().Unix())
}
