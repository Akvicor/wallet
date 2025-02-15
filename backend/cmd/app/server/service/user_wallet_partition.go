package service

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
	"wallet/cmd/app/server/common/types/wallet_partition"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var UserWalletPartition = new(userWalletPartitionService)

type userWalletPartitionService struct {
	base
}

// Create 创建划分
func (w *userWalletPartitionService) Create(uid, walletID, cardID, currencyID int64, name, description string, limit decimal.Decimal, average wallet_partition.AverageType) (userWalletPartition *model.UserWalletPartition, err error) {
	userWalletPartition = model.NewUserWalletPartition(walletID, cardID, currencyID, name, description, limit, decimal.NewFromInt(0), average)
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断钱包是否属于用户
		err = repository.UserWallet.DetectValidByUID(ctx, uid, []int64{walletID})
		if err != nil {
			return err
		}
		// 判断银行卡是否属于用户
		err = repository.UserCard.DetectValidByUID(ctx, uid, []int64{cardID})
		if err != nil {
			return err
		}
		// 判断货币是否属于银行卡
		err = repository.UserCardCurrency.DetectCurrencyValidByCardID(ctx, cardID, []int64{currencyID})
		if err != nil {
			return err
		}
		// 获取序号
		seq, err := repository.UserWalletPartition.GetMaxSequenceByWalletID(ctx, true, walletID)
		if err != nil {
			seq = 0
		}
		userWalletPartition.Sequence = seq + 1
		// 创建钱包
		err = repository.UserWalletPartition.Create(ctx, userWalletPartition)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return userWalletPartition, nil
}

// Update 更新用户的钱包
func (w *userWalletPartitionService) Update(uid, id, walletID int64, name, description string, limit decimal.Decimal, average wallet_partition.AverageType) (err error) {
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断钱包是否属于用户
		err = repository.UserWallet.DetectValidByUID(ctx, uid, []int64{walletID})
		if err != nil {
			return err
		}
		// 判断划分是否属于钱包
		err = repository.UserWalletPartition.DetectValidByWalletID(ctx, walletID, []int64{id})
		if err != nil {
			return err
		}
		// 更新划分
		err = repository.UserWalletPartition.Update(ctx, id, name, description, limit, average)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateSequenceByUID 停用用户钱包
func (w *userWalletPartitionService) UpdateSequenceByUID(uid, walletID, id, targetSequence int64) error {
	err := w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断钱包是否属于用户
		err := repository.UserWallet.DetectValidByUID(ctx, uid, []int64{walletID})
		if err != nil {
			return err
		}
		// 判断划分是否属于钱包
		err = repository.UserWalletPartition.DetectValidByWalletID(ctx, walletID, []int64{id})
		if err != nil {
			return err
		}
		maxSequence, err := repository.UserWalletPartition.GetMaxSequenceByWalletID(ctx, false, walletID)
		if err != nil {
			return nil
		}
		if targetSequence < 1 {
			targetSequence = 1
		}
		if targetSequence > maxSequence {
			targetSequence = maxSequence
		}
		origin, err := repository.UserWalletPartition.FindByWalletID(ctx, false, nil, walletID, id)
		if err != nil {
			return err
		}
		err = repository.UserWalletPartition.UpdatesSequenceByWalletID(ctx, false, walletID, origin.Sequence, targetSequence)
		if err != nil {
			return err
		}
		err = repository.UserWalletPartition.UpdateSequenceByWalletID(ctx, false, walletID, id, targetSequence)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// DisableByID 停用划分
func (w *userWalletPartitionService) DisableByID(uid, walletID, id int64) (err error) {
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断钱包是否属于用户
		err = repository.UserWallet.DetectValidByUID(ctx, uid, []int64{walletID})
		if err != nil {
			return err
		}
		// 判断划分是否属于钱包
		err = repository.UserWalletPartition.DetectValidByWalletID(ctx, walletID, []int64{id})
		if err != nil {
			return err
		}
		// 判断划分是否存在关联的交易
		var count int64
		count, err = repository.UserTransaction.DetectExistByPartitionID(ctx, uid, id)
		if err != nil {
			return err
		}
		if count == 0 {
			// 不存在关联的交易，删除划分
			err = repository.UserWalletPartition.DeleteByID(ctx, true, id)
		} else {
			// 存在关联的交易，停用划分
			err = repository.UserWalletPartition.UpdateDisabledByID(ctx, true, id, time.Now().Unix())
		}
		return err
	})
	return err
}

// EnableByID 启用划分
func (w *userWalletPartitionService) EnableByID(uid, walletID, id int64) (err error) {
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断钱包是否属于用户
		err = repository.UserWallet.DetectValidByUID(ctx, uid, []int64{walletID})
		if err != nil {
			return err
		}
		// 判断划分是否属于钱包
		err = repository.UserWalletPartition.DetectValidByWalletID(ctx, walletID, []int64{id})
		if err != nil {
			return err
		}
		err = repository.UserWalletPartition.UpdateDisabledByID(ctx, false, id, 0)
		return err
	})
	return err
}
