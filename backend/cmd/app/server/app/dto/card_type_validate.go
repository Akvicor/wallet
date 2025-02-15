package dto

import (
	"errors"
	"strings"
)

func (b *CardTypeCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入名称")
	}
	return nil
}

func (b *CardTypeFind) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *CardTypeUpdate) Validate() error {
	return b.CardTypeCreate.Validate()
}

func (b *CardTypeDelete) Validate() error {
	return nil
}
