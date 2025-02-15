package dto

import (
	"errors"
	"strings"
)

func (b *CurrencyCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入名称")
	}
	b.EnglishName = strings.TrimSpace(b.EnglishName)
	if len(b.EnglishName) == 0 {
		return errors.New("请输入英文名称")
	}
	b.Code = strings.TrimSpace(b.Code)
	if len(b.Code) == 0 {
		return errors.New("请输入代码")
	}
	b.Symbol = strings.TrimSpace(b.Symbol)
	if len(b.Symbol) == 0 {
		return errors.New("请输入符号")
	}
	return nil
}

func (b *CurrencyFind) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *CurrencyUpdate) Validate() error {
	return b.CurrencyCreate.Validate()
}

func (b *CurrencyDelete) Validate() error {
	return nil
}
