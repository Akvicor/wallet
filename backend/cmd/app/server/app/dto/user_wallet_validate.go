package dto

import (
	"errors"
	"strings"
)

func (b *UserWalletCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入钱包名称")
	}
	b.Description = strings.TrimSpace(b.Description)
	if !b.WalletType.Valid() {
		return errors.New("请输入正确钱包类型")
	}
	return nil
}

func (b *UserWalletFind) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *UserWalletFindNormal) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *UserWalletFindDebt) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *UserWalletFindWishlist) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *UserWalletUpdate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入钱包名称")
	}
	b.Description = strings.TrimSpace(b.Description)
	if !b.WalletType.Valid() {
		return errors.New("请输入正确钱包类型")
	}
	return nil
}

func (b *UserWalletUpdateSequence) Validate() error {
	return nil
}

func (b *UserWalletDisable) Validate() error {
	return nil
}

func (b *UserWalletEnable) Validate() error {
	return nil
}
