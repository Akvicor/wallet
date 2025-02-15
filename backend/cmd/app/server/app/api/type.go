package api

import (
	"github.com/labstack/echo/v4"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/period"
	"wallet/cmd/app/server/common/types/role"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/common/types/wallet"
	"wallet/cmd/app/server/common/types/wallet_partition"
)

var Type = new(typeApi)

type typeApi struct{}

func (a *typeApi) RoleType(c echo.Context) (err error) {
	return resp.SuccessWithData(c, role.AllType)
}

func (a *typeApi) PeriodType(c echo.Context) (err error) {
	return resp.SuccessWithData(c, period.AllType)
}

func (a *typeApi) TransactionType(c echo.Context) (err error) {
	return resp.SuccessWithData(c, transaction.AllType)
}

func (a *typeApi) WalletType(c echo.Context) (err error) {
	return resp.SuccessWithData(c, wallet.AllType)
}

func (a *typeApi) WalletPartitionAverageType(c echo.Context) (err error) {
	return resp.SuccessWithData(c, wallet_partition.AllAverageType)
}
