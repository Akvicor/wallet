package auth

import (
	"github.com/labstack/echo/v4"
	"wallet/cmd/app/server/common/cache"
	"wallet/cmd/app/server/common/token"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

// GetCacheAuthorizationByToken 从Cache中通过Token获取Authorization
func GetCacheAuthorizationByToken(tokenStr string) *Authorization {
	get, ok := cache.TokenManager.Get(tokenStr)
	if ok {
		return get.(*Authorization)
	}
	return nil
}

// GetAuthorizationByToken 从Cache/DB中通过Token获取Authorization
func GetAuthorizationByToken(tokenStr string) *Authorization {
	if len(tokenStr) == 0 {
		return nil
	}
	authorization := GetCacheAuthorizationByToken(tokenStr)
	if authorization != nil {
		return authorization
	}
	tokenV := token.Parse(tokenStr)
	if tokenV == nil {
		return nil
	}
	if tokenV.Type == token.TypeAccessToken {
		accessToken, err := service.UserAccessToken.FindByToken(true, model.NewPreloaderAccessToken().All(), tokenStr)
		if err != nil || accessToken == nil {
			return nil
		}
		authorization = NewAccessAuthorization(accessToken.Token, accessToken.User)
	}
	return authorization
}

func GetUser(c echo.Context) *model.User {
	authorization := GetAuthorizationByToken(token.GetToken(c))
	if authorization == nil {
		return nil
	}
	return authorization.User
}

func GetAuthorization(c echo.Context) *Authorization {
	authorization := GetAuthorizationByToken(token.GetToken(c))
	return authorization
}
