package dto

import (
	"errors"
	"strings"
)

func (b *UserTransactionCategoryCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入银行名称")
	}
	b.Colour = strings.TrimSpace(b.Colour)
	if len(b.Colour) == 0 {
		return errors.New("请输入颜色")
	}
	if !b.Type.Valid() {
		return errors.New("请输入正确交易类型")
	}
	return nil
}

func (b *UserTransactionCategoryFind) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *UserTransactionCategoryUpdate) Validate() error {
	return b.UserTransactionCategoryCreate.Validate()
}

func (b *UserTransactionCategoryUpdateSequence) Validate() error {
	return nil
}

func (b *UserTransactionCategoryDelete) Validate() error {
	return nil
}
