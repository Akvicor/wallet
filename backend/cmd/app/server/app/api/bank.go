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

var Bank = new(bankApi)

type bankApi struct{}

// Create 创建银行
func (a *bankApi) Create(c echo.Context) (err error) {
	input := new(dto.BankCreate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	_, err = service.Bank.Create(input.Name, input.EnglishName, input.Abbr, input.Phone)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取全部银行
func (a *bankApi) Find(c echo.Context) (err error) {
	input := new(dto.BankFind)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	var banks = make([]*model.Bank, 0)
	if input.Search == "" {
		banks, err = service.Bank.FindAll(&input.PageModel, false, model.NewPreloaderBank().All())
	} else {
		banks, err = service.Bank.FindAllLike(&input.PageModel, false, model.NewPreloaderBank().All(), input.Search)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	return resp.SuccessWithPageData(c, &input.PageModel, banks)
}

// Update 更新银行
func (a *bankApi) Update(c echo.Context) (err error) {
	input := new(dto.BankUpdate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	err = service.Bank.Update(input.ID, input.Name, input.EnglishName, input.Abbr, input.Phone)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// Delete 删除银行
func (a *bankApi) Delete(c echo.Context) (err error) {
	input := new(dto.BankDelete)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	err = service.Bank.Delete(input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "删除失败: "+err.Error())
	}
	return resp.Success(c)
}
