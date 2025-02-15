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

var CardType = new(cardTypeApi)

type cardTypeApi struct{}

// Create 创建银行卡种类
func (a *cardTypeApi) Create(c echo.Context) (err error) {
	input := new(dto.CardTypeCreate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	_, err = service.CardType.Create(input.Name)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取全部银行卡种类
func (a *cardTypeApi) Find(c echo.Context) (err error) {
	input := new(dto.CardTypeFind)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	var cardTypes = make([]*model.CardType, 0)
	if input.Search == "" {
		cardTypes, err = service.CardType.FindAll(&input.PageModel, false, model.NewPreloaderCardType().All())
	} else {
		cardTypes, err = service.CardType.FindAllLike(&input.PageModel, false, model.NewPreloaderCardType().All(), input.Search)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	return resp.SuccessWithPageData(c, &input.PageModel, cardTypes)
}

// Update 更新银行卡种类
func (a *cardTypeApi) Update(c echo.Context) (err error) {
	input := new(dto.CardTypeUpdate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	err = service.CardType.Update(input.ID, input.Name)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// Delete 删除银行卡种类
func (a *cardTypeApi) Delete(c echo.Context) (err error) {
	input := new(dto.CardTypeDelete)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	err = service.CardType.Delete(input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "删除失败: "+err.Error())
	}
	return resp.Success(c)
}
