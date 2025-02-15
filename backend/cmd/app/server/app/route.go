package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"wallet/cmd/app/server/app/api"
	"wallet/cmd/app/server/app/mw"
)

func setupRoutes(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodDelete},
	}))
	e.Use(mw.Error)
	e.Use(middleware.Gzip())

	publicGroup := e.Group("")
	{
		web := getFS()
		publicGroup.FileFS("/", "index.html", web)
		publicGroup.FileFS("/*", "index.html", web)
		publicGroup.FileFS("/index.html", "/index.html", web)
		publicGroup.FileFS("/favicon.ico", "favicon.ico", web)
		publicGroup.FileFS("/manifest.json", "/manifest.json", web)
		publicGroup.FileFS("/asset-manifest.json", "/asset-manifest.json", web)
		publicGroup.GET("/static/*", WrapHandler(http.FileServer(http.FS(web))))
	}

	apiGroup := e.Group("/api")
	// 管理员
	apiGroupAdmin := apiGroup.Group("")
	apiGroupAdmin.Use(mw.AuthAdmin)
	// 管理员、普通用户
	apiGroupUser := apiGroup.Group("")
	apiGroupUser.Use(mw.AuthUser)
	// 管理员、普通用户、浏览者
	apiGroupViewer := apiGroup.Group("")
	apiGroupViewer.Use(mw.AuthViewer)
	// 已登录用户
	apiGroupAuth := apiGroup.Group("")
	apiGroupAuth.Use(mw.Auth)
	// 公开
	apiGroupPublic := apiGroup.Group("")

	// 调试、性能、系统信息
	{
		apiGroupPublic.GET("/sys/info/version", api.Sys.Version, mw.NewIPLimiter(3, 3))
		apiGroupPublic.GET("/sys/info/version/full", api.Sys.VersionFull, mw.NewIPLimiter(3, 3))
		apiGroupPublic.GET("/sys/info/branding", api.Sys.Branding, mw.NewIPLimiter(3, 3))
		apiGroupAdmin.GET("/sys/info/cache", api.Sys.InfoCache, mw.NewIPLimiter(3, 3))
	}

	// 类型数据
	{
		apiGroupPublic.GET("/type/role/type", api.Type.RoleType, mw.NewIPLimiter(3, 3))
		apiGroupPublic.GET("/type/period/type", api.Type.PeriodType, mw.NewIPLimiter(3, 3))
		apiGroupPublic.GET("/type/transaction/type", api.Type.TransactionType, mw.NewIPLimiter(3, 3))
		apiGroupPublic.GET("/type/wallet/type", api.Type.WalletType, mw.NewIPLimiter(3, 3))
		apiGroupPublic.GET("/type/wallet_partition/average", api.Type.WalletPartitionAverageType, mw.NewIPLimiter(3, 3))
	}

	// 用户相关
	{
		// 管理员
		{
			apiGroupAdmin.POST("/admin/user/create", api.Admin.CreateUser, mw.NewIPLimiter(0.5, 1))
			apiGroupAdmin.GET("/admin/user/find", api.Admin.Find, mw.NewIPLimiter(3, 3))
			apiGroupAdmin.POST("/admin/user/update", api.Admin.UpdateUser, mw.NewIPLimiter(1, 1))
			apiGroupAdmin.POST("/admin/user/disable", api.Admin.DisableUser, mw.NewIPLimiter(3, 3))
			apiGroupAdmin.POST("/admin/user/enable", api.Admin.EnableUser, mw.NewIPLimiter(3, 3))
			apiGroupAdmin.GET("/admin/user/access_token/all", api.Admin.AccessTokenInfo, mw.NewIPLimiter(3, 3))
			apiGroupAdmin.GET("/admin/user/login_log/all", api.Admin.LoginLogInfo, mw.NewIPLimiter(3, 3))
		}
		// 公开
		{
			apiGroupPublic.POST("/user/login", api.User.Login, mw.NewIPLimiter(3, 3))
		}
		// 登录用户
		{
			apiGroupAuth.POST("/user/logout", api.User.Logout, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/user/info", api.User.Info, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user/update", api.User.Update, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user/access_token/create", api.User.AccessTokenCreate, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/user/access_token/info", api.User.AccessTokenInfo, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/user/login_log/info", api.User.LoginLogInfo, mw.NewIPLimiter(3, 3))
		}
		// 用户绑定数据
		{
			apiGroupAuth.GET("/user/bind/home/tips/find", api.UserBind.HomeTipsFind, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user/bind/home/tips/save", api.UserBind.HomeTipsSave, mw.NewIPLimiter(3, 3))
		}
		// 银行
		{
			apiGroupAuth.POST("/bank/create", api.Bank.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.GET("/bank/find", api.Bank.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/bank/update", api.Bank.Update, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/bank/delete", api.Bank.Delete, mw.NewIPLimiter(3, 3))
		}
		// 货币类型
		{
			apiGroupAuth.POST("/currency/create", api.Currency.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.GET("/currency/find", api.Currency.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/currency/update", api.Currency.Update, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/currency/delete", api.Currency.Delete, mw.NewIPLimiter(3, 3))
		}
		// 银行卡类型
		{
			apiGroupAuth.POST("/card_type/create", api.CardType.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.GET("/card_type/find", api.CardType.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/card_type/update", api.CardType.Update, mw.NewIPLimiter(1, 1))
			apiGroupAuth.POST("/card_type/delete", api.CardType.Delete, mw.NewIPLimiter(3, 3))
		}
		// 用户银行卡
		{
			apiGroupAuth.POST("/user_card/create", api.UserCard.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.GET("/user_card/find", api.UserCard.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_card/valid/request", api.UserCard.ValidRequest, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_card/valid/input", api.UserCard.ValidInput, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_card/valid/cancel", api.UserCard.ValidCancel, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_card/update", api.UserCard.Update, mw.NewIPLimiter(1, 1))
			apiGroupAuth.POST("/user_card/update/sequence", api.UserCard.UpdateSequence, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_card/disable", api.UserCard.Disable, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_card/enable", api.UserCard.Enable, mw.NewIPLimiter(3, 3))
		}
		// 用户银行卡货币
		{
			apiGroupAuth.POST("/user_card_currency/create", api.UserCardCurrency.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.POST("/user_card_currency/update/balance", api.UserCardCurrency.UpdateBalance, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_card_currency/delete", api.UserCardCurrency.Delete, mw.NewIPLimiter(3, 3))
		}
		// 用户货币汇率
		{
			apiGroupAuth.POST("/user_exchange_rate/create", api.UserExchangeRate.Create, mw.NewIPLimiter(1, 1))
			apiGroupAuth.GET("/user_exchange_rate/find", api.UserExchangeRate.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_exchange_rate/update", api.UserExchangeRate.Update, mw.NewIPLimiter(1, 1))
			apiGroupAuth.POST("/user_exchange_rate/delete", api.UserExchangeRate.Delete, mw.NewIPLimiter(3, 3))
		}
		// 用户钱包
		{
			apiGroupAuth.POST("/wallet/create", api.UserWallet.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.GET("/wallet/find", api.UserWallet.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/wallet/find/normal", api.UserWallet.FindNormal, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/wallet/find/debt", api.UserWallet.FindDebt, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/wallet/find/wishlist", api.UserWallet.FindWishlist, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/wallet/update", api.UserWallet.Update, mw.NewIPLimiter(1, 1))
			apiGroupAuth.POST("/wallet/update/sequence", api.UserWallet.UpdateSequence, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/wallet/disable", api.UserWallet.Disable, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/wallet/enable", api.UserWallet.Enable, mw.NewIPLimiter(3, 3))
		}
		// 用户钱包划分
		{
			apiGroupAuth.POST("/wallet_partition/create", api.UserWalletPartition.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.POST("/wallet_partition/update", api.UserWalletPartition.Update, mw.NewIPLimiter(1, 1))
			apiGroupAuth.POST("/wallet_partition/update/sequence", api.UserWalletPartition.UpdateSequence, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/wallet_partition/disable", api.UserWalletPartition.Disable, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/wallet_partition/enable", api.UserWalletPartition.Enable, mw.NewIPLimiter(3, 3))
		}
		// 用户交易分类
		{
			apiGroupAuth.POST("/user_transaction_category/create", api.UserTransactionCategory.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.GET("/user_transaction_category/find", api.UserTransactionCategory.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_transaction_category/update", api.UserTransactionCategory.Update, mw.NewIPLimiter(1, 1))
			apiGroupAuth.POST("/user_transaction_category/update/sequence", api.UserTransactionCategory.UpdateSequence, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_transaction_category/delete", api.UserTransactionCategory.Delete, mw.NewIPLimiter(3, 3))
		}
		// 用户交易
		{
			apiGroupAuth.POST("/user_transaction/create", api.UserTransaction.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.GET("/user_transaction/find", api.UserTransaction.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_transaction/find", api.UserTransaction.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/user_transaction/find/range", api.UserTransaction.FindRange, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_transaction/update", api.UserTransaction.Update, mw.NewIPLimiter(1, 1))
			apiGroupAuth.POST("/user_transaction/checked", api.UserTransaction.Checked, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_transaction/delete", api.UserTransaction.Delete, mw.NewIPLimiter(3, 3))
		}
		// 交易View
		{
			apiGroupAuth.GET("/user_transaction/view/day", api.UserTransaction.ViewDay, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/user_transaction/view/month", api.UserTransaction.ViewMonth, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/user_transaction/view/year", api.UserTransaction.ViewYear, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/user_transaction/view/total", api.UserTransaction.ViewTotal, mw.NewIPLimiter(3, 3))
		}
		// 交易Chart
		{
			apiGroupAuth.POST("/user_transaction/chart", api.UserTransaction.Chart, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_transaction/chart/pie", api.UserTransaction.PieChart, mw.NewIPLimiter(3, 3))
		}
		// 用户周期性付费
		{
			apiGroupAuth.POST("/user_period_pay/create", api.UserPeriodPay.Create, mw.NewIPLimiter(0.5, 1))
			apiGroupAuth.GET("/user_period_pay/find", api.UserPeriodPay.Find, mw.NewIPLimiter(3, 3))
			apiGroupAuth.GET("/user_period_pay/summary", api.UserPeriodPay.Summary, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_period_pay/update", api.UserPeriodPay.Update, mw.NewIPLimiter(1, 1))
			apiGroupAuth.POST("/user_period_pay/update/next", api.UserPeriodPay.UpdateNextPeriod, mw.NewIPLimiter(3, 3))
			apiGroupAuth.POST("/user_period_pay/delete", api.UserPeriodPay.Delete, mw.NewIPLimiter(3, 3))
		}
	}
}
