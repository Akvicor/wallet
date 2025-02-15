package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/model"
)

var UserTransactionCategory = new(userTransactionCategoryRepository)

type userTransactionCategoryRepository struct {
	base[*model.UserTransactionCategory]
}

/**
查找
*/

// FindAllByUID 获取所有的记录, page为nil时不分页
func (t *userTransactionCategoryRepository) FindAllByUID(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserTransactionCategory, uid int64) (transactionCategory []*model.UserTransactionCategory, err error) {
	return t.paging(page, t.preload(c, alive, preloader).Where("uid = ?", uid).Order("type ASC").Order("sequence ASC"))
}

// FindAllByUIDLike 模糊搜索, page为nil时不分页
func (t *userTransactionCategoryRepository) FindAllByUIDLike(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserTransactionCategory, uid int64, like string) (transactionCategory []*model.UserTransactionCategory, err error) {
	tx := t.preload(c, alive, preloader).Where("uid = ?", uid)

	if len(like) > 0 {
		like = "%" + like + "%"
		tx = tx.Where("name LIKE ?", like)
	}

	return t.paging(page, tx.Order("type ASC").Order("sequence ASC"))
}

// FindByUID 获取记录
func (t *userTransactionCategoryRepository) FindByUID(c context.Context, alive bool, preloader *model.PreloaderUserTransactionCategory, uid, id int64) (transactionCategory *model.UserTransactionCategory, err error) {
	transactionCategory = new(model.UserTransactionCategory)
	err = t.preload(c, alive, preloader).Where("uid = ? AND id = ?", uid, id).First(transactionCategory).Error
	return
}

// GetMaxSequenceByUID 获取最大序号
func (t *userTransactionCategoryRepository) GetMaxSequenceByUID(c context.Context, alive bool, uid int64, transactionType transaction.Type) (seq int64, err error) {
	category := new(model.UserTransactionCategory)
	err = t.WrapResultErr(t.alive(c, alive).Where("uid = ? AND type = ?", uid, transactionType).Order("sequence DESC").First(category))
	if err != nil {
		return 0, err
	}
	return category.Sequence, nil
}

// DetectValidByUID 判断交易类型是否属于某个用户
func (t *userTransactionCategoryRepository) DetectValidByUID(c context.Context, uid int64, ids []int64) (err error) {
	count := int64(0)
	err = t.dbt(c).Where("uid = ? AND id IN (?)", uid, ids).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("inconsistent quantity")
	}
	return nil
}

// DetectValidByUIDType 判断交易类型是否属于某个用户的某个类型
func (t *userTransactionCategoryRepository) DetectValidByUIDType(c context.Context, uid int64, transactionType transaction.Type, ids []int64) (err error) {
	count := int64(0)
	err = t.dbt(c).Where("uid = ? AND type = ? AND id IN (?)", uid, transactionType, ids).Count(&count).Error
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
func (t *userTransactionCategoryRepository) Create(c context.Context, transactionCategory *model.UserTransactionCategory) error {
	return t.WrapResultErr(t.db(c).Create(transactionCategory))
}

/**
更新
*/

// UpdateByUID 更新记录
func (t *userTransactionCategoryRepository) UpdateByUID(c context.Context, uid int64, transactionCategory *model.UserTransactionCategory) error {
	return t.WrapResultErr(t.db(c).Where("uid = ? AND id = ?", uid, transactionCategory.ID).Select("*").Omit("sequence").Updates(transactionCategory))
}

// UpdatesSequenceByUID 更新范围序号，需配合 UpdateSequenceByUID 使用
func (t *userTransactionCategoryRepository) UpdatesSequenceByUID(c context.Context, alive bool, transactionType transaction.Type, uid int64, origin, target int64) (err error) {
	if origin == target {
		return nil
	} else if origin < target {
		// update (origin, target] -= 1, origin = target
		err = t.WrapResultErr(t.alive(c, alive).Where("uid = ? AND type = ? AND sequence > ? AND sequence <= ?", uid, transactionType, origin, target).UpdateColumn("sequence", gorm.Expr("sequence - ?", 1)))
	} else if target < origin {
		// update [target, origin) += 1, origin = target
		err = t.WrapResultErr(t.alive(c, alive).Where("uid = ? AND type = ? AND sequence >= ? AND sequence < ?", uid, transactionType, target, origin).UpdateColumn("sequence", gorm.Expr("sequence + ?", 1)))
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
func (t *userTransactionCategoryRepository) UpdateSequenceByUID(c context.Context, alive bool, transactionType transaction.Type, uid, id, targetSequence int64) error {
	return t.WrapResultErr(t.alive(c, alive).Where("uid = ? AND type = ? AND id = ?", uid, transactionType, id).UpdateColumn("sequence", targetSequence))
}

/**
删除
*/

// DeleteByUID 删除记录
func (t *userTransactionCategoryRepository) DeleteByUID(c context.Context, uid, id int64) error {
	return t.WrapResultErr(t.db(c).Where("uid = ? AND id = ?", uid, id).Delete(&model.UserTransactionCategory{}))
}
