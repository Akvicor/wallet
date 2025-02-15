package api

import (
	"github.com/labstack/echo/v4"
	"wallet/cmd/app/server/app/dto"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/app/server/service"
)

var UserWalletPartition = new(userWalletPartitionApi)

type userWalletPartitionApi struct{}

// Create 创建划分
func (a *userWalletPartitionApi) Create(c echo.Context) (err error) {
	input := new(dto.UserWalletPartitionCreate)
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
	_, err = service.UserWalletPartition.Create(self.ID, input.WalletID, input.CardID, input.CurrencyID, input.Name, input.Description, input.Limit, input.Average)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Update 更新划分
func (a *userWalletPartitionApi) Update(c echo.Context) (err error) {
	input := new(dto.UserWalletPartitionUpdate)
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
	err = service.UserWalletPartition.Update(self.ID, input.ID, input.WalletID, input.Name, input.Description, input.Limit, input.Average)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// UpdateSequence 更新序号
func (a *userWalletPartitionApi) UpdateSequence(c echo.Context) (err error) {
	input := new(dto.UserWalletPartitionUpdateSequence)
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
	err = service.UserWalletPartition.UpdateSequenceByUID(self.ID, input.WalletID, input.ID, input.Target)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "停用失败: "+err.Error())
	}
	return resp.Success(c)
}

// Disable 停用划分
func (a *userWalletPartitionApi) Disable(c echo.Context) (err error) {
	input := new(dto.UserWalletPartitionDisable)
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
	err = service.UserWalletPartition.DisableByID(self.ID, input.WalletID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "停用失败: "+err.Error())
	}
	return resp.Success(c)
}

// Enable 启用划分
func (a *userWalletPartitionApi) Enable(c echo.Context) (err error) {
	input := new(dto.UserWalletPartitionEnable)
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
	err = service.UserWalletPartition.EnableByID(self.ID, input.WalletID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "启用失败: "+err.Error())
	}
	return resp.Success(c)
}
