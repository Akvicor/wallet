package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/wallet"
	"wallet/cmd/app/server/model"
)

var UserWallet = new(userWalletRepository)

type userWalletRepository struct {
	base[*model.UserWallet]
}

/**
查找
*/

// FindAll 获取所有的钱包, page为nil时不分页
func (u *userWalletRepository) FindAll(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserWallet) (wallets []*model.UserWallet, err error) {
	return u.paging(page, u.preload(c, alive, preloader).Order("disabled ASC").Order("sequence ASC"))
}

// FindAllByUIDType 获取所有的钱包, page为nil时不分页
func (u *userWalletRepository) FindAllByUIDType(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserWallet, uid int64, walletType []wallet.Type) (wallets []*model.UserWallet, err error) {
	tx := u.preload(c, alive, preloader).Where("uid = ?", uid)
	if len(walletType) != 0 {
		tx = tx.Where("wallet_type IN (?)", walletType)
	}
	wallets, err = u.paging(page, tx.Order("disabled ASC").Order(fmt.Sprintf("CASE WHEN wallet_type IN (%d,%d) THEN %d ELSE wallet_type END ASC", wallet.TypeNormal, wallet.TypeHideBalance, wallet.TypeNormal)).Order("sequence ASC"))
	return wallets, err
}

// FindAllByUIDTypeLike 模糊搜索, page为nil时不分页
func (u *userWalletRepository) FindAllByUIDTypeLike(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserWallet, uid int64, walletType []wallet.Type, like string) (wallets []*model.UserWallet, err error) {
	tx := u.preload(c, alive, preloader).Where("uid = ?", uid)
	if len(walletType) != 0 {
		tx = tx.Where("wallet_type IN (?)", walletType)
	}

	if len(like) > 0 {
		like = "%" + like + "%"
		tx = tx.Where("name LIKE ?", like)
	}

	wallets, err = u.paging(page, tx.Order("disabled ASC").Order(fmt.Sprintf("CASE WHEN wallet_type IN (%d,%d) THEN %d ELSE wallet_type END ASC", wallet.TypeNormal, wallet.TypeHideBalance, wallet.TypeNormal)).Order("sequence ASC"))
	return wallets, err
}

// FindByUID 获取所有的记录, page为nil时不分页
func (u *userWalletRepository) FindByUID(c context.Context, alive bool, preloader *model.PreloaderUserWallet, uid, id int64) (wallet *model.UserWallet, err error) {
	wallet = new(model.UserWallet)
	err = u.preload(c, alive, preloader).Where("uid = ? AND id = ?", uid, id).First(wallet).Error
	return
}

// GetMaxSequenceByUID 获取最大序号
func (u *userWalletRepository) GetMaxSequenceByUID(c context.Context, alive bool, uid int64, walletType []wallet.Type) (seq int64, err error) {
	userWallet := new(model.UserWallet)
	err = u.WrapResultErr(u.alive(c, alive).Where("uid = ?", uid).Where("wallet_type IN (?)", walletType).Order("sequence DESC").First(userWallet))
	if err != nil {
		return 0, err
	}
	return userWallet.Sequence, nil
}

// DetectValidByUID 判断Wallet是否属于某个用户
func (u *userWalletRepository) DetectValidByUID(c context.Context, uid int64, ids []int64) (err error) {
	count := int64(0)
	err = u.dbt(c).Where("uid = ? AND id IN (?)", uid, ids).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("inconsistent quantity")
	}
	return nil
}

/**
创建
*/

// Create 创建记录
func (u *userWalletRepository) Create(c context.Context, wallet *model.UserWallet) error {
	return u.WrapResultErr(u.db(c).Create(wallet))
}

/**
更新
*/

// UpdateByUID 更新记录
func (u *userWalletRepository) UpdateByUID(c context.Context, alive bool, uid int64, wallet *model.UserWallet) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, wallet.ID).Select("*").Omit("sequence", "disabled").Updates(wallet))
}

// UpdatesSequenceByUID 更新范围序号，需配合 UpdateSequenceByUID 使用
func (u *userWalletRepository) UpdatesSequenceByUID(c context.Context, alive bool, uid int64, walletType []wallet.Type, origin, target int64) (err error) {
	if origin == target {
		return nil
	} else if origin < target {
		// update (origin, target] -= 1, origin = target
		err = u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND sequence > ? AND sequence <= ?", uid, origin, target).Where("wallet_type IN (?)", walletType).UpdateColumn("sequence", gorm.Expr("sequence - ?", 1)))
	} else if target < origin {
		// update [target, origin) += 1, origin = target
		err = u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND sequence >= ? AND sequence < ?", uid, target, origin).Where("wallet_type IN (?)", walletType).UpdateColumn("sequence", gorm.Expr("sequence + ?", 1)))
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// UpdateSequenceByUID 更新序号，需配合 UpdatesSequenceByUID 使用
func (u *userWalletRepository) UpdateSequenceByUID(c context.Context, alive bool, uid, id, targetSequence int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, id).UpdateColumn("sequence", targetSequence))
}

// UpdateDisabledByUID 更新停用时间
func (u *userWalletRepository) UpdateDisabledByUID(c context.Context, alive bool, uid, id, timestamp int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, id).UpdateColumn("disabled", timestamp))
}

/**
删除
*/

// DeleteByID 更新停用时间
func (u *userWalletRepository) DeleteByID(c context.Context, alive bool, uid, id int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, id).Delete(&model.UserWallet{}))
}
