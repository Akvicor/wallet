package api

import (
	"errors"
	"fmt"
	"github.com/Akvicor/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"time"
	"wallet/cmd/app/server/app/dto"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/sign"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/app/server/global/mail"
	"wallet/cmd/app/server/global/sms"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

var UserCard = new(userCardApi)

type userCardApi struct{}

// Create 创建银行卡
func (a *userCardApi) Create(c echo.Context) (err error) {
	input := new(dto.UserCardCreate)
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
	_, err = service.UserCard.Create(self.ID, input.BankID, input.TypeID, input.Name, input.Description, input.Number, input.ExpDate, input.CVV, input.StatementClosingDay, input.PaymentDueDay, input.Password, input.MasterCurrencyID, input.CurrencyID, input.Limit, input.Fee, input.HideBalance)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取全部银行卡
func (a *userCardApi) Find(c echo.Context) (err error) {
	input := new(dto.UserCardFind)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	authorization := auth.GetAuthorization(c)
	if authorization == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	if input.ID == 0 {
		var cards = make([]*model.UserCard, 0)
		if input.Search == "" {
			cards, err = service.UserCard.FindAllByUID(&input.PageModel, !input.All, model.NewPreloaderUserCard().All(), self.ID)
		} else {
			cards, err = service.UserCard.FindAllByUIDLike(&input.PageModel, !input.All, model.NewPreloaderUserCard().All(), self.ID, input.Search)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		if authorization.Session != nil {
			if authorization.Session.CardCVVValid != sign.SessionCardCVVValid {
				for _, card := range cards {
					card.CVV = "000"
				}
			}
			if authorization.Session.CardPasswordValid != sign.SessionCardPasswordValid {
				for _, card := range cards {
					card.Password = "000000"
				}
			}
		} else {
			for _, card := range cards {
				card.CVV = "000"
			}
			for _, card := range cards {
				card.Password = "000000"
			}
		}
		return resp.SuccessWithPageData(c, &input.PageModel, cards)
	} else {
		var card *model.UserCard
		card, err = service.UserCard.FindByUID(!input.All, model.NewPreloaderUserCard().All(), self.ID, input.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		if authorization.Session != nil {
			if authorization.Session.CardCVVValid != sign.SessionCardCVVValid {
				card.CVV = "000"
			}
			if authorization.Session.CardPasswordValid != sign.SessionCardPasswordValid {
				card.Password = "000000"
			}
		} else {
			card.CVV = "000"
			card.Password = "000000"
		}
		return resp.SuccessWithData(c, card)
	}
}

func (a *userCardApi) ValidRequest(c echo.Context) (err error) {
	input := new(dto.UserCardValidRequest)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	authorization := auth.GetAuthorization(c)
	if authorization == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	if authorization.Session == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "无效请求")
	}
	if authorization.User == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "无效请求")
	}
	switch input.Method {
	case "phone":
		if len(authorization.User.Phone) == 0 {
			return resp.FailWithMsg(c, resp.UnAuthorized, "请填写手机号")
		}
	case "mail":
		if len(authorization.User.Mail) == 0 {
			return resp.FailWithMsg(c, resp.UnAuthorized, "请填写邮箱")
		}
	default:
		return resp.FailWithMsg(c, resp.UnAuthorized, "无效请求")
	}

	switch input.Key {
	case sign.SessionCardCVVKey:
		if authorization.Session.CardCVVValidLastRequest > time.Now().Add(-60*time.Second).Unix() {
			return resp.FailWithMsg(c, resp.UnAuthorized, "重复请求")
		}
		authorization.Session.CardCVVValid = util.RandomString(6, util.RandomDigit)
		authorization.Session.CardCVVValidLastRequest = time.Now().Unix()
		content := fmt.Sprintf("您请求获取银行卡CVV的验证码为(15分钟内有效): %s", authorization.Session.CardCVVValid)
		if input.Method == "phone" {
			sms.Send(authorization.User.Phone, content)
		} else if input.Method == "mail" {
			err = mail.Send(authorization.User.Mail, "银行卡CVV请求验证", mail.Plain, content)
		}
	case sign.SessionCardPasswordKey:
		if authorization.Session.CardPasswordValidLastRequest > time.Now().Add(-60*time.Second).Unix() {
			return resp.FailWithMsg(c, resp.UnAuthorized, "重复请求")
		}
		authorization.Session.CardPasswordValid = util.RandomString(6, util.RandomDigit)
		authorization.Session.CardPasswordValidLastRequest = time.Now().Unix()
		content := fmt.Sprintf("您请求获取银行卡密码的验证码为(15分钟内有效): %s", authorization.Session.CardPasswordValid)
		if input.Method == "phone" {
			sms.Send(authorization.User.Phone, content)
		} else if input.Method == "mail" {
			err = mail.Send(authorization.User.Mail, "银行卡密码请求验证", mail.Plain, content)
		}
	default:
		return resp.FailWithMsg(c, resp.UnAuthorized, "无效请求")
	}
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, err.Error())
	}
	return resp.Success(c)
}

