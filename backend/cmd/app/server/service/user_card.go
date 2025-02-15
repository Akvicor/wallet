package service

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"slices"
	"time"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/repository"
)

var UserCard = new(userCardService)

type userCardService struct {
	base
}

// FindAll 获取所有银行卡
func (u *userCardService) FindAll(page *resp.PageModel, alive bool, preload *model.PreloaderUserCard) (cards []*model.UserCard, err error) {
	return repository.UserCard.FindAll(context.Background(), page, alive, preload)
}

// FindAllByUID 通过UID获取所有银行卡
func (u *userCardService) FindAllByUID(page *resp.PageModel, alive bool, preload *model.PreloaderUserCard, uid int64) (cards []*model.UserCard, err error) {
	return repository.UserCard.FindAllByUID(context.Background(), page, alive, preload, uid)
}

// FindAllByUIDLike 通过UID获取所有银行卡
func (u *userCardService) FindAllByUIDLike(page *resp.PageModel, alive bool, preload *model.PreloaderUserCard, uid int64, like string) (cards []*model.UserCard, err error) {
	return repository.UserCard.FindAllByUIDLike(context.Background(), page, alive, preload, uid, like)
}

// FindByUID 获取所有的银行卡, page为nil时不分页
func (u *userCardService) FindByUID(alive bool, preloader *model.PreloaderUserCard, uid, id int64) (card *model.UserCard, err error) {
	return repository.UserCard.FindByUID(context.Background(), alive, preloader, uid, id)
}

// FindExpiringSoon 获取所有的即将过期的银行卡
func (u *userCardService) FindExpiringSoon(alive bool, preloader *model.PreloaderUserCard, expDate int64) (cards []*model.UserCard, err error) {
	return repository.UserCard.FindAllByExpDate(context.Background(), nil, alive, preloader, expDate, true)
}

// FindAllByStatementClosingDay 获取所有指定账单日的银行卡
func (u *userCardService) FindAllByStatementClosingDay(alive bool, preloader *model.PreloaderUserCard, statementClosingDay int64) (cards []*model.UserCard, err error) {
	return repository.UserCard.FindAllByStatementClosingDay(context.Background(), nil, alive, preloader, statementClosingDay)
}

// FindAllByPaymentDueDay 获取所有指定还款日的银行卡
func (u *userCardService) FindAllByPaymentDueDay(alive bool, preloader *model.PreloaderUserCard, paymentDueDay int64) (cards []*model.UserCard, err error) {
	return repository.UserCard.FindAllByPaymentDueDay(context.Background(), nil, alive, preloader, paymentDueDay)
}

// Create 创建用户银行卡
func (u *userCardService) Create(uid, bid int64, typeID []int64, name, description, number string, expDate int64, cvv string, statementClosingDay, paymentDueDay int64, password string, masterCurrencyID int64, currencyID []int64, limit decimal.Decimal, fee string, hideBalance bool) (card *model.UserCard, err error) {
	card = model.NewUserCard(uid, bid, name, description, number, expDate, cvv, statementClosingDay, paymentDueDay, password, masterCurrencyID, limit, fee, hideBalance)
	cardTypes := make([]*model.UserCardRelTypes, len(typeID))
	cardCurrency := make([]*model.UserCardCurrency, len(currencyID))
	err = u.transaction(context.Background(), func(ctx context.Context) error {
		// 查看主货币是否包含在货币之中
		if !slices.Contains(currencyID, masterCurrencyID) {
			return errors.New("master currency is not in currency")
		}
		// 获取序号
		seq, err := repository.UserCard.GetMaxSequenceByUID(ctx, true, uid)
		if err != nil {
			seq = 0
		}
		card.Sequence = seq + 1
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 创建用户银行卡
		e = repository.UserCard.Create(ctx, card)
		if e != nil {
			return e
		}
		// 绑定卡类型
		exist, e = repository.CardType.ExistById(ctx, true, typeID)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("卡种类不存在")
		}
		for k, v := range typeID {
			cardTypes[k] = model.NewUserCardRelTypes(card.ID, v)
		}
		e = repository.UserCardRelTypes.Creates(ctx, cardTypes)
		if e != nil {
			return e
		}
		// 绑定卡货币种类
		exist, e = repository.Currency.ExistById(ctx, true, currencyID)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("货币种类不存在")
		}
		for k, v := range currencyID {
			cardCurrency[k] = model.NewUserCardCurrency(card.ID, v, decimal.NewFromInt(0))
		}
		e = repository.UserCardCurrency.Creates(ctx, cardCurrency)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return card, nil
}

