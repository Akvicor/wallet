package api

import (
	"errors"
	"github.com/jinzhu/now"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"time"
	"wallet/cmd/app/server/app/dro"
	"wallet/cmd/app/server/app/dto"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/common/types/transaction"
	"wallet/cmd/app/server/global/auth"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

var UserTransaction = new(userTransactionApi)

type userTransactionApi struct{}

// Create 创建交易
func (a *userTransactionApi) Create(c echo.Context) (err error) {
	input := new(dto.UserTransactionCreate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	switch input.TransactionType {
	case transaction.TypeIncome:
		_, err = service.UserTransaction.CreateIncome(transaction.TypeIncome, self.ID, input.ToPartitionID, input.CurrencyID, input.CategoryID, input.Description, input.ToValue, input.Created)
	case transaction.TypeExpense:
		_, err = service.UserTransaction.CreateExpense(transaction.TypeExpense, self.ID, input.FromPartitionID, input.CurrencyID, input.CategoryID, input.Description, input.ToValue, input.Created)
	case transaction.TypeTransfer:
		_, err = service.UserTransaction.CreateTransfer(transaction.TypeTransfer, self.ID, input.FromPartitionID, input.ToPartitionID, input.CategoryID, input.Description, input.FromValue, input.ToValue, input.Created)
	case transaction.TypeAutoTransfer:
		_, err = service.UserTransaction.CreateAutoTransfer(transaction.TypeAutoTransfer, self.ID, input.FromPartitionID, input.ToPartitionID, input.CurrencyID, input.CategoryID, input.Description, input.FromValue, input.Created)
	case transaction.TypeExchange:
		_, err = service.UserTransaction.CreateExchange(transaction.TypeExchange, self.ID, input.FromPartitionID, input.ToPartitionID, input.CategoryID, input.Description, input.FromValue, input.ToValue, input.Created)
	case transaction.TypeAutoExchange:
		_, err = service.UserTransaction.CreateAutoExchange(transaction.TypeAutoExchange, self.ID, input.FromPartitionID, input.ToPartitionID, input.CategoryID, input.Description, input.ToValue, input.Created)
	default:
		return resp.FailWithMsg(c, resp.BadRequest, "错误交易类型")
	}
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "创建失败: "+err.Error())
	}
	return resp.Success(c)
}

// Find 获取全部交易
func (a *userTransactionApi) Find(c echo.Context) (err error) {
	input := new(dto.UserTransactionFind)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	if input.ID == 0 {
		var transactions = make([]*model.UserTransaction, 0)
		transactions, err = service.UserTransaction.FindAllByUID(&input.PageModel, !input.All, model.NewPreloaderUserTransaction().All(), self.ID, input.Search, input.FromPartitionID, input.ToPartitionID, input.FromCurrencyID, input.ToCurrencyID, input.CategoryID, input.TransactionType)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		return resp.SuccessWithPageData(c, &input.PageModel, transactions)
	} else {
		var card *model.UserTransaction
		card, err = service.UserTransaction.FindByUID(!input.All, model.NewPreloaderUserTransaction().All(), self.ID, input.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		return resp.SuccessWithData(c, card)
	}
}

// FindRange 获取时间范围全部交易
func (a *userTransactionApi) FindRange(c echo.Context) (err error) {
	input := new(dto.UserTransactionFindRange)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	var transactions = make([]*model.UserTransaction, 0)
	transactions, err = service.UserTransaction.FindAllByUIDCreated(nil, !input.All, model.NewPreloaderUserTransaction().All(), self.ID, input.Start, input.End)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	return resp.SuccessWithData(c, transactions)
}

// Update 更新交易
func (a *userTransactionApi) Update(c echo.Context) (err error) {
	input := new(dto.UserTransactionUpdate)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	err = service.UserTransaction.UpdateByUID(self.ID, input.ID, input.CategoryID, input.Description, input.Created)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "更新失败: "+err.Error())
	}
	return resp.Success(c)
}

// Checked 确认交易
func (a *userTransactionApi) Checked(c echo.Context) (err error) {
	input := new(dto.UserTransactionChecked)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	err = service.UserTransaction.UpdateCheckedByUID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "确认失败: "+err.Error())
	}
	return resp.Success(c)
}

// Delete 删除交易
func (a *userTransactionApi) Delete(c echo.Context) (err error) {
	input := new(dto.UserTransactionDelete)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	err = service.UserTransaction.DeleteByUID(self.ID, input.ID)
	if err != nil {
		return resp.FailWithMsg(c, resp.Failed, "删除失败: "+err.Error())
	}
	return resp.Success(c)
}