func (a *userCardApi) ValidInput(c echo.Context) (err error) {
	input := new(dto.UserCardValidInput)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	authorization := auth.GetAuthorization(c)
	if authorization == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	if authorization.Session == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "无效请求")
	}

	switch input.Key {
	case sign.SessionCardCVVKey:
		if authorization.Session.CardCVVValid != input.VerifyCode {
			return resp.FailWithMsg(c, resp.UnAuthorized, "验证失败")
		}
		if authorization.Session.CardCVVValidLastRequest < time.Now().Add(-15*time.Minute).Unix() {
			return resp.FailWithMsg(c, resp.UnAuthorized, "验证超时")
		}
		authorization.Session.CardCVVValid = sign.SessionCardCVVValid
	case sign.SessionCardPasswordKey:
		if authorization.Session.CardPasswordValid != input.VerifyCode {
			return resp.FailWithMsg(c, resp.UnAuthorized, "验证失败")
		}
		if authorization.Session.CardPasswordValidLastRequest < time.Now().Add(-15*time.Minute).Unix() {
			return resp.FailWithMsg(c, resp.UnAuthorized, "验证超时")
		}
		authorization.Session.CardPasswordValid = sign.SessionCardPasswordValid
	default:
		return resp.FailWithMsg(c, resp.UnAuthorized, "无效请求")
	}
	return resp.Success(c)
}

func (a *userCardApi) ValidCancel(c echo.Context) (err error) {
	input := new(dto.UserCardValidCancel)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	authorization := auth.GetAuthorization(c)
	if authorization == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	if authorization.Session == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "无效请求")
	}

	switch input.Key {
	case sign.SessionCardCVVKey:
		authorization.Session.CardCVVValid = sign.SessionEmpty
	case sign.SessionCardPasswordKey:
		authorization.Session.CardPasswordValid = sign.SessionEmpty
	default:
		return resp.FailWithMsg(c, resp.UnAuthorized, "无效请求")
	}
	return resp.Success(c)
}

// Update 更新银行卡
func (a *userCardApi) Update(c echo.Context) (err error) {
	input := new(dto.UserCardUpdate)
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
	err = service.UserCard.UpdateByUID(self.ID, input.ID, input.BankID, input.TypeID, input.Name, input.Description, input.Number, input.ExpDate, input.CVV, input.StatementClosingDay, input.PaymentDueDay, input.Password, input.MasterCurrencyID, input.Limit, input.Fee, input.HideBalance)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// UpdateSequence 更新序号
func (a *userCardApi) UpdateSequence(c echo.Context) (err error) {
	input := new(dto.UserCardUpdateSequence)
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
	err = service.UserCard.UpdateSequenceByUID(self.ID, input.ID, input.Target)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "序号更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// Disable 停用银行卡
func (a *userCardApi) Disable(c echo.Context) (err error) {
	input := new(dto.UserCardDisable)
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
	err = service.UserCard.DisableByUID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "停用失败: "+err.Error())
	}
	return resp.Success(c)
}

// Enable 启用银行卡
func (a *userCardApi) Enable(c echo.Context) (err error) {
	input := new(dto.UserCardEnable)
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
	err = service.UserCard.EnableByUID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "启用失败: "+err.Error())
	}
	return resp.Success(c)
}
