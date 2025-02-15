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

var UserTransactionCategory = new(userTransactionCategoryApi)

type userTransactionCategoryApi struct{}

// Create 创建交易分类
func (a *userTransactionCategoryApi) Create(c echo.Context) (err error) {
	input := new(dto.UserTransactionCategoryCreate)
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
	_, err = service.UserTransactionCategory.Create(self.ID, input.Type, input.Name, input.Colour)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取全部交易分类
func (a *userTransactionCategoryApi) Find(c echo.Context) (err error) {
	input := new(dto.UserTransactionCategoryFind)
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
	var transactionCategory = make([]*model.UserTransactionCategory, 0)
	if input.Search == "" {
		transactionCategory, err = service.UserTransactionCategory.FindAllByUID(&input.PageModel, false, model.NewPreloaderUserTransactionCategory().All(), self.ID)
	} else {
		transactionCategory, err = service.UserTransactionCategory.FindAllByUIDLike(&input.PageModel, false, model.NewPreloaderUserTransactionCategory().All(), self.ID, input.Search)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	return resp.SuccessWithPageData(c, &input.PageModel, transactionCategory)
}

// Update 更新交易分类
func (a *userTransactionCategoryApi) Update(c echo.Context) (err error) {
	input := new(dto.UserTransactionCategoryUpdate)
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
	err = service.UserTransactionCategory.Update(self.ID, input.ID, input.Type, input.Name, input.Colour)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// UpdateSequence 更新序号
func (a *userTransactionCategoryApi) UpdateSequence(c echo.Context) (err error) {
	input := new(dto.UserTransactionCategoryUpdateSequence)
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
	err = service.UserTransactionCategory.UpdateSequenceByUID(self.ID, input.ID, input.Target)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "序号更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// Delete 删除交易分类
func (a *userTransactionCategoryApi) Delete(c echo.Context) (err error) {
	input := new(dto.UserTransactionCategoryDelete)
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
	err = service.UserTransactionCategory.DeleteByUID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "删除失败: "+err.Error())
	}
	return resp.Success(c)
}
