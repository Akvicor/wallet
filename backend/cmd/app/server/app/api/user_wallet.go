package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"wallet/cmd/app/server/app/dto"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/sign"
	"wallet/cmd/app/server/common/types/wallet"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

var UserWallet = new(userWalletApi)

type userWalletApi struct{}

// Create 创建钱包
func (a *userWalletApi) Create(c echo.Context) (err error) {
	input := new(dto.UserWalletCreate)
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
	_, err = service.UserWallet.Create(self.ID, input.Name, input.Description, input.WalletType, input.CardsID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取用户全部钱包
func (a *userWalletApi) Find(c echo.Context) (err error) {
	input := new(dto.UserWalletFind)
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
		var wallets = make([]*model.UserWallet, 0)
		if input.Search == "" {
			wallets, err = service.UserWallet.FindAllByUIDType(&input.PageModel, !input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, nil)
		} else {
			wallets, err = service.UserWallet.FindAllByUIDTypeLike(&input.PageModel, !input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, nil, input.Search)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		if authorization.Session != nil {
			if authorization.Session.CardCVVValid != sign.SessionCardCVVValid {
				for _, wal := range wallets {
					for _, part := range wal.Partition {
						part.Card.CVV = "000"
					}
				}
			}
			if authorization.Session.CardPasswordValid != sign.SessionCardPasswordValid {
				for _, wal := range wallets {
					for _, part := range wal.Partition {
						part.Card.Password = "000"
					}
				}
			}
		} else {
			for _, wal := range wallets {
				for _, part := range wal.Partition {
					part.Card.CVV = "000"
				}
			}
			for _, wal := range wallets {
				for _, part := range wal.Partition {
					part.Card.Password = "000"
				}
			}
		}
		return resp.SuccessWithPageData(c, &input.PageModel, wallets)
	} else {
		var userWallet *model.UserWallet
		userWallet, err = service.UserWallet.FindByUID(!input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, input.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		return resp.SuccessWithData(c, userWallet)
	}
}

// FindNormal 获取用户全部普通钱包
func (a *userWalletApi) FindNormal(c echo.Context) (err error) {
	input := new(dto.UserWalletFindNormal)
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

	var wallets = make([]*model.UserWallet, 0)
	walletTypes := []wallet.Type{wallet.TypeNormal, wallet.TypeHideBalance}
	if input.Search == "" {
		wallets, err = service.UserWallet.FindAllByUIDType(&input.PageModel, !input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, walletTypes)
	} else {
		wallets, err = service.UserWallet.FindAllByUIDTypeLike(&input.PageModel, !input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, walletTypes, input.Search)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	if authorization.Session != nil {
		if authorization.Session.CardCVVValid != sign.SessionCardCVVValid {
			for _, wal := range wallets {
				for _, part := range wal.Partition {
					part.Card.CVV = "000"
				}
			}
		}
		if authorization.Session.CardPasswordValid != sign.SessionCardPasswordValid {
			for _, wal := range wallets {
				for _, part := range wal.Partition {
					part.Card.Password = "000"
				}
			}
		}
	} else {
		for _, wal := range wallets {
			for _, part := range wal.Partition {
				part.Card.CVV = "000"
			}
		}
		for _, wal := range wallets {
			for _, part := range wal.Partition {
				part.Card.Password = "000"
			}
		}
	}
	return resp.SuccessWithPageData(c, &input.PageModel, wallets)
}

// FindDebt 获取用户全部债务钱包
func (a *userWalletApi) FindDebt(c echo.Context) (err error) {
	input := new(dto.UserWalletFindDebt)
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

	var wallets = make([]*model.UserWallet, 0)
	walletTypes := []wallet.Type{wallet.TypeDebt}
	if input.Search == "" {
		wallets, err = service.UserWallet.FindAllByUIDType(&input.PageModel, !input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, walletTypes)
	} else {
		wallets, err = service.UserWallet.FindAllByUIDTypeLike(&input.PageModel, !input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, walletTypes, input.Search)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	if authorization.Session != nil {
		if authorization.Session.CardCVVValid != sign.SessionCardCVVValid {
			for _, wal := range wallets {
				for _, part := range wal.Partition {
					part.Card.CVV = "000"
				}
			}
		}
		if authorization.Session.CardPasswordValid != sign.SessionCardPasswordValid {
			for _, wal := range wallets {
				for _, part := range wal.Partition {
					part.Card.Password = "000"
				}
			}
		}
	} else {
		for _, wal := range wallets {
			for _, part := range wal.Partition {
				part.Card.CVV = "000"
			}
		}
		for _, wal := range wallets {
			for _, part := range wal.Partition {
				part.Card.Password = "000"
			}
		}
	}
	return resp.SuccessWithPageData(c, &input.PageModel, wallets)
}

// FindWishlist 获取用户全部愿望单钱包
func (a *userWalletApi) FindWishlist(c echo.Context) (err error) {
	input := new(dto.UserWalletFindWishlist)
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

	var wallets = make([]*model.UserWallet, 0)
	walletTypes := []wallet.Type{wallet.TypeWishlist}
	if input.Search == "" {
		wallets, err = service.UserWallet.FindAllByUIDType(&input.PageModel, !input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, walletTypes)
	} else {
		wallets, err = service.UserWallet.FindAllByUIDTypeLike(&input.PageModel, !input.All, model.NewPreloaderUserWallet().All(!input.AllPartition), self.ID, walletTypes, input.Search)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	if authorization.Session != nil {
		if authorization.Session.CardCVVValid != sign.SessionCardCVVValid {
			for _, wal := range wallets {
				for _, part := range wal.Partition {
					part.Card.CVV = "000"
				}
			}
		}
		if authorization.Session.CardPasswordValid != sign.SessionCardPasswordValid {
			for _, wal := range wallets {
				for _, part := range wal.Partition {
					part.Card.Password = "000"
				}
			}
		}
	} else {
		for _, wal := range wallets {
			for _, part := range wal.Partition {
				part.Card.CVV = "000"
			}
		}
		for _, wal := range wallets {
			for _, part := range wal.Partition {
				part.Card.Password = "000"
			}
		}
	}
	return resp.SuccessWithPageData(c, &input.PageModel, wallets)
}

// Update 更新钱包
func (a *userWalletApi) Update(c echo.Context) (err error) {
	input := new(dto.UserWalletUpdate)
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
	err = service.UserWallet.Update(self.ID, input.ID, input.Name, input.Description, input.WalletType)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// UpdateSequence 更新序号
func (a *userWalletApi) UpdateSequence(c echo.Context) (err error) {
	input := new(dto.UserWalletUpdateSequence)
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
	err = service.UserWallet.UpdateSequenceByUID(self.ID, input.ID, input.Target)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "序号更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// Disable 停用钱包
func (a *userWalletApi) Disable(c echo.Context) (err error) {
	input := new(dto.UserWalletDisable)
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
	err = service.UserWallet.DisableByUID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "停用失败: "+err.Error())
	}
	return resp.Success(c)
}

// Enable 启用钱包
func (a *userWalletApi) Enable(c echo.Context) (err error) {
	input := new(dto.UserWalletEnable)
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
	err = service.UserWallet.EnableByUID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "启用失败: "+err.Error())
	}
	return resp.Success(c)
}
