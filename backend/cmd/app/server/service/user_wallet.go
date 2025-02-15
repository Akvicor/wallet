package service

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/wallet"
	"wallet/cmd/app/server/common/types/wallet_partition"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var UserWallet = new(userWalletService)

type userWalletService struct {
	base
}

// FindAllByUIDType 获取所有的钱包, page为nil时不分页
func (w *userWalletService) FindAllByUIDType(page *resp.PageModel, alive bool, preloader *model.PreloaderUserWallet, uid int64, walletType []wallet.Type) (wallets []*model.UserWallet, err error) {
	return repository.UserWallet.FindAllByUIDType(context.Background(), page, alive, preloader, uid, walletType)
}

// FindAllByUIDTypeLike 获取所有的钱包, page为nil时不分页
func (w *userWalletService) FindAllByUIDTypeLike(page *resp.PageModel, alive bool, preloader *model.PreloaderUserWallet, uid int64, walletType []wallet.Type, like string) (wallets []*model.UserWallet, err error) {
	return repository.UserWallet.FindAllByUIDTypeLike(context.Background(), page, alive, preloader, uid, walletType, like)
}

// FindByUID 获取钱包
func (w *userWalletService) FindByUID(alive bool, preloader *model.PreloaderUserWallet, uid, id int64) (wallet *model.UserWallet, err error) {
	return repository.UserWallet.FindByUID(context.Background(), alive, preloader, uid, id)
}

// Create 创建记录
func (w *userWalletService) Create(uid int64, name, description string, walletType wallet.Type, userCardID []int64) (userWallet *model.UserWallet, err error) {
	if !walletType.Valid() {
		return nil, errors.New("wallet type is invalid")
	}
	userWallet = model.NewUserWallet(uid, name, description, walletType)
	var cards []*model.UserCard
	var partitions = make([]*model.UserWalletPartition, 0, len(userCardID))
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户输入的银行卡ID是否属于用户
		if len(userCardID) > 0 {
			cards, err = repository.UserCard.FindAllByID(ctx, nil, true, nil, uid, userCardID)
			if err != nil {
				return err
			}
			if len(cards) != len(userCardID) {
				return errors.New("input user card invalid")
			}
		}
		// 获取序号
		seq, err := repository.UserWallet.GetMaxSequenceByUID(ctx, true, uid, walletType.SameTypes())
		if err != nil {
			seq = 0
		}
		userWallet.Sequence = seq + 1
		// 创建钱包
		err = repository.UserWallet.Create(ctx, userWallet)
		if err != nil {
			return err
		}
		if walletType.IsNormal() && len(cards) > 0 {
			// 创建资金分区
			for i, card := range cards {
				part := model.NewUserWalletPartition(userWallet.ID, card.ID, card.MasterCurrencyID, "Master", "", decimal.NewFromInt(0), decimal.NewFromInt(0), wallet_partition.TypeAverageNormal)
				part.Sequence = int64(i) + 1
				partitions = append(partitions, part)
			}
			err = repository.UserWalletPartition.Creates(ctx, partitions)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return userWallet, nil
}

// Update 更新用户的钱包
func (w *userWalletService) Update(uid, id int64, name, description string, walletType wallet.Type) (err error) {
	if !walletType.Valid() {
		return errors.New("wallet type is invalid")
	}
	userWallet := model.NewUserWallet(uid, name, description, walletType)
	userWallet.ID = id

	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 更新钱包
		err = repository.UserWallet.UpdateByUID(ctx, false, uid, userWallet)
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

// UpdateSequenceByUID 更新序号
func (w *userWalletService) UpdateSequenceByUID(uid, id, targetSequence int64) error {
	err := w.transaction(context.Background(), func(ctx context.Context) error {
		origin, err := repository.UserWallet.FindByUID(ctx, false, nil, uid, id)
		if err != nil {
			return err
		}
		maxSequence, err := repository.UserWallet.GetMaxSequenceByUID(ctx, false, uid, origin.WalletType.SameTypes())
		if err != nil {
			return nil
		}
		if targetSequence < 1 {
			targetSequence = 1
		}
		if targetSequence > maxSequence {
			targetSequence = maxSequence
		}
		err = repository.UserWallet.UpdatesSequenceByUID(ctx, false, uid, origin.WalletType.SameTypes(), origin.Sequence, targetSequence)
		if err != nil {
			return err
		}
		err = repository.UserWallet.UpdateSequenceByUID(ctx, false, uid, id, targetSequence)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// DisableByUID 停用用户钱包
func (w *userWalletService) DisableByUID(uid, id int64) (err error) {
	err = w.transaction(context.Background(), func(ctx context.Context) error {
		// 判断钱包是否属于用户
		err = repository.UserWallet.DetectValidByUID(ctx, uid, []int64{id})
		if err != nil {
			return err
		}
		// 判断划分是否存在关联的划分
		var count int64
		count, err = repository.UserWalletPartition.DetectExistByWalletID(ctx, id)
		if err != nil {
			return err
		}
		if count == 0 {
			// 不存在关联的划分，删除钱包
			err = repository.UserWallet.DeleteByID(ctx, true, uid, id)
		} else {
			// 存在关联的划分，停用钱包
			err = repository.UserWallet.UpdateDisabledByUID(ctx, true, uid, id, time.Now().Unix())
		}
		return err
	})
	return err
}

// EnableByUID 启用用户钱包
func (w *userWalletService) EnableByUID(uid, id int64) error {
	return repository.UserWallet.UpdateDisabledByUID(context.Background(), false, uid, id, 0)
}
