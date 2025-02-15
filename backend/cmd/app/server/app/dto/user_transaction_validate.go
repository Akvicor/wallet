package dto

import (
	"errors"
	"github.com/jinzhu/now"
	"strings"
	"time"
)

func (b *UserTransactionCreate) Validate() error {
	if !b.TransactionType.Valid() {
		return errors.New("请输入正确交易类型")
	}
	if b.CategoryID == 0 {
		return errors.New("请输入类型")
	}
	b.Description = strings.TrimSpace(b.Description)
	return nil
}

func (b *UserTransactionFind) Validate() error {
	b.Search = strings.TrimSpace(b.Search)
	return nil
}

func (b *UserTransactionFindRange) Validate() error {
	return nil
}

func (b *UserTransactionUpdate) Validate() error {
	if b.CategoryID == 0 {
		return errors.New("请输入类型")
	}
	b.Description = strings.TrimSpace(b.Description)
	return nil
}

func (b *UserTransactionChecked) Validate() error {
	return nil
}

func (b *UserTransactionDelete) Validate() error {
	return nil
}

func (b *UserTransactionViewDay) Validate() error {
	return nil
}

func (b *UserTransactionViewMonth) Validate() error {
	return nil
}

func (b *UserTransactionViewYear) Validate() error {
	return nil
}

func (b *UserTransactionViewTotal) Validate() error {
	return nil
}

func (b *UserTransactionChart) Validate() error {
	if b.Unit != "day" && b.Unit != "month" && b.Unit != "year" {
		return errors.New("unknown unit")
	}
	if b.Start > b.End {
		return errors.New("wrong range")
	}
	switch b.Unit {
	case "day":
		b.Start = now.With(time.Unix(b.Start, 0)).BeginningOfDay().Unix()
		b.End = now.With(time.Unix(b.End, 0)).EndOfDay().Unix()
	case "month":
		b.Start = now.With(time.Unix(b.Start, 0)).BeginningOfMonth().Unix()
		b.End = now.With(time.Unix(b.End, 0)).EndOfMonth().Unix()
	case "year":
		b.Start = now.With(time.Unix(b.Start, 0)).BeginningOfYear().Unix()
		b.End = now.With(time.Unix(b.End, 0)).EndOfYear().Unix()
	}
	return nil
}

func (b *UserTransactionPieChart) Validate() error {
	if b.Unit != "day" && b.Unit != "month" && b.Unit != "year" {
		return errors.New("unknown unit")
	}
	if b.Start > b.End {
		return errors.New("wrong range")
	}
	switch b.Unit {
	case "day":
		b.Start = now.With(time.Unix(b.Start, 0)).BeginningOfDay().Unix()
		b.End = now.With(time.Unix(b.End, 0)).EndOfDay().Unix()
	case "month":
		b.Start = now.With(time.Unix(b.Start, 0)).BeginningOfMonth().Unix()
		b.End = now.With(time.Unix(b.End, 0)).EndOfMonth().Unix()
	case "year":
		b.Start = now.With(time.Unix(b.Start, 0)).BeginningOfYear().Unix()
		b.End = now.With(time.Unix(b.End, 0)).EndOfYear().Unix()
	}
	return nil
}
