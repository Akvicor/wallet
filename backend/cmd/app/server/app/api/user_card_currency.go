package api

import (
	"github.com/labstack/echo/v4"
	"wallet/cmd/app/server/app/dto"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/app/server/service"
)

var UserCardCurrency = new(userCardCurrencyApi)

type userCardCurrencyApi struct{}

// Create 创建
func (a *userCardCurrencyApi) Create(c echo.Context) (err error) {
	input := new(dto.UserCardCurrencyCreate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	_, err = service.UserCardCurrency.Create(self.ID, input.UserCardID, input.CurrencyID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// UpdateBalance 更新
func (a *userCardCurrencyApi) UpdateBalance(c echo.Context) (err error) {
	input := new(dto.UserCardCurrencyUpdateBalance)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}

	return resp.FailWithMsg(c, resp.Failed, "接口停用")
	// err = input.Validate()
	// if err != nil {
	// 	return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	// }
	// self := auth.GetUser(c)
	// if self == nil {
	// 	return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	// }
	// err = service.UserCardCurrency.UpdateBalance(self.ID, input.ID, input.Balance)
	// if err != nil {
	// 	return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	// }
	// return resp.Success(c)
}

// Delete 删除
func (a *userCardCurrencyApi) Delete(c echo.Context) (err error) {
	input := new(dto.UserCardCurrencyDelete)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	err = service.UserCardCurrency.DeleteByID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "删除失败: "+err.Error())
	}
	return resp.Success(c)
}
