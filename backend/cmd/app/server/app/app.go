package app

import (
	"fmt"
	"github.com/Akvicor/glog"
	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/config"
)

var app *App

type App struct {
	Cron   gocron.Scheduler
	Server *echo.Echo
}

func newApp() *App {
	return &App{}
}

func init() {
	app = newApp()
}

func Run() error {
	var err error

	if err = auth.LoadFromDBToCache(); err != nil {
		glog.Fatal("cannot load token from login log. %v", err)
	}

	app.Server = echo.New()
	app.Server.HideBanner = true
	setupRoutes(app.Server)

	app.Cron, err = gocron.NewScheduler()
	if err != nil {
		glog.Fatal("cannot create cron. %v", err)
	}
	err = setupCron(app.Cron)
	if err != nil {
		glog.Fatal("cannot setup cron. %v", err)
	}

	if config.Global.Server.EnableHttps && config.Global.Server.CrtFile != "" && config.Global.Server.KeyFile != "" {
		return app.Server.StartTLS(fmt.Sprintf("%s:%d", config.Global.Server.HttpIp, config.Global.Server.HttpPort), config.Global.Server.CrtFile, config.Global.Server.KeyFile)
	} else {
		return app.Server.Start(fmt.Sprintf("%s:%d", config.Global.Server.HttpIp, config.Global.Server.HttpPort))
	}
}
