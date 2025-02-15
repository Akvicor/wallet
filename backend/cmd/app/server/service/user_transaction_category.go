package service

import (
	"context"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var UserTransactionCategory = new(userTransactionCategoryService)

type userTransactionCategoryService struct {
	base
}

// FindAllByUID 获取所有记录
func (u *userTransactionCategoryService) FindAllByUID(page *resp.PageModel, alive bool, preload *model.PreloaderUserTransactionCategory, uid int64) (transactionCategory []*model.UserTransactionCategory, err error) {
	return repository.UserTransactionCategory.FindAllByUID(context.Background(), page, alive, preload, uid)
}

// FindAllByUIDLike 获取所有记录
func (u *userTransactionCategoryService) FindAllByUIDLike(page *resp.PageModel, alive bool, preload *model.PreloaderUserTransactionCategory, uid int64, like string) (transactionCategory []*model.UserTransactionCategory, err error) {
	return repository.UserTransactionCategory.FindAllByUIDLike(context.Background(), page, alive, preload, uid, like)
}

// Create 创建记录
func (u *userTransactionCategoryService) Create(uid int64, transactionType transaction.Type, name, colour string) (transactionCategory *model.UserTransactionCategory, err error) {
	transactionCategory = model.NewUserTransactionCategory(uid, transactionType, name, colour)
	err = u.transaction(context.Background(), func(ctx context.Context) error {
		// 获取序号
		seq, err := repository.UserTransactionCategory.GetMaxSequenceByUID(ctx, true, uid, transactionType)
		if err != nil {
			seq = 0
		}
		transactionCategory.Sequence = seq + 1
		err = repository.UserTransactionCategory.Create(ctx, transactionCategory)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactionCategory, nil
}

// Update 更新记录
func (u *userTransactionCategoryService) Update(uid, id int64, transactionType transaction.Type, name, colour string) (err error) {
	transactionCategory := model.NewUserTransactionCategory(uid, transactionType, name, colour)
	transactionCategory.ID = id
	err = repository.UserTransactionCategory.UpdateByUID(context.Background(), uid, transactionCategory)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSequenceByUID 更新序号
func (u *userTransactionCategoryService) UpdateSequenceByUID(uid, id, targetSequence int64) error {
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		origin, err := repository.UserTransactionCategory.FindByUID(ctx, false, nil, uid, id)
		if err != nil {
			return err
		}
		maxSequence, err := repository.UserTransactionCategory.GetMaxSequenceByUID(ctx, false, uid, origin.Type)
		if err != nil {
			return nil
		}
		if targetSequence < 1 {
			targetSequence = 1
		}
		if targetSequence > maxSequence {
			targetSequence = maxSequence
		}
		err = repository.UserTransactionCategory.UpdatesSequenceByUID(ctx, false, origin.Type, uid, origin.Sequence, targetSequence)
		if err != nil {
			return err
		}
		err = repository.UserTransactionCategory.UpdateSequenceByUID(ctx, false, origin.Type, uid, id, targetSequence)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByUID 删除记录
func (u *userTransactionCategoryService) DeleteByUID(uid, id int64) error {
	return repository.UserTransactionCategory.DeleteByUID(context.Background(), uid, id)
}
