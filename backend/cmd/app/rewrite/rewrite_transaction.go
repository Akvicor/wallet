package rewrite

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

// 清空wallet和card余额，重新执行交易
func rewriteTransaction() error {
	zero := decimal.NewFromInt(0)
	err := repository.Common.Transaction(context.Background(), func(ctx context.Context) error {
		// 处理钱包
		{
			// 获取所有钱包信息
			wallets, err := repository.UserWallet.FindAll(ctx, nil, false, model.NewPreloaderUserWallet().All(false))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			// 清空钱包划分余额
			for _, wallet := range wallets {
				partition := wallet.Partition
				for _, part := range partition {
					//part.Balance = zero
					err = repository.UserWalletPartition.UpdateBalance(ctx, false, part.WalletID, part.ID, zero)
					if err != nil {
						return err
					}
				}
			}
		}
		// 处理银行卡
		{
			// 获取所有银行卡信息
			cards, err := repository.UserCard.FindAll(ctx, nil, false, model.NewPreloaderUserCard().All())
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			// 清空银行卡各货币余额
			for _, card := range cards {
				currency := card.Currency
				for _, cur := range currency {
					//cur.Balance = zero
					err = repository.UserCardCurrency.UpdateBalance(ctx, false, cur.ID, zero)
					if err != nil {
						return err
					}
				}
			}
		}
		// 获取所有交易
		transactions, err := repository.UserTransaction.FindAll(ctx, nil, false, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		// 执行所有交易
		for _, tx := range transactions {
			if tx.Type == transaction.TypeIncome {
				toPartition, err := repository.UserWalletPartition.FindByID(ctx, false, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), tx.ToPartitionID)
				if err != nil {
					return err
				}
				if toPartition.Wallet == nil {
					return errors.New("cannot find wallet for partition")
				}
				if toPartition.Card == nil {
					return errors.New("cannot find wallet for partition card")
				}
				if toPartition.Card.Currency == nil {
					return errors.New("cannot find wallet for partition card currency")
				}
				// 增加目标划分余额
				toPartition.Balance = toPartition.Balance.Add(tx.ToValue)
				err = repository.UserWalletPartition.UpdateBalanceAdd(ctx, false, toPartition.Wallet.ID, toPartition.ID, tx.ToValue)
				if err != nil {
					return fmt.Errorf("update to partition balance failed: %v", err)
				}
				// 增加目标划分所属银行卡余额
				changed := false
				for _, currency := range toPartition.Card.Currency {
					if currency.CurrencyID == toPartition.CurrencyID {
						changed = true
						//currency.Balance = currency.Balance.Add(tx.ToValue)
						err = repository.UserCardCurrency.UpdateBalanceAdd(ctx, false, currency.ID, tx.ToValue)
						if err != nil {
							return fmt.Errorf("update to partition card currency balance failed: %v", err)
						}
						break
					}
				}
				if !changed {
					return errors.New("user card not changed")
				}
			} else if tx.Type == transaction.TypeExpense {
				fromPartition, err := repository.UserWalletPartition.FindByID(ctx, false, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), tx.FromPartitionID)
				if err != nil {
					return err
				}
				if fromPartition.Wallet == nil {
					return errors.New("cannot find wallet for partition")
				}
				if fromPartition.Card == nil {
					return errors.New("cannot find wallet for partition card")
				}
				if fromPartition.Card.Currency == nil {
					return errors.New("cannot find wallet for partition card currency")
				}
				// 减少源划分余额
				//fromPartition.Balance = fromPartition.Balance.Sub(tx.FromValue)
				err = repository.UserWalletPartition.UpdateBalanceSub(ctx, false, fromPartition.Wallet.ID, fromPartition.ID, tx.FromValue)
				if err != nil {
					return fmt.Errorf("update to partition balance failed: %v", err)
				}
				// 减少源划分所属银行卡余额
				changed := false
				for _, currency := range fromPartition.Card.Currency {
					if currency.CurrencyID == fromPartition.CurrencyID {
						changed = true
						//currency.Balance = currency.Balance.Sub(tx.FromValue)
						err = repository.UserCardCurrency.UpdateBalanceSub(ctx, false, currency.ID, tx.FromValue)
						if err != nil {
							return fmt.Errorf("update to partition card currency balance failed: %v", err)
						}
						break
					}
				}
				if !changed {
					return errors.New("user card not changed")
				}
			} else if tx.Type == transaction.TypeTransfer ||
				tx.Type == transaction.TypeAutoTransfer ||
				tx.Type == transaction.TypeExchange ||
				tx.Type == transaction.TypeAutoExchange {
				// 源
				fromPartition, err := repository.UserWalletPartition.FindByID(ctx, false, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), tx.FromPartitionID)
				if err != nil {
					return err
				}
				if fromPartition.Wallet == nil {
					return errors.New("cannot find wallet for partition")
				}
				if fromPartition.Card == nil {
					return errors.New("cannot find wallet for partition card")
				}
				if fromPartition.Card.Currency == nil {
					return errors.New("cannot find wallet for partition card currency")
				}
				// 减少源划分余额
				//fromPartition.Balance = fromPartition.Balance.Sub(tx.FromValue)
				err = repository.UserWalletPartition.UpdateBalanceSub(ctx, false, fromPartition.Wallet.ID, fromPartition.ID, tx.FromValue)
				if err != nil {
					return fmt.Errorf("update to partition balance failed: %v", err)
				}
				// 减少源划分所属银行卡余额
				changed := false
				for _, currency := range fromPartition.Card.Currency {
					if currency.CurrencyID == fromPartition.CurrencyID {
						changed = true
						//currency.Balance = currency.Balance.Sub(tx.FromValue)
						err = repository.UserCardCurrency.UpdateBalanceSub(ctx, false, currency.ID, tx.FromValue)
						if err != nil {
							return fmt.Errorf("update to partition card currency balance failed: %v", err)
						}
						break
					}
				}
				if !changed {
					return errors.New("user card not changed")
				}

				// 目标
				toPartition, err := repository.UserWalletPartition.FindByID(ctx, false, model.NewPreloaderUserWalletPartition().Wallet().Card().CardCurrency(), tx.ToPartitionID)
				if err != nil {
					return err
				}
				if toPartition.Wallet == nil {
					return errors.New("cannot find wallet for partition")
				}
				if toPartition.Card == nil {
					return errors.New("cannot find wallet for partition card")
				}
				if toPartition.Card.Currency == nil {
					return errors.New("cannot find wallet for partition card currency")
				}
				// 增加目标划分余额
				//toPartition.Balance = toPartition.Balance.Add(tx.ToValue)
				err = repository.UserWalletPartition.UpdateBalanceAdd(ctx, false, toPartition.Wallet.ID, toPartition.ID, tx.ToValue)
				if err != nil {
					return fmt.Errorf("update to partition balance failed: %v", err)
				}
				// 增加目标划分所属银行卡余额
				changed = false
				for _, currency := range toPartition.Card.Currency {
					if currency.CurrencyID == toPartition.CurrencyID {
						changed = true
						//currency.Balance = currency.Balance.Add(tx.ToValue)
						err = repository.UserCardCurrency.UpdateBalanceAdd(ctx, false, currency.ID, tx.ToValue)
						if err != nil {
							return fmt.Errorf("update to partition card currency balance failed: %v", err)
						}
						break
					}
				}
				if !changed {
					return errors.New("user card not changed")
				}
			}
		}
		return nil
	})
	return err
}
