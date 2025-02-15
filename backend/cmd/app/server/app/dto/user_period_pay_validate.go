package dto

import (
	"errors"
	"strings"
)

func (b *UserPeriodPayCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入名称")
	}
	b.Description = strings.TrimSpace(b.Description)
	if !b.PeriodType.Valid() {
		return errors.New("请输入正确周期类型")
	}
	return nil
}

func (b *UserPeriodPayFind) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *UserPeriodPayUpdate) Validate() error {
	return b.UserPeriodPayCreate.Validate()
}

func (b *UserPeriodPayUpdateNextPeriod) Validate() error {
	return nil
}

func (b *UserPeriodPayDisable) Validate() error {
	return nil
}

func (b *UserPeriodPayEnable) Validate() error {
	return nil
}
