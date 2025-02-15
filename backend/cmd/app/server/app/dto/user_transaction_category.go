package dto

import (
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/transaction"
)

type UserTransactionCategoryCreate struct {
	Type   transaction.Type `json:"type" form:"type" query:"type"`       // 交易类型
	Name   string           `json:"name" form:"name" query:"name"`       // 分类名称
	Colour string           `json:"colour" form:"colour" query:"colour"` // 分类颜色
}

type UserTransactionCategoryFind struct {
	resp.PageModel
	Search string `json:"search" form:"search" query:"search"`
}

type UserTransactionCategoryUpdate struct {
	Id
	UserTransactionCategoryCreate
}

type UserTransactionCategoryUpdateSequence struct {
	Id
	Target int64 `json:"target" form:"target" query:"target"` // 目标序号
}

type UserTransactionCategoryDelete struct {
	Id
}
