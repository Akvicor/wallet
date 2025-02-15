package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/model"
)

var UserTransaction = new(userTransactionRepository)

type userTransactionRepository struct {
	base[*model.UserTransaction]
}

/**
查找
*/

// FindAll 获取所有的交易, page为nil时不分页
func (t *userTransactionRepository) FindAll(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserTransaction) (transaction []*model.UserTransaction, err error) {
	transaction, err = t.paging(page, t.preload(c, alive, preloader).Order("checked = 0 DESC").Order("created DESC"))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return transaction, nil
}

// FindAllByUID 获取所有的交易, page为nil时不分页
func (t *userTransactionRepository) FindAllByUID(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserTransaction, uid int64, like string, fromPartitionID, toPartitionID, fromCurrencyID, toCurrencyID, categoryID []int64, transactionType []transaction.Type) (transaction []*model.UserTransaction, err error) {
	tx := t.preload(c, alive, preloader).Where("uid = ?", uid)
	if len(fromPartitionID) != 0 {
		tx = tx.Where("from_partition_id IN ?", fromPartitionID)
	}
	if len(fromCurrencyID) != 0 {
		tx = tx.Where("from_currency_id IN ?", fromCurrencyID)
	}
	if len(toPartitionID) != 0 {
		tx = tx.Where("to_partition_id IN ?", toPartitionID)
	}
	if len(toCurrencyID) != 0 {
		tx = tx.Where("to_currency_id IN ?", toCurrencyID)
	}
	if len(categoryID) != 0 {
		tx = tx.Where("category_id IN ?", categoryID)
	}
	if len(transactionType) != 0 {
		tx = tx.Where("type IN ?", transactionType)
	}
	if len(like) > 0 {
		like = "%" + like + "%"
		tx = tx.Where("description LIKE ?", like)
	}
	transaction, err = t.paging(page, tx.Order("checked = 0 DESC").Order("created DESC"))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return transaction, nil
}

// FindAllByUIDCreated 获取所有的记录
func (t *userTransactionRepository) FindAllByUIDCreated(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserTransaction, uid, start, end int64) (transactions []*model.UserTransaction, err error) {
	return t.paging(page, t.preload(c, alive, preloader).Where("uid = ? AND created >= ? AND created <= ?", uid, start, end).Order("created DESC"))
}

// FindAllByUIDCreatedFromPartition 获取所有的记录
func (t *userTransactionRepository) FindAllByUIDCreatedFromPartition(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserTransaction, uid int64, fromPartitionID []int64, start, end int64) (transactions []*model.UserTransaction, err error) {
	if len(fromPartitionID) == 0 {
		return t.paging(page, t.preload(c, alive, preloader).Where("uid = ? AND created >= ? AND created <= ?", uid, start, end).Order("created DESC"))
	} else {
		return t.paging(page, t.preload(c, alive, preloader).Where("uid = ? AND from_partition_id IN ? AND created >= ? AND created <= ?", uid, fromPartitionID, start, end).Order("created DESC"))
	}
}

// FindAllByUIDUnchecked 获取所有未确认的记录
func (t *userTransactionRepository) FindAllByUIDUnchecked(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserTransaction, uid int64) (transactions []*model.UserTransaction, err error) {
	return t.paging(page, t.preload(c, alive, preloader).Where("uid = ? AND checked = ?", uid, 0).Order("created DESC"))
}

// FindByUID 获取记录
func (t *userTransactionRepository) FindByUID(c context.Context, alive bool, preloader *model.PreloaderUserTransaction, uid, id int64) (transaction *model.UserTransaction, err error) {
	transaction = new(model.UserTransaction)
	err = t.WrapResultErr(t.preload(c, alive, preloader).Where("uid = ? AND id = ?", uid, id).First(transaction))
	return
}

// DetectValidByUID 判断Transaction是否属于某个用户
func (t *userTransactionRepository) DetectValidByUID(c context.Context, uid int64, ids []int64) (err error) {
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

// DetectExistByPartitionID 判断是否存在某个划分的交易
func (t *userTransactionRepository) DetectExistByPartitionID(c context.Context, uid int64, pid int64) (count int64, err error) {
	err = t.dbt(c).Where("uid = ? AND (from_partition_id = ? OR to_partition_id = ?)", uid, pid, pid).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

/**
创建
*/

// Create 创建记录
func (t *userTransactionRepository) Create(c context.Context, transaction *model.UserTransaction) error {
	return t.WrapResultErr(t.db(c).Create(transaction))
}

/**
更新
*/

// UpdateByUID 更新记录
func (t *userTransactionRepository) UpdateByUID(c context.Context, uid, id, categoryID int64, description string, created int64) error {
	return t.WrapResultErr(t.dbm(c).Where("uid = ? AND id = ?", uid, id).Updates(map[string]any{
		"category_id": categoryID,
		"description": description,
		"created":     created,
		"checked":     0,
	}))
}

// UpdateValueByUID 更新余额
func (t *userTransactionRepository) UpdateValueByUID(c context.Context, uid, id int64, value float64) error {
	return t.WrapResultErr(t.dbm(c).Where("uid = ? AND id = ?", uid, id).UpdateColumn("value", value))
}

// UpdateCheckedByUID 更新Checked
func (t *userTransactionRepository) UpdateCheckedByUID(c context.Context, uid, id int64, timestamp int64) error {
	return t.WrapResultErr(t.dbm(c).Where("uid = ? AND id = ?", uid, id).UpdateColumn("checked", timestamp))
}

// DeleteByUID 删除
func (t *userTransactionRepository) DeleteByUID(c context.Context, uid, id int64) error {
	return t.WrapResultErr(t.db(c).Where("uid = ? AND id = ?", uid, id).Delete(&model.UserTransaction{}))
}
