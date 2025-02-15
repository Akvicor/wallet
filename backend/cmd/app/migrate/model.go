package migrate

import (
	"wallet/cmd/app/server/model"
)

var list = []any{
	&model.User{}, &model.UserAccessToken{}, &model.LoginLog{},
	&model.Currency{}, &model.CardType{}, &model.Bank{},
	&model.UserCard{}, &model.UserCardCurrency{}, &model.UserExchangeRate{},
	&model.UserWallet{}, &model.UserWalletPartition{},
	&model.UserTransaction{}, &model.UserTransactionCategory{},
	&model.UserPeriodPay{}, &model.UserBindHomeTips{},
}