// ViewDay 获取天全部交易
func (a *userTransactionApi) ViewDay(c echo.Context) (err error) {
	input := new(dto.UserTransactionViewDay)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	target := now.With(now.With(time.Unix(input.Target, 0)).BeginningOfDay())
	var start, end int64
	result := new(dro.UserTransactionViewDayResp)
	result.Merge = new(dro.UserTransactionViewResp)
	result.Days = make(map[string]*dro.UserTransactionViewResp)
	if input.Direction > 0 {
		for i := 0; i < input.Direction; i++ {
			targetI := now.With(target.AddDate(0, 0, i))
			start = targetI.BeginningOfDay().Unix()
			end = targetI.EndOfDay().Unix()
			list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return resp.FailWithMsg(c, resp.Failed, "系统错误")
			}
			item := new(dro.UserTransactionViewResp)
			item.ParseTransactions(list)
			result.Days[targetI.Format(time.DateOnly)] = item
			result.Merge.ParseTransactions(list)
		}
		return resp.SuccessWithData(c, result)
	} else if input.Direction < 0 {
		for i := 0; i > input.Direction; i-- {
			targetI := now.With(target.AddDate(0, 0, i))
			start = targetI.BeginningOfDay().Unix()
			end = targetI.EndOfDay().Unix()
			list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return resp.FailWithMsg(c, resp.Failed, "系统错误")
			}
			item := new(dro.UserTransactionViewResp)
			item.ParseTransactions(list)
			result.Days[targetI.Format(time.DateOnly)] = item
			result.Merge.ParseTransactions(list)
		}
		return resp.SuccessWithData(c, result)
	} else {
		start = target.BeginningOfDay().Unix()
		end = target.EndOfDay().Unix()
		list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.FailWithMsg(c, resp.Failed, "系统错误")
		}
		item := new(dro.UserTransactionViewResp)
		item.ParseTransactions(list)
		return resp.SuccessWithData(c, item)
	}
}

// ViewMonth 获取月全部交易
func (a *userTransactionApi) ViewMonth(c echo.Context) (err error) {
	input := new(dto.UserTransactionViewMonth)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	target := now.With(now.With(time.Unix(input.Target, 0)).BeginningOfMonth())
	var start, end int64
	result := new(dro.UserTransactionViewMonthResp)
	result.Merge = new(dro.UserTransactionViewResp)
	result.Months = make(map[string]*dro.UserTransactionViewDayResp)

	if input.Direction > 0 {
		for i := 0; i < input.Direction; i++ {
			targetI := now.With(target.AddDate(0, i, 0))
			resultMonth := new(dro.UserTransactionViewDayResp)
			resultMonth.Merge = new(dro.UserTransactionViewResp)
			resultMonth.Days = make(map[string]*dro.UserTransactionViewResp)
			result.Months[targetI.Format("2006-01")] = resultMonth

			startDay := targetI.BeginningOfMonth().Day()
			endDay := targetI.EndOfMonth().Day()
			targetStart := targetI.BeginningOfMonth()
			for ; startDay <= endDay; startDay++ {
				targetDay := now.With(targetStart.AddDate(0, 0, startDay-1))
				start = targetDay.BeginningOfDay().Unix()
				end = targetDay.EndOfDay().Unix()
				list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return resp.FailWithMsg(c, resp.Failed, "系统错误")
				}
				item := new(dro.UserTransactionViewResp)
				item.ParseTransactions(list)
				resultMonth.Days[targetDay.Format(time.DateOnly)] = item
				resultMonth.Merge.ParseTransactions(list)
				result.Merge.ParseTransactions(list)
			}
		}
		return resp.SuccessWithData(c, result)
	} else if input.Direction < 0 {
		for i := 0; i > input.Direction; i-- {
			targetI := now.With(target.AddDate(0, i, 0))
			resultMonth := new(dro.UserTransactionViewDayResp)
			resultMonth.Merge = new(dro.UserTransactionViewResp)
			resultMonth.Days = make(map[string]*dro.UserTransactionViewResp)
			result.Months[targetI.Format("2006-01")] = resultMonth

			startDay := targetI.BeginningOfMonth().Day()
			endDay := targetI.EndOfMonth().Day()
			targetStart := targetI.BeginningOfMonth()
			for ; startDay <= endDay; startDay++ {
				targetDay := now.With(targetStart.AddDate(0, 0, startDay-1))
				start = targetDay.BeginningOfDay().Unix()
				end = targetDay.EndOfDay().Unix()
				list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return resp.FailWithMsg(c, resp.Failed, "系统错误")
				}
				item := new(dro.UserTransactionViewResp)
				item.ParseTransactions(list)
				resultMonth.Days[targetDay.Format(time.DateOnly)] = item
				resultMonth.Merge.ParseTransactions(list)
				result.Merge.ParseTransactions(list)
			}
		}
		return resp.SuccessWithData(c, result)
	} else {
		resultMonth := new(dro.UserTransactionViewDayResp)
		resultMonth.Merge = new(dro.UserTransactionViewResp)
		resultMonth.Days = make(map[string]*dro.UserTransactionViewResp)

		startDay := target.BeginningOfMonth().Day()
		endDay := target.EndOfMonth().Day()
		targetStart := target.BeginningOfMonth()
		for ; startDay <= endDay; startDay++ {
			targetDay := now.With(targetStart.AddDate(0, 0, startDay-1))
			start = targetDay.BeginningOfDay().Unix()
			end = targetDay.EndOfDay().Unix()
			list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return resp.FailWithMsg(c, resp.Failed, "系统错误")
			}
			item := new(dro.UserTransactionViewResp)
			item.ParseTransactions(list)
			resultMonth.Days[targetDay.Format(time.DateOnly)] = item
			resultMonth.Merge.ParseTransactions(list)
		}
		return resp.SuccessWithData(c, resultMonth)
	}
}