// UpdateByUID 更新
func (u *userCardService) UpdateByUID(uid, id, bid int64, typeID []int64, name, description, number string, expDate int64, cvv string, statementClosingDay, paymentDueDay int64, password string, masterCurrencyID int64, limit decimal.Decimal, fee string, hideBalance bool) (err error) {
	card := model.NewUserCard(uid, bid, name, description, number, expDate, cvv, statementClosingDay, paymentDueDay, password, masterCurrencyID, limit, fee, hideBalance)
	card.ID = id
	cardTypes := make([]*model.UserCardRelTypes, len(typeID))
	err = u.transaction(context.Background(), func(ctx context.Context) error {
		// 判断用户是否存在且未被停用
		exist, e := repository.User.ExistById(ctx, true, uid)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("用户已被停用")
		}
		// 判断用户输入的银行卡ID是否属于用户
		e = repository.UserCard.DetectValidByUID(ctx, uid, []int64{card.ID})
		if e != nil {
			return e
		}
		// 更新用户银行卡
		e = repository.UserCard.UpdateByUID(ctx, false, uid, card)
		if e != nil {
			return e
		}
		// 查看主货币是否包含在货币之中
		cardCurrencyIDS, e := repository.UserCardCurrency.GetCurrencyIDSByCardID(ctx, true, nil, card.ID)
		if e != nil {
			return e
		}
		if !slices.Contains(cardCurrencyIDS, masterCurrencyID) {
			return errors.New("master currency is not in currency")
		}
		// 绑定卡类型
		exist, e = repository.CardType.ExistById(ctx, true, typeID)
		if !exist || errors.Is(e, gorm.ErrRecordNotFound) {
			return errors.New("卡种类不存在")
		}
		// 解除旧绑定
		e = repository.UserCardRelTypes.DeleteByCardID(ctx, card.ID)
		if e != nil {
			return e
		}
		// 创建新绑定
		for k, v := range typeID {
			cardTypes[k] = model.NewUserCardRelTypes(card.ID, v)
		}
		e = repository.UserCardRelTypes.Creates(ctx, cardTypes)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateWalletByUID 更新银行卡绑定钱包
func (u *userCardService) UpdateWalletByUID(uid, id, walletID int64) error {
	return repository.UserCard.UpdateWalletByUID(context.Background(), true, uid, id, walletID)
}

// UpdateSequenceByUID 更新序号
func (u *userCardService) UpdateSequenceByUID(uid, id, targetSequence int64) error {
	err := u.transaction(context.Background(), func(ctx context.Context) error {
		maxSequence, err := repository.UserCard.GetMaxSequenceByUID(ctx, false, uid)
		if err != nil {
			return nil
		}
		if targetSequence < 1 {
			targetSequence = 1
		}
		if targetSequence > maxSequence {
			targetSequence = maxSequence
		}
		origin, err := repository.UserCard.FindByUID(ctx, false, nil, uid, id)
		if err != nil {
			return err
		}
		err = repository.UserCard.UpdatesSequenceByUID(ctx, false, uid, origin.Sequence, targetSequence)
		if err != nil {
			return err
		}
		err = repository.UserCard.UpdateSequenceByUID(ctx, false, uid, id, targetSequence)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// DisableByUID 停用用户银行卡
func (u *userCardService) DisableByUID(uid, id int64) error {
	return repository.UserCard.UpdateDisabledByUID(context.Background(), true, uid, id, time.Now().Unix())
}

// EnableByUID 启用用户银行卡
func (u *userCardService) EnableByUID(uid, id int64) error {
	return repository.UserCard.UpdateDisabledByUID(context.Background(), false, uid, id, 0)
}
