package api

import (
	"github.com/labstack/echo/v4"
	"wallet/cmd/app/server/app/dro"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/global/sys"
	"wallet/cmd/config"
	"wallet/cmd/def"
)

var Sys = new(sysApi)

type sysApi struct{}

func (a *sysApi) Version(c echo.Context) (err error) {
	return resp.SuccessWithData(c, def.Version)
}

func (a *sysApi) VersionFull(c echo.Context) (err error) {
	return resp.SuccessWithData(c, def.AppVersion())
}

func (a *sysApi) Branding(c echo.Context) error {
	return resp.SuccessWithData(c, resp.Map{
		"name":      config.Global.AppName,
		"copyright": def.Copyright(),
	})
}

func (a *sysApi) InfoCache(c echo.Context) (err error) {
	tokenManager := sys.TokenManagerItems()
	loginFailedManager := sys.LoginFailedManagerItems()

	return resp.SuccessWithData(c, resp.Map{
		"token_manager_count":        len(tokenManager),
		"token_manager":              tokenManager,
		"login_failed_manager_count": len(loginFailedManager),
		"login_failed_manager":       loginFailedManager,
	})
}

func (a *sysApi) Health(c echo.Context) error {
	health := &dro.SysHealth{
		Status: "healthy",
	}
	if health.Status == "healthy" {
		return resp.Healthy(c, health)
	} else {
		return resp.Unhealthy(c, health)
	}
}
