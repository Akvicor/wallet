package cron

import (
	"errors"
	"fmt"
	"github.com/Akvicor/glog"
	"github.com/go-co-op/gocron/v2"
	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"strings"
	"wallet/cmd/app/server/global/mail"
	"wallet/cmd/app/server/global/sms"
	"wallet/cmd/app/server/service"
	"wallet/cmd/config"
)

var Sys = new(sysCron)

type sysCron struct{}

// DailyReminderNotify 每日记账通知
func (s *sysCron) DailyReminderNotify(scheduler gocron.Scheduler) (err error) {
	// 执行时间
	definition := gocron.DailyJob(
		config.Global.Cron.Offset.DailyReminderNotify.Interval,
		gocron.NewAtTimes(
			gocron.NewAtTime(config.Global.Cron.Offset.DailyReminderNotify.Hour, config.Global.Cron.Offset.DailyReminderNotify.Minute, config.Global.Cron.Offset.DailyReminderNotify.Second),
		),
	)
	// 任务
	task := gocron.NewTask(
		func() {
			glog.Debug("daily reminder notify")
			users, err := service.User.FindAll(nil, true, nil)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return
				}
				glog.Warning("find users failed, %v", err)
				return
			}

			start := now.BeginningOfDay().Unix()
			end := now.EndOfDay().Unix()

			for _, user := range users {
				if user.Phone == "" && user.Mail == "" {
					continue
				}
				incPhone := user.Phone != ""
				incMail := user.Mail != ""

				unchecked, income, expense, subject, table := _util.TransactionDetails(incMail, user.ID, "今日", start, end)

				if income == 0 && expense == 0 {
					if user.Phone != "" {
						content := strings.Builder{}
						content.WriteString("[Wallet]\n今日无收入和支出")
						if unchecked > 0 {
							content.WriteString(fmt.Sprintf("，有%d个未确认交易", unchecked))
						}
						sms.Send(user.Phone, content.String())
					}
					continue
				}

				if incPhone {
					content := fmt.Sprintf("[Wallet]\n%s", subject)
					sms.Send(user.Phone, content)
				}
				if incMail {
					_ = mail.Send(user.Mail, subject, table.ContentType(), table.Content())
				}
			}
		},
	)
	// 创建cron
	_, err = scheduler.NewJob(definition, task)
	return err
}

// ScheduleNotify 月末资金调度
func (s *sysCron) ScheduleNotify(scheduler gocron.Scheduler) (err error) {
	// 执行时间
	definition := gocron.MonthlyJob(config.Global.Cron.Offset.ScheduleNotify.Interval,
		gocron.NewDaysOfTheMonth(config.Global.Cron.Offset.ScheduleNotify.DaysOfTheMonth),
		gocron.NewAtTimes(
			gocron.NewAtTime(config.Global.Cron.Offset.ScheduleNotify.Hour, config.Global.Cron.Offset.ScheduleNotify.Minute, config.Global.Cron.Offset.ScheduleNotify.Second),
		),
	)
	// 任务
	task := gocron.NewTask(
		func() {
			glog.Info("schedule notify")
			users, err := service.User.FindAll(nil, true, nil)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return
				}
				glog.Warning("find users failed, %v", err)
				return
			}

			start := now.BeginningOfMonth().Unix()
			end := now.EndOfMonth().Unix()
			for _, user := range users {
				if user.Phone == "" && user.Mail == "" {
					continue
				}
				incPhone := user.Phone != ""
				incMail := user.Mail != ""

				unchecked, income, expense, subject, table := _util.TransactionDetails(incMail, user.ID, "本月", start, end)

				if income == 0 && expense == 0 {
					if user.Phone != "" {
						content := strings.Builder{}
						content.WriteString("[Wallet]\n今天是当月最后一天, 记得清账")
						if unchecked > 0 {
							content.WriteString(fmt.Sprintf("，有%d个未确认交易", unchecked))
						}
						sms.Send(user.Phone, content.String())
					}
					continue
				}

				if incPhone {
					content := fmt.Sprintf("[Wallet]\n%s", subject)
					sms.Send(user.Phone, content)
				}
				if incMail {
					_ = mail.Send(user.Mail, subject, table.ContentType(), table.Content())
				}
			}
		},
	)
	// 创建cron
	_, err = scheduler.NewJob(definition, task)
	return err
}
