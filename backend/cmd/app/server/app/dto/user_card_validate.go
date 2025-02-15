package dto

import (
	"errors"
	"slices"
	"strings"
)

func (b *UserCardCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入银行卡名称")
	}
	b.Number = strings.TrimSpace(b.Number)
	if len(b.Number) == 0 {
		return errors.New("请输入银行卡号")
	}
	b.Description = strings.TrimSpace(b.Description)
	b.CVV = strings.TrimSpace(b.CVV)
	b.Password = strings.TrimSpace(b.Password)
	if !slices.Contains(b.CurrencyID, b.MasterCurrencyID) {
		return errors.New("master currency must in currency")
	}
	return nil
}

func (b *UserCardFind) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *UserCardValidRequest) Validate() error {
	b.Key = strings.TrimSpace(b.Key)
	if len(b.Key) == 0 {
		return errors.New("请输入Key")
	}
	return nil
}

func (b *UserCardValidInput) Validate() error {
	b.Key = strings.TrimSpace(b.Key)
	if len(b.Key) == 0 {
		return errors.New("请输入Key")
	}
	b.VerifyCode = strings.TrimSpace(b.VerifyCode)
	if len(b.Key) == 0 {
		return errors.New("请输入验证码")
	}
	return nil
}

func (b *UserCardValidCancel) Validate() error {
	b.Key = strings.TrimSpace(b.Key)
	if len(b.Key) == 0 {
		return errors.New("请输入Key")
	}
	return nil
}

func (b *UserCardUpdate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入银行卡名称")
	}
	b.Number = strings.TrimSpace(b.Number)
	if len(b.Number) == 0 {
		return errors.New("请输入银行卡号")
	}
	b.Description = strings.TrimSpace(b.Description)
	b.CVV = strings.TrimSpace(b.CVV)
	b.Password = strings.TrimSpace(b.Password)
	return nil
}

func (b *UserCardUpdateSequence) Validate() error {
	return nil
}

func (b *UserCardDisable) Validate() error {
	return nil
}

func (b *UserCardEnable) Validate() error {
	return nil
}
