package app

import (
	"github.com/go-co-op/gocron/v2"
	"wallet/cmd/app/server/app/cron"
)

func setupCron(scheduler gocron.Scheduler) error {
	var err error

	// 每日记账通知
	err = cron.Sys.DailyReminderNotify(scheduler)
	if err != nil {
		return err
	}

	// 每月资金调度通知
	err = cron.Sys.ScheduleNotify(scheduler)
	if err != nil {
		return err
	}

	// 银行卡到期通知
	err = cron.UserCard.ExpDateNotify(scheduler)
	if err != nil {
		return err
	}

	// 银行卡账单日通知
	err = cron.UserCard.StatementClosingDayNotify(scheduler)
	if err != nil {
		return err
	}

	// 银行卡还款日通知
	err = cron.UserCard.PaymentDueDayNotify(scheduler)
	if err != nil {
		return err
	}

	// 周期订阅通知
	err = cron.UserPeriodPay.Notify(scheduler)
	if err != nil {
		return err
	}

	scheduler.Start()
	return nil
}