// ViewYear 获取年全部交易
func (a *userTransactionApi) ViewYear(c echo.Context) (err error) {
	input := new(dto.UserTransactionViewYear)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	target := now.With(now.With(time.Unix(input.Target, 0)).BeginningOfYear())
	var start, end int64
	result := new(dro.UserTransactionViewYearResp)
	result.Merge = new(dro.UserTransactionViewResp)
	result.Years = make(map[string]*dro.UserTransactionViewMonthResp)

	if input.Direction > 0 {
		for year := 0; year < input.Direction; year++ {
			targetYear := now.With(target.AddDate(year, 0, 0))
			resultYear := new(dro.UserTransactionViewMonthResp)
			resultYear.Merge = new(dro.UserTransactionViewResp)
			resultYear.Months = make(map[string]*dro.UserTransactionViewDayResp)
			result.Years[targetYear.Format("2006")] = resultYear

			startMonth := int(targetYear.BeginningOfYear().Month())
			endMonth := int(targetYear.EndOfYear().Month())
			for ; startMonth <= endMonth; startMonth++ {
				targetDay := now.With(targetYear.AddDate(0, startMonth-1, 0))
				resultMonth := new(dro.UserTransactionViewDayResp)
				resultMonth.Merge = new(dro.UserTransactionViewResp)
				resultMonth.Days = make(map[string]*dro.UserTransactionViewResp)
				resultYear.Months[targetDay.Format("2006-01")] = resultMonth

				startDay := targetDay.BeginningOfMonth().Day()
				endDay := targetDay.EndOfMonth().Day()
				targetStart := targetDay.BeginningOfMonth()
				for ; startDay <= endDay; startDay++ {
					targetStartDay := now.With(targetStart.AddDate(0, 0, startDay-1))
					start = targetStartDay.BeginningOfDay().Unix()
					end = targetStartDay.EndOfDay().Unix()
					list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return resp.FailWithMsg(c, resp.Failed, "系统错误")
					}
					item := new(dro.UserTransactionViewResp)
					item.ParseTransactions(list)
					resultMonth.Days[targetStartDay.Format(time.DateOnly)] = item
					resultMonth.Merge.ParseTransactions(list)
					resultYear.Merge.ParseTransactions(list)
					result.Merge.ParseTransactions(list)
				}
			}
		}
		return resp.SuccessWithData(c, result)
	} else if input.Direction < 0 {
		for year := 0; year > input.Direction; year-- {
			targetYear := now.With(target.AddDate(year, 0, 0))
			resultYear := new(dro.UserTransactionViewMonthResp)
			resultYear.Merge = new(dro.UserTransactionViewResp)
			resultYear.Months = make(map[string]*dro.UserTransactionViewDayResp)
			result.Years[targetYear.Format("2006")] = resultYear

			startMonth := int(targetYear.BeginningOfYear().Month())
			endMonth := int(targetYear.EndOfYear().Month())
			for ; startMonth <= endMonth; startMonth++ {
				targetDay := now.With(targetYear.AddDate(0, startMonth-1, 0))
				resultMonth := new(dro.UserTransactionViewDayResp)
				resultMonth.Merge = new(dro.UserTransactionViewResp)
				resultMonth.Days = make(map[string]*dro.UserTransactionViewResp)
				resultYear.Months[targetDay.Format("2006-01")] = resultMonth

				startDay := targetDay.BeginningOfMonth().Day()
				endDay := targetDay.EndOfMonth().Day()
				targetStart := targetDay.BeginningOfMonth()
				for ; startDay <= endDay; startDay++ {
					targetStartDay := now.With(targetStart.AddDate(0, 0, startDay-1))
					start = targetStartDay.BeginningOfDay().Unix()
					end = targetStartDay.EndOfDay().Unix()
					list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return resp.FailWithMsg(c, resp.Failed, "系统错误")
					}
					item := new(dro.UserTransactionViewResp)
					item.ParseTransactions(list)
					resultMonth.Days[targetStartDay.Format(time.DateOnly)] = item
					resultMonth.Merge.ParseTransactions(list)
					resultYear.Merge.ParseTransactions(list)
					result.Merge.ParseTransactions(list)
				}
			}
		}
		return resp.SuccessWithData(c, result)
	} else {
		resultYear := new(dro.UserTransactionViewMonthResp)
		resultYear.Merge = new(dro.UserTransactionViewResp)
		resultYear.Months = make(map[string]*dro.UserTransactionViewDayResp)

		startMonth := int(target.BeginningOfYear().Month())
		endMonth := int(target.EndOfYear().Month())
		for ; startMonth <= endMonth; startMonth++ {
			targetDay := now.With(target.AddDate(0, startMonth-1, 0))
			resultMonth := new(dro.UserTransactionViewDayResp)
			resultMonth.Merge = new(dro.UserTransactionViewResp)
			resultMonth.Days = make(map[string]*dro.UserTransactionViewResp)
			resultYear.Months[targetDay.Format("2006-01")] = resultMonth

			startDay := targetDay.BeginningOfMonth().Day()
			endDay := targetDay.EndOfMonth().Day()
			targetStart := targetDay.BeginningOfMonth()
			for ; startDay <= endDay; startDay++ {
				targetStartDay := now.With(targetStart.AddDate(0, 0, startDay-1))
				start = targetStartDay.BeginningOfDay().Unix()
				end = targetStartDay.EndOfDay().Unix()
				list, err := service.UserTransaction.FindAllByUIDCreated(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, start, end)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return resp.FailWithMsg(c, resp.Failed, "系统错误")
				}
				item := new(dro.UserTransactionViewResp)
				item.ParseTransactions(list)
				resultMonth.Days[targetStartDay.Format(time.DateOnly)] = item
				resultMonth.Merge.ParseTransactions(list)
				resultYear.Merge.ParseTransactions(list)
			}
		}
		return resp.SuccessWithData(c, resultYear)
	}
}

