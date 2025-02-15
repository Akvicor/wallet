package dro

import (
	"github.com/shopspring/decimal"
	"slices"
	"sort"
	"time"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

/*
View Day
*/

func (view *UserTransactionViewResp) ParseTransactions(transactions []*model.UserTransaction) {
	if view.List == nil || len(view.List) == 0 {
		view.IncomeMerge = make([]*UserTransactionViewRespCategorySum, 0)
		view.ExpenseMerge = make([]*UserTransactionViewRespCategorySum, 0)
		view.Income = make([]*UserTransactionViewRespCategory, 0)
		view.Expense = make([]*UserTransactionViewRespCategory, 0)
		view.List = transactions
	} else {
		view.List = append(view.List, transactions...)
	}

	for _, tx := range transactions {
		if tx.Type == transaction.TypeIncome {
			foundMergeSum := false
			for _, income := range view.IncomeMerge {
				if income.Currency.ID == tx.ToCurrency.ID {
					foundMergeSum = true
					income.Value = income.Value.Add(tx.ToValue)
					break
				}
			}
			if !foundMergeSum {
				view.IncomeMerge = append(view.IncomeMerge, &UserTransactionViewRespCategorySum{
					Currency: tx.ToCurrency,
					Value:    tx.ToValue,
				})
			}
			foundCategory := false
			for _, income := range view.Income {
				if income.Category.ID == tx.Category.ID {
					foundCategory = true
					foundSum := false
					for _, sum := range income.Sum {
						if sum.Currency.ID == tx.ToCurrency.ID {
							foundSum = true
							sum.Value = sum.Value.Add(tx.ToValue)
							break
						}
					}
					if !foundSum {
						income.Sum = append(income.Sum, &UserTransactionViewRespCategorySum{
							Currency: tx.ToCurrency,
							Value:    tx.ToValue,
						})
					}
					break
				}
			}
			if !foundCategory {
				income := &UserTransactionViewRespCategory{
					Category: tx.Category,
					Sum:      make([]*UserTransactionViewRespCategorySum, 0),
				}
				income.Sum = append(income.Sum, &UserTransactionViewRespCategorySum{
					Currency: tx.ToCurrency,
					Value:    tx.ToValue,
				})
				view.Income = append(view.Income, income)
			}
		}
		if tx.Type == transaction.TypeExpense {
			foundMergeSum := false
			for _, expense := range view.ExpenseMerge {
				if expense.Currency.ID == tx.FromCurrency.ID {
					foundMergeSum = true
					expense.Value = expense.Value.Add(tx.FromValue)
					break
				}
			}
			if !foundMergeSum {
				view.ExpenseMerge = append(view.ExpenseMerge, &UserTransactionViewRespCategorySum{
					Currency: tx.FromCurrency,
					Value:    tx.FromValue,
				})
			}
			foundCategory := false
			for _, expense := range view.Expense {
				if expense.Category.ID == tx.Category.ID {
					foundCategory = true
					foundSum := false
					for _, sum := range expense.Sum {
						if sum.Currency.ID == tx.FromCurrency.ID {
							foundSum = true
							sum.Value = sum.Value.Add(tx.FromValue)
							break
						}
					}
					if !foundSum {
						expense.Sum = append(expense.Sum, &UserTransactionViewRespCategorySum{
							Currency: tx.FromCurrency,
							Value:    tx.FromValue,
						})
					}
					break
				}
			}
			if !foundCategory {
				expense := &UserTransactionViewRespCategory{
					Category: tx.Category,
					Sum:      make([]*UserTransactionViewRespCategorySum, 0),
				}
				expense.Sum = append(expense.Sum, &UserTransactionViewRespCategorySum{
					Currency: tx.FromCurrency,
					Value:    tx.FromValue,
				})
				view.Expense = append(view.Income, expense)
			}
		}
	}
}

func (chart *UserTransactionChartResp) ParseTransactions(uid int64, masterCurrency *model.Currency, unit string, start, end int64, transactions []*model.UserTransaction) {
	slices.SortFunc(transactions, func(a, b *model.UserTransaction) int {
		return int(a.Created - b.Created)
	})
	*chart = make(UserTransactionChartResp, 0)

	// 根据时间戳获取日期字符串
	getDate := func(created int64) string {
		switch unit {
		case "day":
			return time.Unix(created, 0).Format("2006-01-02")
		case "month":
			return time.Unix(created, 0).Format("2006-01")
		case "year":
			return time.Unix(created, 0).Format("2006")
		}
		return ""
	}
	// 区分货币
	getValue := func(tx *model.UserTransaction) decimal.Decimal {
		if tx.Type == transaction.TypeIncome {
			return tx.ToValue
		} else if tx.Type == transaction.TypeExpense {
			return tx.FromValue
		}
		return decimal.NewFromInt(0)
	}
	getCurrency := func(tx *model.UserTransaction) string {
		if tx.Type == transaction.TypeIncome {
			return tx.ToCurrency.Name
		} else if tx.Type == transaction.TypeExpense {
			return tx.FromCurrency.Name
		}
		return "Unknown"
	}
	// 日期相差
	dateDiff := int64(0)
	{
		startT := time.Unix(start, 0)
		endT := time.Unix(end, 0)
		switch unit {
		case "day":
			dateDiff = int64(endT.Sub(startT)/time.Second)/60/60/24 + 1
		case "month":
			dateDiff = int64(12)*int64(endT.Year()-startT.Year()) + int64(endT.Month()) - int64(startT.Month()) + 1
		case "year":
			dateDiff = int64(12)*int64(endT.Year()-startT.Year()) + 1
		}
		if dateDiff == 0 {
			dateDiff = 1
		}
	}
	// 计算汇率后的值
	//const CnyID = 1
	rateCache := make(map[int64]*model.UserExchangeRate)
	getSumValue := func(tx *model.UserTransaction) decimal.Decimal {
		var value decimal.Decimal
		var currency int64
		if tx.Type == transaction.TypeIncome {
			value = tx.ToValue
			currency = tx.ToCurrencyID
		} else if tx.Type == transaction.TypeExpense {
			value = tx.FromValue
			currency = tx.FromCurrencyID
		}
		if currency == masterCurrency.ID {
			return value
		}
		var rate *model.UserExchangeRate
		var ok bool
		var err error
		rate, ok = rateCache[currency]
		if !ok {
			rate, err = service.UserExchangeRate.FindByUID(false, nil, uid, currency, masterCurrency.ID)
			if err != nil {
				return value
			}
			rateCache[currency] = rate
		}
		return value.Mul(rate.Rate)
	}
	// 用于快速查找并合并拥有相同key的数据
	currencyMap := make(map[string]struct{}) // 记录出现过的货币
	totalExpense := decimal.NewFromInt(0)
	recordMap := make(map[string]map[string]*UserTransactionChartRespItem) // Date => Currency => Value
	for _, tx := range transactions {
		if tx.Type != transaction.TypeExpense {
			continue
		}
		iDate := getDate(tx.Created)
		var iCurrency string
		var iValue decimal.Decimal
		{
			// 不同类型货币单独计算
			iCurrency = getCurrency(tx)
			iValue = getValue(tx)
			_, ok := recordMap[iDate]
			if !ok {
				recordMap[iDate] = make(map[string]*UserTransactionChartRespItem)
			}
			item, ok := recordMap[iDate][iCurrency]
			if !ok {
				item = new(UserTransactionChartRespItem)
				item.Date = iDate
				item.Currency = iCurrency
				item.TValue = iValue
				recordMap[iDate][iCurrency] = item
			} else {
				item.TValue = item.TValue.Add(iValue)
			}
			currencyMap[iCurrency] = struct{}{}
		}
		{
			// 将所有货币按照人民币计算
			iCurrency = "O-SUM"
			iValue = getSumValue(tx)
			_, ok := recordMap[iDate]
			if !ok {
				recordMap[iDate] = make(map[string]*UserTransactionChartRespItem)
			}
			item, ok := recordMap[iDate][iCurrency]
			if !ok {
				item = new(UserTransactionChartRespItem)
				item.Date = iDate
				item.Currency = iCurrency
				item.TValue = iValue
				recordMap[iDate][iCurrency] = item
			} else {
				item.TValue = item.TValue.Add(iValue)
			}
			currencyMap[iCurrency] = struct{}{}
			totalExpense = totalExpense.Add(iValue)
		}
	}

	currencyList := make([]string, 0, len(currencyMap))
	for currency := range currencyMap {
		currencyList = append(currencyList, currency)
	}
	sort.Strings(currencyList)
	avgExpense := totalExpense.Div(decimal.NewFromInt(dateDiff))

	startT := time.Unix(start, 0)
	endT := time.Unix(end, 0)
	for startT.Before(endT) {
		iDate := getDate(startT.Unix())
		// 补足缺失的日期
		record, ok := recordMap[iDate]
		if !ok {
			record = make(map[string]*UserTransactionChartRespItem)
		}
		// 补足日期里缺失的货币
		for _, currency := range currencyList {
			_, ok := record[currency]
			if !ok {
				item := new(UserTransactionChartRespItem)
				item.Date = iDate
				item.Currency = currency
				item.TValue = decimal.NewFromInt(0)
				record[currency] = item
			}
		}

		// avg
		avg := new(UserTransactionChartRespItem)
		avg.Date = iDate
		avg.Currency = "AVG"
		avg.TValue = avgExpense
		*chart = append(*chart, avg)

		item := record["O-SUM"]
		if item != nil {
			*chart = append(*chart, item)
		}

		// currency
		for _, currency := range currencyList {
			if currency == "O-SUM" {
				continue
			}
			item = record[currency]
			if item == nil {
				continue
			}
			*chart = append(*chart, item)
		}

		addY, addM, addD := 0, 0, 0
		switch unit {
		case "day":
			dateDiff = int64(endT.Sub(startT)/time.Second)/60/60/24 + 1
			addY = 0
			addM = 0
			addD = 1
		case "month":
			dateDiff = int64(12)*int64(endT.Year()-startT.Year()) + int64(endT.Month()) - int64(startT.Month()) + 1
			addY = 0
			addM = 1
			addD = 0
		case "year":
			dateDiff = int64(12)*int64(endT.Year()-startT.Year()) + 1
			addY = 1
			addM = 0
			addD = 0
		default:
			addY = 0
			addM = 0
			addD = 1
		}
		startT = startT.AddDate(addY, addM, addD)
	}
	for _, item := range *chart {
		item.Value = item.TValue.Truncate(2).InexactFloat64()
	}
}

func (chart *UserTransactionPieChartResp) ParseTransactions(masterCurrency *model.Currency, transactions []*model.UserTransaction) {
	*chart = make(UserTransactionPieChartResp, 0)
	union := make(map[int64]*UserTransactionPieChartRespItem)

	rateCache := make(map[int64]*model.UserExchangeRate)
	getRateValue := func(tx *model.UserTransaction) decimal.Decimal {
		if tx.FromCurrencyID == masterCurrency.ID {
			return tx.FromValue
		}
		var rate *model.UserExchangeRate
		var ok bool
		var err error
		rate, ok = rateCache[tx.FromCurrencyID]
		if !ok {
			rate, err = service.UserExchangeRate.FindByUID(false, nil, tx.UID, tx.FromCurrencyID, masterCurrency.ID)
			if err != nil {
				return tx.FromValue
			}
			rateCache[tx.FromCurrencyID] = rate
		}
		return tx.FromValue.Mul(rate.Rate)
	}

	for _, tx := range transactions {
		if tx.Type != transaction.TypeExpense {
			continue
		}
		item, ok := union[tx.CategoryID]
		if !ok {
			item = &UserTransactionPieChartRespItem{
				Type:   tx.Category.Name,
				Colour: tx.Category.Colour,
				TValue: getRateValue(tx),
			}
			union[tx.CategoryID] = item
			*chart = append(*chart, item)
			continue
		}
		item.TValue = item.TValue.Add(getRateValue(tx))
	}
	for _, item := range *chart {
		item.Value = item.TValue.Truncate(2).InexactFloat64()
	}
}
