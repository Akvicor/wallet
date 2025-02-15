package dto

import (
	"github.com/shopspring/decimal"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/transaction"
)

type UserTransactionCreate struct {
	FromPartitionID int64            `json:"from_partition_id" form:"from_partition_id" query:"from_partition_id"` // 交易源划分
	ToPartitionID   int64            `json:"to_partition_id" form:"to_partition_id" query:"to_partition_id"`       // 交易目标划分
	CurrencyID      int64            `json:"currency_id" form:"currency_id" query:"currency_id"`                   // 交易货币
	CategoryID      int64            `json:"category_id" form:"category_id" query:"category_id"`                   // 交易分类
	TransactionType transaction.Type `json:"transaction_type" form:"transaction_type" query:"transaction_type"`    // 交易类型
	Description     string           `json:"description" form:"description" query:"description"`                   // 描述补充信息
	FromValue       decimal.Decimal  `json:"from_value" form:"from_value" query:"from_value"`                      // 交易源金额
	ToValue         decimal.Decimal  `json:"to_value" form:"to_value" query:"to_value"`                            // 交易目标金额
	Created         int64            `json:"created" form:"created" query:"created"`                               // 交易发生时间
}

type UserTransactionFind struct {
	resp.PageModel
	Id
	All             bool               `json:"all" form:"all" query:"all"`
	Search          string             `json:"search" form:"search" query:"search"`
	FromPartitionID []int64            `json:"from_partition_id" form:"from_partition_id" query:"from_partition_id"` // 交易源划分
	ToPartitionID   []int64            `json:"to_partition_id" form:"to_partition_id" query:"to_partition_id"`       // 交易目标划分
	FromCurrencyID  []int64            `json:"from_currency_id" form:"from_currency_id" query:"from_currency_id"`    // 交易货币
	ToCurrencyID    []int64            `json:"to_currency_id" form:"to_currency_id" query:"to_currency_id"`          // 交易货币
	CategoryID      []int64            `json:"category_id" form:"category_id" query:"category_id"`                   // 交易分类
	TransactionType []transaction.Type `json:"transaction_type" form:"transaction_type" query:"transaction_type"`    // 交易类型
}

type UserTransactionFindRange struct {
	Start int64 `json:"start" form:"start" query:"start"`
	End   int64 `json:"end" form:"end" query:"end"`
	All   bool  `json:"all" form:"all" query:"all"`
}

type UserTransactionUpdate struct {
	Id
	CategoryID  int64  `json:"category_id" form:"category_id" query:"category_id"` // 交易分类
	Description string `json:"description" form:"description" query:"description"` // 描述补充信息
	Created     int64  `json:"created" form:"created" query:"created"`             // 交易发生时间
}

type UserTransactionChecked struct {
	Id
}

type UserTransactionDelete struct {
	Id
}

type UserTransactionViewDay struct {
	Target    int64 `json:"target" form:"target" query:"target"`          // 目标日期
	Direction int   `json:"direction" form:"direction" query:"direction"` // >0 目标日期及后面的日期; <0 目标日期及前面的日期
}

type UserTransactionViewMonth struct {
	Target    int64 `json:"target" form:"target" query:"target"`          // 目标日期
	Direction int   `json:"direction" form:"direction" query:"direction"` // >0 目标日期及后面的日期; <0 目标日期及前面的日期
}

type UserTransactionViewYear struct {
	Target    int64 `json:"target" form:"target" query:"target"`          // 目标日期
	Direction int   `json:"direction" form:"direction" query:"direction"` // >0 目标日期及后面的日期; <0 目标日期及前面的日期
}

type UserTransactionViewTotal struct {
}

type UserTransactionChart struct {
	FromPartitionID []int64 `json:"from_partition_id" form:"from_partition_id" query:"from_partition_id"` // 交易源划分
	Unit            string  `json:"unit" form:"unit" query:"unit"`                                        // 统计单位: day, month, year
	Start           int64   `json:"start" form:"start" query:"start"`                                     // 开始时间
	End             int64   `json:"end" form:"end" query:"end"`                                           // 结束时间
}

type UserTransactionPieChart struct {
	FromPartitionID []int64 `json:"from_partition_id" form:"from_partition_id" query:"from_partition_id"` // 交易源划分
	Unit            string  `json:"unit" form:"unit" query:"unit"`                                        // 统计单位: day, month, year
	Start           int64   `json:"start" form:"start" query:"start"`                                     // 开始时间
	End             int64   `json:"end" form:"end" query:"end"`                                           // 结束时间
}
