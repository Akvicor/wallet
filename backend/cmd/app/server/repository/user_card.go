package repository

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/global/encrypt_card"
	"wallet/cmd/app/server/model"
)

var UserCard = new(userCardRepository)

type userCardRepository struct {
	base[*model.UserCard]
}

/**
查找
*/

// FindAll 获取所有的记录, page为nil时不分页
func (u *userCardRepository) FindAll(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserCard) (cards []*model.UserCard, err error) {
	cards, err = u.paging(page, u.preload(c, alive, preloader).Order("disabled ASC").Order("sequence ASC"))
	if err != nil {
		return
	}
	_ = encrypt_card.DecryptCards(cards)
	return cards, nil
}

// FindAllByUID 获取所有的记录, page为nil时不分页
func (u *userCardRepository) FindAllByUID(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserCard, uid int64) (cards []*model.UserCard, err error) {
	cards, err = u.paging(page, u.preload(c, alive, preloader).Where("uid = ?", uid).Order("disabled ASC").Order("sequence ASC"))
	if err != nil {
		return
	}
	_ = encrypt_card.DecryptCards(cards)
	return cards, nil
}

// FindAllByUIDLike 模糊搜索, page为nil时不分页
func (u *userCardRepository) FindAllByUIDLike(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserCard, uid int64, like string) (cards []*model.UserCard, err error) {
	tx := u.preload(c, alive, preloader).Where("uid = ?", uid)

	if len(like) > 0 {
		like = "%" + like + "%"
		tx = tx.Where("name LIKE ? OR number LIKE ?", like, like)
	}

	cards, err = u.paging(page, tx.Order("disabled ASC").Order("sequence ASC"))
	if err != nil {
		return
	}
	_ = encrypt_card.DecryptCards(cards)
	return cards, nil
}

// FindAllByID 获取所有的记录, page为nil时不分页
func (u *userCardRepository) FindAllByID(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserCard, uid int64, ids []int64) (cards []*model.UserCard, err error) {
	cards, err = u.paging(page, u.preload(c, alive, preloader).Where("uid = ? AND id IN (?)", uid, ids).Order("disabled ASC").Order("sequence ASC"))
	if err != nil {
		return
	}
	_ = encrypt_card.DecryptCards(cards)
	return cards, nil
}

// FindAllByExpDate 获取所有的记录, page为nil时不分页
func (u *userCardRepository) FindAllByExpDate(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserCard, expDate int64, le bool) (cards []*model.UserCard, err error) {
	if le {
		cards, err = u.paging(page, u.preload(c, alive, preloader).Where("exp_date <= ?", expDate))
	} else {
		cards, err = u.paging(page, u.preload(c, alive, preloader).Where("exp_date >= ?", expDate))
	}
	if err != nil {
		return
	}
	_ = encrypt_card.DecryptCards(cards)
	return cards, nil
}

// FindAllByStatementClosingDay 获取所有的记录, page为nil时不分页
func (u *userCardRepository) FindAllByStatementClosingDay(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserCard, statementClosingDay int64) (cards []*model.UserCard, err error) {
	cards, err = u.paging(page, u.preload(c, alive, preloader).Where("statement_closing_day = ?", statementClosingDay))
	if err != nil {
		return
	}
	_ = encrypt_card.DecryptCards(cards)
	return cards, nil
}

// FindAllByPaymentDueDay 获取所有的记录, page为nil时不分页
func (u *userCardRepository) FindAllByPaymentDueDay(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderUserCard, paymentDueDay int64) (cards []*model.UserCard, err error) {
	cards, err = u.paging(page, u.preload(c, alive, preloader).Where("payment_due_day = ?", paymentDueDay))
	if err != nil {
		return
	}
	_ = encrypt_card.DecryptCards(cards)
	return cards, nil
}

// FindByUID 获取所有的记录, page为nil时不分页
func (u *userCardRepository) FindByUID(c context.Context, alive bool, preloader *model.PreloaderUserCard, uid, id int64) (card *model.UserCard, err error) {
	card = new(model.UserCard)
	err = u.preload(c, alive, preloader).Where("uid = ? AND id = ?", uid, id).First(card).Error
	if err != nil {
		return
	}
	_ = encrypt_card.DecryptCard(card)
	return
}

// GetMaxSequenceByUID 获取最大序号
func (u *userCardRepository) GetMaxSequenceByUID(c context.Context, alive bool, uid int64) (seq int64, err error) {
	card := new(model.UserCard)
	err = u.WrapResultErr(u.alive(c, alive).Where("uid = ?", uid).Order("sequence DESC").First(card))
	if err != nil {
		return 0, err
	}
	return card.Sequence, nil
}

// DetectValidByUID 判断UserCard是否属于某个用户
func (u *userCardRepository) DetectValidByUID(c context.Context, uid int64, ids []int64) (err error) {
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
func (u *userCardRepository) Create(c context.Context, card *model.UserCard) error {
	err := encrypt_card.EncryptCard(card)
	if err != nil {
		return err
	}
	return u.WrapResultErr(u.db(c).Create(card))
}

/**
更新
*/

// UpdateByUID 更新记录
func (u *userCardRepository) UpdateByUID(c context.Context, alive bool, uid int64, card *model.UserCard) (err error) {
	tx := u.alive(c, alive).Where("uid = ? AND id = ?", uid, card.ID).Select("*")
	omits := []string{"balance", "sequence", "disabled"}
	if card.CVV == "---" {
		omits = append(omits, "cvv")
	}
	if card.Password == "------" {
		omits = append(omits, "password")
	}
	err = encrypt_card.EncryptCard(card)
	if err != nil {
		return err
	}
	return u.WrapResultErr(tx.Omit(omits...).Updates(card))
}

// UpdateWalletByUID 更新绑定钱包
func (u *userCardRepository) UpdateWalletByUID(c context.Context, alive bool, uid, id, walletID int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, id).UpdateColumn("wallet_id", walletID))
}

// UpdatesSequenceByUID 更新范围序号，需配合 UpdateSequenceByUID 使用
func (u *userCardRepository) UpdatesSequenceByUID(c context.Context, alive bool, uid int64, origin, target int64) (err error) {
	if origin == target {
		return nil
	} else if origin < target {
		// update (origin, target] -= 1, origin = target
		err = u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND sequence > ? AND sequence <= ?", uid, origin, target).UpdateColumn("sequence", gorm.Expr("sequence - ?", 1)))
	} else if target < origin {
		// update [target, origin) += 1, origin = target
		err = u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND sequence >= ? AND sequence < ?", uid, target, origin).UpdateColumn("sequence", gorm.Expr("sequence + ?", 1)))
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
func (u *userCardRepository) UpdateSequenceByUID(c context.Context, alive bool, uid, id, targetSequence int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, id).UpdateColumn("sequence", targetSequence))
}

// UpdateBalanceByUID 更新余额
func (u *userCardRepository) UpdateBalanceByUID(c context.Context, alive bool, uid, id int64, value decimal.Decimal) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, id).UpdateColumn("balance", value))
}

// UpdateDisabledByUID 更新停用时间
func (u *userCardRepository) UpdateDisabledByUID(c context.Context, alive bool, uid, id, timestamp int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ? AND id = ?", uid, id).UpdateColumn("disabled", timestamp))
}