// ViewTotal 获取全部交易
func (a *userTransactionApi) ViewTotal(c echo.Context) (err error) {
	input := new(dto.UserTransactionViewTotal)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}

	list, err := service.UserTransaction.FindAllByUID(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, "", nil, nil, nil, nil, nil, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	result := new(dro.UserTransactionViewResp)
	result.ParseTransactions(list)
	return resp.SuccessWithData(c, result)
}

// Chart 获取时间段交易总额
func (a *userTransactionApi) Chart(c echo.Context) (err error) {
	input := new(dto.UserTransactionChart)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	if self.MasterCurrency == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请配置默认货币")
	}

	list, err := service.UserTransaction.FindAllByUIDCreatedFromPartition(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, input.FromPartitionID, input.Start, input.End)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	result := new(dro.UserTransactionChartResp)
	result.ParseTransactions(self.ID, self.MasterCurrency, input.Unit, input.Start, input.End, list)
	return resp.SuccessWithData(c, result)
}

// PieChart 获取时间段交易总额
func (a *userTransactionApi) PieChart(c echo.Context) (err error) {
	input := new(dto.UserTransactionPieChart)
	if err = c.Bind(input); err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, "错误输入: "+err.Error())
	}
	err = input.Validate()
	if err != nil {
		return resp.FailWithMsg(c, resp.BadRequest, err.Error())
	}
	self := auth.GetUser(c)
	if self == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请登录")
	}
	if self.MasterCurrency == nil {
		return resp.FailWithMsg(c, resp.UnAuthorized, "请配置默认货币")
	}

	list, err := service.UserTransaction.FindAllByUIDCreatedFromPartition(nil, true, model.NewPreloaderUserTransaction().All(), self.ID, input.FromPartitionID, input.Start, input.End)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp.FailWithMsg(c, resp.Failed, "系统错误")
	}
	result := new(dro.UserTransactionPieChartResp)
	result.ParseTransactions(self.MasterCurrency, list)
	return resp.SuccessWithData(c, result)
}
