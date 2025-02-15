package dto

import (
	"errors"
)

func (b *UserExchangeRateCreate) Validate() error {
	if b.FromCurrencyID == 0 {
		return errors.New("请输入源货币")
	}
	if b.ToCurrencyID == 0 {
		return errors.New("请输入目标货币")
	}
	return nil
}

func (b *UserExchangeRateFind) Validate() error {
	return nil
}

func (b *UserExchangeRateUpdate) Validate() error {
	return b.UserExchangeRateCreate.Validate()
}

func (b *UserExchangeRateDelete) Validate() error {
	return nil
}
