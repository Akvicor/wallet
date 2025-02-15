package cron

import (
	"errors"
	"fmt"
	"github.com/Akvicor/glog"
	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm"
	"time"
	"wallet/cmd/app/server/global/mail"
	"wallet/cmd/app/server/global/sms"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
	"wallet/cmd/config"
)

var UserPeriodPay = new(userPeriodPayCron)

type userPeriodPayCron struct{}

// Notify 周期订阅通知
func (u *userPeriodPayCron) Notify(scheduler gocron.Scheduler) (err error) {
	// 执行时间
	definition := gocron.DailyJob(
		config.Global.Cron.Offset.UserPeriodPayNotify.Interval,
		gocron.NewAtTimes(
			gocron.NewAtTime(config.Global.Cron.Offset.UserPeriodPayNotify.Hour, config.Global.Cron.Offset.UserPeriodPayNotify.Minute, config.Global.Cron.Offset.UserPeriodPayNotify.Second),
		),
	)
	// 任务
	task := gocron.NewTask(
		func() {
			glog.Debug("find all period pay in %d days", config.Global.Cron.Offset.UserPeriodPayNotify.DayLimit)
			pays, err := service.UserPeriodPay.FindAllByPeriodByDayLimit(true, model.NewPreloaderUserPeriodPay().All(), config.Global.Cron.Offset.UserPeriodPayNotify.DayLimit)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return
				}
				glog.Warning("find period pay, %v", err)
				return
			}

			userMap := make(map[int64]*model.User)
			userPaysMap := make(map[int64][]*model.UserPeriodPay)
			for _, pay := range pays {
				user, ok := userMap[pay.UID]
				if ok {
					pay.User = user
					userPaysMap[pay.UID] = append(userPaysMap[pay.UID], pay)
				} else {
					userMap[pay.UID] = pay.User
					user = pay.User
					userPaysMap[pay.UID] = make([]*model.UserPeriodPay, 0, 1)
					userPaysMap[pay.UID] = append(userPaysMap[pay.UID], pay)
				}
			}
			for _, user := range userMap {
				userPays := userPaysMap[user.ID]
				subject := fmt.Sprintf("您在未来%d天内有%d个循环付款", config.Global.Cron.Offset.UserPeriodPayNotify.DayLimit, len(userPays))
				if user.Mail != "" {
					content := mail.NewHtmlTable()
					content.SetHeader([]*mail.HtmlTableHeader{
						{0, true, "Name"},
						{0, true, "Description"},
						{0, true, "Value"},
						{0, true, "Date"},
					})
					for _, pay := range userPays {
						content.AddRow([]*mail.HtmlTableRow{
							{false, "", pay.Name},
							{false, "", pay.Description},
							{false, "", fmt.Sprintf("%s %s", pay.Value.String(), pay.Currency.Code)},
							{false, "", time.Unix(pay.NextOfPeriod, 0).Format("2006-01-02 15:04:05")},
						})
					}
					_ = mail.Send(user.Mail, subject, mail.Html, content.Content())
				}
				if user.Phone != "" {
					sms.Send(user.Phone, subject)
				}
			}
		},
	)
	// 创建cron
	_, err = scheduler.NewJob(definition, task)
	return err
}
