package repository

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wallet/cmd/app/server/common/types/wallet_partition"
	"wallet/cmd/app/server/model"
)

var UserWalletPartition = new(userWalletPartitionRepository)

type userWalletPartitionRepository struct {
	base[*model.UserWalletPartition]
}

/**
查找
*/

// FindByWalletID 通过钱包ID获取
func (u *userWalletPartitionRepository) FindByWalletID(c context.Context, alive bool, preloader *model.PreloaderUserWalletPartition, walletID, id int64) (partition *model.UserWalletPartition, err error) {
	partition = new(model.UserWalletPartition)
	err = u.preload(c, alive, preloader).Where("wallet_id = ? AND id = ?", walletID, id).First(partition).Error
	return
}

// FindByID 通过钱包ID获取
func (u *userWalletPartitionRepository) FindByID(c context.Context, alive bool, preloader *model.PreloaderUserWalletPartition, id int64) (partition *model.UserWalletPartition, err error) {
	partition = new(model.UserWalletPartition)
	err = u.preload(c, alive, preloader).Where("id = ?", id).First(partition).Error
	return
}

// GetMaxSequenceByWalletID 获取最大序号
func (u *userWalletPartitionRepository) GetMaxSequenceByWalletID(c context.Context, alive bool, walletID int64) (seq int64, err error) {
	partition := new(model.UserWalletPartition)
	err = u.WrapResultErr(u.alive(c, alive).Where("wallet_id = ?", walletID).Order("sequence DESC").First(partition))
	if err != nil {
		return 0, err
	}
	return partition.Sequence, nil
}

// DetectValidByWalletID 判断WalletPartition是否属于某个钱包
func (u *userWalletPartitionRepository) DetectValidByWalletID(c context.Context, wid int64, ids []int64) (err error) {
	count := int64(0)
	err = u.dbt(c).Where("wallet_id = ? AND id IN (?)", wid, ids).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("inconsistent quantity")
	}
	return nil
}

// DetectExistByWalletID 判断是否存在某个钱包的划分
func (u *userWalletPartitionRepository) DetectExistByWalletID(c context.Context, wid int64) (count int64, err error) {
	err = u.dbt(c).Where("wallet_id = ?", wid).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

/**
创建
*/

// Create 创建记录
func (u *userWalletPartitionRepository) Create(c context.Context, partition *model.UserWalletPartition) error {
	return u.WrapResultErr(u.db(c).Create(partition))
}

// Creates 批量创建记录
func (u *userWalletPartitionRepository) Creates(c context.Context, partitions []*model.UserWalletPartition) error {
	return u.WrapResultErr(u.db(c).Create(partitions))
}

/**
更新
*/

// Update 更新记录
func (u *userWalletPartitionRepository) Update(c context.Context, id int64, name, description string, limit decimal.Decimal, average wallet_partition.AverageType) error {
	return u.WrapResultErr(u.dbm(c).Where("id = ?", id).Updates(map[string]any{
		"name":        name,
		"description": description,
		"limit":       limit,
		"average":     average,
	}))
}

// UpdateBalance 更新余额
func (u *userWalletPartitionRepository) UpdateBalance(c context.Context, alive bool, walletID, id int64, value decimal.Decimal) error {
	return u.WrapResultErr(u.alive(c, alive).Where("wallet_id = ? AND id = ?", walletID, id).UpdateColumn("balance", value))
}

// UpdateBalanceAdd 增加余额
func (u *userWalletPartitionRepository) UpdateBalanceAdd(c context.Context, alive bool, walletID, id int64, value decimal.Decimal) error {
	err := Common.Transaction(c, func(ctx context.Context) error {
		partition := new(model.UserWalletPartition)
		var err error
		err = u.alive(ctx, alive).Clauses(clause.Locking{Strength: "UPDATE"}).Where("wallet_id = ? AND id = ?", walletID, id).First(partition).Error
		if err != nil {
			return err
		}
		partition.Balance = partition.Balance.Add(value)
		return u.alive(ctx, alive).Where("wallet_id = ? AND id = ?", walletID, id).UpdateColumn("balance", partition.Balance).Error
	})
	return err
}

// UpdateBalanceSub 减少余额
func (u *userWalletPartitionRepository) UpdateBalanceSub(c context.Context, alive bool, walletID, id int64, value decimal.Decimal) error {
	err := Common.Transaction(c, func(ctx context.Context) error {
		partition := new(model.UserWalletPartition)
		var err error
		err = u.alive(ctx, alive).Clauses(clause.Locking{Strength: "UPDATE"}).Where("wallet_id = ? AND id = ?", walletID, id).First(partition).Error
		if err != nil {
			return err
		}
		partition.Balance = partition.Balance.Sub(value)
		return u.alive(ctx, alive).Where("wallet_id = ? AND id = ?", walletID, id).UpdateColumn("balance", partition.Balance).Error
	})
	return err
}

// UpdatesSequenceByWalletID 更新范围序号，需配合 UpdateSequenceByWalletID 使用
func (u *userWalletPartitionRepository) UpdatesSequenceByWalletID(c context.Context, alive bool, walletID int64, origin, target int64) (err error) {
	if origin == target {
		return nil
	} else if origin < target {
		// update (origin, target] -= 1, origin = target
		err = u.WrapResultErr(u.alive(c, alive).Where("wallet_id = ? AND sequence > ? AND sequence <= ?", walletID, origin, target).UpdateColumn("sequence", gorm.Expr("sequence - ?", 1)))
	} else if target < origin {
		// update [target, origin) += 1, origin = target
		err = u.WrapResultErr(u.alive(c, alive).Where("wallet_id = ? AND sequence >= ? AND sequence < ?", walletID, target, origin).UpdateColumn("sequence", gorm.Expr("sequence + ?", 1)))
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// UpdateSequenceByWalletID 更新序号，需配合 UpdatesSequenceByWalletID 使用
func (u *userWalletPartitionRepository) UpdateSequenceByWalletID(c context.Context, alive bool, walletID, id, targetSequence int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("wallet_id = ? AND id = ?", walletID, id).UpdateColumn("sequence", targetSequence))
}

// UpdateDisabledByID 通过ID更新停用时间
func (u *userWalletPartitionRepository) UpdateDisabledByID(c context.Context, alive bool, id, timestamp int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("id = ?", id).UpdateColumn("disabled", timestamp))
}

// UpdateDisabledByWallet 通过钱包ID更新停用时间
func (u *userWalletPartitionRepository) UpdateDisabledByWallet(c context.Context, alive bool, walletID, timestamp int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("wallet_id = ?", walletID).UpdateColumn("disabled", timestamp))
}

// UpdateDisabledByCard 通过银行卡ID更新停用时间
func (u *userWalletPartitionRepository) UpdateDisabledByCard(c context.Context, alive bool, cardID, timestamp int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("card_id = ?", cardID).UpdateColumn("disabled", timestamp))
}

/**
删除
*/

// DeleteByID 通过ID删除
func (u *userWalletPartitionRepository) DeleteByID(c context.Context, alive bool, id int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("id = ?", id).Delete(&model.UserWalletPartition{}))
}
