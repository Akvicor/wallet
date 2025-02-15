package cron

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"time"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/global/mail"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

var _util utilModel

type utilModel struct{}

func (u *utilModel) TransactionDetails(incMail bool, userId int64, prefix string, start, end int64) (uncheckedR, incomeR, expenseR int, subjectR string, tableR *mail.HtmlTable) {
	income := 0                                         // 收入数量
	expense := 0                                        // 支出数量
	todayChecked := make([]*model.UserTransaction, 0)   // 今天的交易
	todayUnchecked := make([]*model.UserTransaction, 0) // 今天未确认的交易
	pastUnchecked := make([]*model.UserTransaction, 0)  // 过去的未确认的交易

	{
		today, _ := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), userId, start, end)
		todayUnique := make(map[int64]struct{})
		for _, t := range today {
			todayUnique[t.ID] = struct{}{}
			if t.Checked != 0 {
				todayChecked = append(todayChecked, t)
			} else {
				todayUnchecked = append(todayUnchecked, t)
			}
			switch t.Type {
			case transaction.TypeIncome:
				income++
			case transaction.TypeExpense:
				expense++
			}
		}
		unchecked, _ := service.UserTransaction.FindAllByUIDUnchecked(nil, true, model.NewPreloaderUserTransaction().All(), userId)
		for _, t := range unchecked {
			if _, ok := todayUnique[t.ID]; ok {
				continue
			}
			pastUnchecked = append(pastUnchecked, t)
		}
	}
	unchecked := len(todayUnchecked) + len(pastUnchecked)

	if income == 0 && expense == 0 {
		return unchecked, 0, 0, "", nil
	}

	subject := strings.Builder{}
	subject.WriteString(prefix)
	subject.WriteString("共")
	if income != 0 {
		subject.WriteString(fmt.Sprintf("%d个收入", income))
	}
	if expense != 0 {
		if income != 0 {
			subject.WriteString("和")
		}
		subject.WriteString(fmt.Sprintf("%d个支出", expense))
	}
	if unchecked > 0 {
		subject.WriteString(fmt.Sprintf("，有%d个未确认交易", unchecked))
	}

	if !incMail {
		return unchecked, income, expense, subject.String(), nil
	}
	content := mail.NewHtmlTable()
	{
		content.SetHeader([]*mail.HtmlTableHeader{
			{0, true, "类型"},
			{0, true, "分类"},
			{0, true, "来源"},
			{0, true, "目标"},
			{0, true, "描述"},
			{0, true, "交易时间"},
			{0, true, "确认时间"},
		})
		partitionF := func(f string, part *model.UserWalletPartition, currency *model.Currency, value decimal.Decimal) string {
			if part == nil || part.Wallet == nil || currency == nil {
				return "-"
			}
			return fmt.Sprintf("%s %s %s [%s|%s]", value, currency.Code, f, part.Wallet.Name, part.Name)
		}
		dateF := func(unix int64) string {
			if unix == 0 {
				return "-"
			}
			return time.Unix(unix, 0).Format("2006-01-02 15:04:05")
		}
		for _, tx := range todayChecked {
			content.AddRow([]*mail.HtmlTableRow{
				{false, tx.Type.Colour(), tx.Type.String()},
				{false, tx.Category.Colour, tx.Category.Name},
				{false, "", partitionF("<-", tx.FromPartition, tx.FromCurrency, tx.FromValue)},
				{false, "", partitionF("->", tx.ToPartition, tx.ToCurrency, tx.ToValue)},
				{false, "", tx.Description},
				{false, "", dateF(tx.Created)},
				{false, "", dateF(tx.Checked)},
			})
		}
		for _, tx := range todayUnchecked {
			content.AddRow([]*mail.HtmlTableRow{
				{false, tx.Type.Colour(), tx.Type.String()},
				{false, tx.Category.Colour, tx.Category.Name},
				{false, "", partitionF("<-", tx.FromPartition, tx.FromCurrency, tx.FromValue)},
				{false, "", partitionF("->", tx.ToPartition, tx.ToCurrency, tx.ToValue)},
				{false, "", tx.Description},
				{false, "", dateF(tx.Created)},
				{false, "", dateF(tx.Checked)},
			})
		}
		for _, tx := range pastUnchecked {
			content.AddRow([]*mail.HtmlTableRow{
				{false, tx.Type.Colour(), tx.Type.String()},
				{false, tx.Category.Colour, tx.Category.Name},
				{false, "", partitionF("<-", tx.FromPartition, tx.FromCurrency, tx.FromValue)},
				{false, "", partitionF("->", tx.ToPartition, tx.ToCurrency, tx.ToValue)},
				{false, "", tx.Description},
				{false, "", dateF(tx.Created)},
				{false, "", dateF(tx.Checked)},
			})
		}
	}
	return unchecked, income, expense, subject.String(), content
}
