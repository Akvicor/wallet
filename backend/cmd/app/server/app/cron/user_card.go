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

var UserCard = new(userCardCron)

type userCardCron struct{}

// ExpDateNotify 银行卡过期通知
func (u *userCardCron) ExpDateNotify(scheduler gocron.Scheduler) (err error) {
	// 执行时间
	definition := gocron.DailyJob(
		config.Global.Cron.Offset.UserCardExpDateNotify.Interval,
		gocron.NewAtTimes(
			gocron.NewAtTime(config.Global.Cron.Offset.UserCardExpDateNotify.Hour, config.Global.Cron.Offset.UserCardExpDateNotify.Minute, config.Global.Cron.Offset.UserCardExpDateNotify.Second),
		),
	)
	// 几天内过期的银行卡
	afterDay := config.Global.Cron.Offset.UserCardExpDateNotify.DayLimit
	// 任务
	task := gocron.NewTask(
		func(after int) {
			glog.Debug("find all expiring soon user card")
			cards, err := service.UserCard.FindExpiringSoon(true, model.NewPreloaderUserCard().User().Bank(), time.Now().AddDate(0, 0, after).Unix())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return
				}
				glog.Warning("find expiring soon failed, %v", err)
				return
			}
			userMap := make(map[int64]*model.User)
			userCardsMap := make(map[int64][]*model.UserCard)
			for _, card := range cards {
				user, ok := userMap[card.UID]
				if ok {
					card.User = user
					userCardsMap[card.UID] = append(userCardsMap[card.UID], card)
				} else {
					userMap[card.UID] = card.User
					user = card.User
					userCardsMap[card.UID] = make([]*model.UserCard, 0, 1)
					userCardsMap[card.UID] = append(userCardsMap[card.UID], card)
				}
			}
			for _, user := range userMap {
				if user.Mail == "" {
					continue
				}
				userCards := userCardsMap[user.ID]
				subject := fmt.Sprintf("您有%d张银行卡在%d天内过期", len(userCards), after)
				if user.Mail != "" {
					content := mail.NewHtmlTable()
					content.SetHeader([]*mail.HtmlTableHeader{
						{0, true, "Name"},
						{0, true, "Bank"},
						{0, true, "Number"},
						{0, true, "Expiring"},
					})
					for _, card := range userCards {
						content.AddRow([]*mail.HtmlTableRow{
							{false, "", card.Name},
							{false, "", card.Bank.Name},
							{false, "", card.Number},
							{false, "", time.Unix(card.ExpDate, 0).Format("2006-01-02 15:04:05 -0700")},
						})
					}
					_ = mail.Send(user.Mail, subject, mail.Html, content.Content())
				}
				if user.Phone != "" {
					sms.Send(user.Phone, subject)
				}
			}
		},
		afterDay,
	)
	// 创建cron
	_, err = scheduler.NewJob(definition, task)
	return err
}

// StatementClosingDayNotify 银行卡账单日通知
func (u *userCardCron) StatementClosingDayNotify(scheduler gocron.Scheduler) (err error) {
	// 执行时间
	definition := gocron.DailyJob(
		config.Global.Cron.Offset.UserCardStatementClosingDayNotify.Interval,
		gocron.NewAtTimes(
			gocron.NewAtTime(config.Global.Cron.Offset.UserCardStatementClosingDayNotify.Hour, config.Global.Cron.Offset.UserCardStatementClosingDayNotify.Minute, config.Global.Cron.Offset.UserCardStatementClosingDayNotify.Second),
		),
	)
	// 任务
	task := gocron.NewTask(
		func() {
			glog.Debug("find all statement closing day user card")
			cards, err := service.UserCard.FindAllByStatementClosingDay(true, model.NewPreloaderUserCard().User().Bank(), int64(time.Now().AddDate(0, 0, -config.Global.Cron.Offset.UserCardStatementClosingDayNotify.DelayDay).Day()))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return
				}
				glog.Warning("find statement closing day user card failed, %v", err)
				return
			}
			userMap := make(map[int64]*model.User)
			userCardsMap := make(map[int64][]*model.UserCard)
			for _, card := range cards {
				user, ok := userMap[card.UID]
				if ok {
					card.User = user
					userCardsMap[card.UID] = append(userCardsMap[card.UID], card)
				} else {
					userMap[card.UID] = card.User
					user = card.User
					userCardsMap[card.UID] = make([]*model.UserCard, 0, 1)
					userCardsMap[card.UID] = append(userCardsMap[card.UID], card)
				}
			}
			for _, user := range userMap {
				if user.Mail == "" {
					continue
				}
				userCards := userCardsMap[user.ID]
				subject := fmt.Sprintf("您有%d张信用卡处于账单日", len(userCards))
				if user.Mail != "" {
					content := mail.NewHtmlTable()
					content.SetHeader([]*mail.HtmlTableHeader{
						{0, true, "Name"},
						{0, true, "Bank"},
						{0, true, "Number"},
					})
					for _, card := range userCards {
						content.AddRow([]*mail.HtmlTableRow{
							{false, "", card.Name},
							{false, "", card.Bank.Name},
							{false, "", card.Number},
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

// PaymentDueDayNotify 银行卡还款日通知
func (u *userCardCron) PaymentDueDayNotify(scheduler gocron.Scheduler) (err error) {
	// 执行时间
	definition := gocron.DailyJob(
		config.Global.Cron.Offset.UserCardPaymentDueDayNotify.Interval,
		gocron.NewAtTimes(
			gocron.NewAtTime(config.Global.Cron.Offset.UserCardPaymentDueDayNotify.Hour, config.Global.Cron.Offset.UserCardPaymentDueDayNotify.Minute, config.Global.Cron.Offset.UserCardPaymentDueDayNotify.Second),
		),
	)
	// 任务
	task := gocron.NewTask(
		func() {
			glog.Debug("find all payment due day user card")
			cards, err := service.UserCard.FindAllByPaymentDueDay(true, model.NewPreloaderUserCard().User().Bank(), int64(time.Now().AddDate(0, 0, -config.Global.Cron.Offset.UserCardPaymentDueDayNotify.DelayDay).Day()))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return
				}
				glog.Warning("find payment due day user card failed, %v", err)
				return
			}
			userMap := make(map[int64]*model.User)
			userCardsMap := make(map[int64][]*model.UserCard)
			for _, card := range cards {
				user, ok := userMap[card.UID]
				if ok {
					card.User = user
					userCardsMap[card.UID] = append(userCardsMap[card.UID], card)
				} else {
					userMap[card.UID] = card.User
					user = card.User
					userCardsMap[card.UID] = make([]*model.UserCard, 0, 1)
					userCardsMap[card.UID] = append(userCardsMap[card.UID], card)
				}
			}
			for _, user := range userMap {
				if user.Mail == "" {
					continue
				}
				userCards := userCardsMap[user.ID]
				subject := fmt.Sprintf("您有%d张信用卡处于还款日", len(userCards))
				if user.Mail != "" {
					content := mail.NewHtmlTable()
					content.SetHeader([]*mail.HtmlTableHeader{
						{0, true, "Name"},
						{0, true, "Bank"},
						{0, true, "Number"},
					})
					for _, card := range userCards {
						content.AddRow([]*mail.HtmlTableRow{
							{false, "", card.Name},
							{false, "", card.Bank.Name},
							{false, "", card.Number},
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
