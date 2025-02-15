package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"wallet/cmd/app/server/app/dto"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

var UserExchangeRate = new(userExchangeRateApi)

type userExchangeRateApi struct{}

// Create 创建用户货币汇率
func (a *userExchangeRateApi) Create(c echo.Context) (err error) {
	input := new(dto.UserExchangeRateCreate)
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
	_, err = service.UserExchangeRate.Create(self.ID, input.FromCurrencyID, input.ToCurrencyID, input.Rate)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取用户全部货币汇率
func (a *userExchangeRateApi) Find(c echo.Context) (err error) {
	input := new(dto.UserExchangeRateFind)
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
	var userExchangeRate = make([]*model.UserExchangeRate, 0)
	userExchangeRate, err = service.UserExchangeRate.FindAllByUID(&input.PageModel, false, model.NewPreloaderUserExchangeRate().All(), self.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	return resp.SuccessWithPageData(c, &input.PageModel, userExchangeRate)
}

// Update 更新用户货币汇率
func (a *userExchangeRateApi) Update(c echo.Context) (err error) {
	input := new(dto.UserExchangeRateUpdate)
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
	err = service.UserExchangeRate.Update(self.ID, input.ID, input.FromCurrencyID, input.ToCurrencyID, input.Rate)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// Delete 删除用户货币汇率
func (a *userExchangeRateApi) Delete(c echo.Context) (err error) {
	input := new(dto.UserExchangeRateUpdate)
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
	err = service.UserExchangeRate.DeleteWithUID(true, self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "删除失败: "+err.Error())
	}
	return resp.Success(c)
}
