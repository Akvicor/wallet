package dto

import (
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/wallet"
)

type UserWalletCreate struct {
	Name        string      `json:"name" form:"name" query:"name"`                      // 名称
	Description string      `json:"description" form:"description" query:"description"` // 描述
	WalletType  wallet.Type `json:"wallet_type" form:"wallet_type" query:"wallet_type"` // 钱包类型
	CardsID     []int64     `json:"cards_id" form:"cards_id" query:"cards_id"`          // 绑定的银行卡
}

type UserWalletFind struct {
	resp.PageModel
	Id
	All          bool   `json:"all" form:"all" query:"all"`
	AllPartition bool   `json:"all_partition" form:"all_partition" query:"all_partition"`
	Search       string `json:"search" form:"search" query:"search"`
}

type UserWalletFindNormal struct {
	resp.PageModel
	All          bool   `json:"all" form:"all" query:"all"`
	AllPartition bool   `json:"all_partition" form:"all_partition" query:"all_partition"`
	Search       string `json:"search" form:"search" query:"search"`
}

type UserWalletFindDebt struct {
	resp.PageModel
	All          bool   `json:"all" form:"all" query:"all"`
	AllPartition bool   `json:"all_partition" form:"all_partition" query:"all_partition"`
	Search       string `json:"search" form:"search" query:"search"`
}

type UserWalletFindWishlist struct {
	resp.PageModel
	All          bool   `json:"all" form:"all" query:"all"`
	AllPartition bool   `json:"all_partition" form:"all_partition" query:"all_partition"`
	Search       string `json:"search" form:"search" query:"search"`
}

type UserWalletUpdate struct {
	Id
	Name        string      `json:"name" form:"name" query:"name"`                      // 名称
	Description string      `json:"description" form:"description" query:"description"` // 描述
	WalletType  wallet.Type `json:"wallet_type" form:"wallet_type" query:"wallet_type"` // 钱包类型
}

type UserWalletUpdateSequence struct {
	Id
	Target int64 `json:"target" form:"target" query:"target"` // 目标序号
}

type UserWalletDisable struct {
	Id
}

type UserWalletEnable struct {
	Id
}
