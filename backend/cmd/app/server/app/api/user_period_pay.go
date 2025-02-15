package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"wallet/cmd/app/server/app/dro"
	"wallet/cmd/app/server/app/dto"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

var UserPeriodPay = new(userPeriodPayApi)

type userPeriodPayApi struct{}

// Create 创建周期付费
func (a *userPeriodPayApi) Create(c echo.Context) (err error) {
	input := new(dto.UserPeriodPayCreate)
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
	_, err = service.UserPeriodPay.Create(self.ID, input.Name, input.Description, input.CurrencyID, input.Value, input.PeriodType, input.StartAt, input.NextOfPeriod, input.IntervalOfPeriod, input.ExpirationDate, input.ExpirationTimes)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取用户全部周期付费
func (a *userPeriodPayApi) Find(c echo.Context) (err error) {
	input := new(dto.UserPeriodPayFind)
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
	if input.ID == 0 {
		var wallets = make([]*model.UserPeriodPay, 0)
		if input.Search == "" {
			wallets, err = service.UserPeriodPay.FindAllByUID(&input.PageModel, !input.All, model.NewPreloaderUserPeriodPay().All(), self.ID)
		} else {
			wallets, err = service.UserPeriodPay.FindAllByUIDLike(&input.PageModel, !input.All, model.NewPreloaderUserPeriodPay().All(), self.ID, input.Search)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		return resp.SuccessWithPageData(c, &input.PageModel, wallets)
	} else {
		var userWallet *model.UserPeriodPay
		userWallet, err = service.UserPeriodPay.FindByUID(!input.All, model.NewPreloaderUserPeriodPay().All(), self.ID, input.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		return resp.SuccessWithData(c, userWallet)
	}
}

// Summary 统计用户周期付费
func (a *userPeriodPayApi) Summary(c echo.Context) (err error) {
	input := new(dto.UserPeriodPaySummary)
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

	allPeriod := make([]*model.UserPeriodPay, 0)
	{
		allPeriod, err = service.UserPeriodPay.FindAllByUID(&input.PageModel, !input.All, model.NewPreloaderUserPeriodPay().All(), self.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
	}
	searchPeriod := make([]*model.UserPeriodPay, 0)
	{
		if input.ID == 0 {
			if input.Search == "" {
				searchPeriod, err = service.UserPeriodPay.FindAllByUID(&input.PageModel, !input.All, model.NewPreloaderUserPeriodPay().All(), self.ID)
			} else {
				searchPeriod, err = service.UserPeriodPay.FindAllByUIDLike(&input.PageModel, !input.All, model.NewPreloaderUserPeriodPay().All(), self.ID, input.Search)
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return resp.FailWithMsg(c, resp.Failed, "系统错误")
			}
		} else {
			var userWallet *model.UserPeriodPay
			userWallet, err = service.UserPeriodPay.FindByUID(!input.All, model.NewPreloaderUserPeriodPay().All(), self.ID, input.ID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return resp.FailWithMsg(c, resp.Failed, "系统错误")
			}
			searchPeriod = append(searchPeriod, userWallet)
		}
	}

	data := &dro.UserPeriodPayViewResp{}
	data.ParsePeriodPay(self.ID, allPeriod, searchPeriod)

	return resp.SuccessWithData(c, data)
}

// Update 更新周期付费
func (a *userPeriodPayApi) Update(c echo.Context) (err error) {
	input := new(dto.UserPeriodPayUpdate)
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
	err = service.UserPeriodPay.Update(self.ID, input.ID, input.Name, input.Description, input.CurrencyID, input.Value, input.PeriodType, input.StartAt, input.NextOfPeriod, input.IntervalOfPeriod, input.ExpirationDate, input.ExpirationTimes)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// UpdateNextPeriod 确认周期付费
func (a *userPeriodPayApi) UpdateNextPeriod(c echo.Context) (err error) {
	input := new(dto.UserPeriodPayUpdateNextPeriod)
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
	err = service.UserPeriodPay.UpdatesNextPeriodByID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "确认失败: "+err.Error())
	}
	return resp.Success(c)
}

// Delete 删除周期付费
func (a *userPeriodPayApi) Delete(c echo.Context) (err error) {
	input := new(dto.UserPeriodPayDisable)
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
	err = service.UserPeriodPay.DeleteByUID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "删除失败: "+err.Error())
	}
	return resp.Success(c)
}
