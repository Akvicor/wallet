package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"wallet/cmd/app/server/app/dto"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

var Currency = new(currencyApi)

type currencyApi struct{}

// Create 创建货币种类
func (a *currencyApi) Create(c echo.Context) (err error) {
	input := new(dto.CurrencyCreate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	_, err = service.Currency.Create(input.Name, input.EnglishName, input.Code, input.Symbol)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取全部货币种类
func (a *currencyApi) Find(c echo.Context) (err error) {
	input := new(dto.CurrencyFind)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	var currency = make([]*model.Currency, 0)
	if input.Search == "" {
		currency, err = service.Currency.FindAll(&input.PageModel, false, model.NewPreloaderCurrency().All())
	} else {
		currency, err = service.Currency.FindAllLike(&input.PageModel, false, model.NewPreloaderCurrency().All(), input.Search)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	return resp.SuccessWithPageData(c, &input.PageModel, currency)
}

// Update 更新银行
func (a *currencyApi) Update(c echo.Context) (err error) {
	input := new(dto.CurrencyUpdate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	err = service.Currency.Update(input.ID, input.Name, input.EnglishName, input.Code, input.Symbol)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// Delete 删除货币种类
func (a *currencyApi) Delete(c echo.Context) (err error) {
	input := new(dto.CurrencyDelete)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	err = service.Currency.Delete(input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "删除失败: "+err.Error())
	}
	return resp.Success(c)
}
