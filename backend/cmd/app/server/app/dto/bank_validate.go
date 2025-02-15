package dto

import (
	"errors"
	"strings"
)

func (b *BankCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入名称")
	}
	b.EnglishName = strings.TrimSpace(b.EnglishName)
	if len(b.EnglishName) == 0 {
		return errors.New("请输入英文名称")
	}
	b.Abbr = strings.TrimSpace(b.Abbr)
	if len(b.Abbr) == 0 {
		return errors.New("请输入缩写")
	}
	b.Phone = strings.TrimSpace(b.Phone)
	return nil
}

func (b *BankFind) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *BankUpdate) Validate() error {
	return b.BankCreate.Validate()
}

func (b *BankDelete) Validate() error {
	return nil
}
