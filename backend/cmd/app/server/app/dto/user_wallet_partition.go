package dto

import (
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/common/types/wallet_partition"
)

type UserWalletPartitionCreate struct {
	WalletID    int64                        `json:"wallet_id" form:"wallet_id" query:"wallet_id"`       // 钱包ID
	CardID      int64                        `json:"card_id" form:"card_id" query:"card_id"`             // 银行卡ID
	CurrencyID  int64                        `json:"currency_id" form:"currency_id" query:"currency_id"` // 货币ID
	Name        string                       `json:"name" form:"name" query:"name"`                      // 划分名字
	Description string                       `json:"description" form:"description" query:"description"` // 划分描述
	Limit       decimal.Decimal              `json:"limit" form:"limit" query:"limit"`                   // 划分资金上限
	Average     wallet_partition.AverageType `json:"average" form:"average" query:"average"`             // 余额均分
}

type UserWalletPartitionUpdate struct {
	Id
	WalletID    int64                        `json:"wallet_id" form:"wallet_id" query:"wallet_id"`       // 钱包ID
	Name        string                       `json:"name" form:"name" query:"name"`                      // 划分名字
	Description string                       `json:"description" form:"description" query:"description"` // 划分描述
	Limit       decimal.Decimal              `json:"limit" form:"limit" query:"limit"`                   // 划分资金上限
	Average     wallet_partition.AverageType `json:"average" form:"average" query:"average"`             // 余额均分
}

type UserWalletPartitionDisable struct {
	WalletID int64 `json:"wallet_id" form:"wallet_id" query:"wallet_id"` // 钱包ID
	Id
}

type UserWalletPartitionUpdateSequence struct {
	WalletID int64 `json:"wallet_id" form:"wallet_id" query:"wallet_id"` // 钱包ID
	Id
	Target int64 `json:"target" form:"target" query:"target"` // 目标序号
}

type UserWalletPartitionEnable struct {
	WalletID int64 `json:"wallet_id" form:"wallet_id" query:"wallet_id"` // 钱包ID
	Id
}
