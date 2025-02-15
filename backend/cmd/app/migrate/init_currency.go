package migrate

import (
	"encoding/json"
	"github.com/Akvicor/glog"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

// 创建默认货币
func initCurrency() (err error) {
	array := []*model.Currency{
		model.NewCurrency("人民币", "Chinese Yuan", "CNY", "¥"),
		model.NewCurrency("美元", "US Dollar", "USD", "$"),
		model.NewCurrency("日元", "Japanese Yen", "JPY", "¥"),
		model.NewCurrency("港元", "Hong Kong Dollar", "HKD", "$"),
		model.NewCurrency("新台币", "New Taiwan Dollar", "TWD", "$"),
		model.NewCurrency("欧元", "Euro", "EUR", "€"),
		model.NewCurrency("英镑", "British Pound Sterling", "GBP", "£"),
		model.NewCurrency("加元", "Canadian Dollar", "CAD", "$"),
		model.NewCurrency("印度卢比", "Indian Rupee", "INR", "₹"),
		model.NewCurrency("澳大利亚元", "Australian Dollar", "AUD", "$"),
		model.NewCurrency("新西兰元", "New Zealand Dollar", "NZD", "$"),
		model.NewCurrency("新加坡元", "Singapore Dollar", "SGD", "$"),
		model.NewCurrency("瑞士法郎", "Swiss Franc", "CHF", "₣"),
		model.NewCurrency("黄金", "Gold", "AU", "g"),
		model.NewCurrency("白银", "Silver", "AG", "g"),
		model.NewCurrency("比特币", "Bitcoin", "BTC", "₿"),
		model.NewCurrency("以太币", "Ethereum", "ETH", "Ξ"),
		model.NewCurrency("泰达币", "Tether", "USDT", "₮"),
	}
	// 初始化货币
	for _, v := range array {
		exist, err := service.Currency.Exist(v.Name, v.EnglishName)
		if err != nil {
			glog.Warning("判断货币[%s][%s]存在性失败: %v", v.Name, v.EnglishName, err)
			return err
		}
		if exist {
			continue
		}
		_, err = service.Currency.Create(v.Name, v.EnglishName, v.Code, v.Symbol)
		if err != nil {
			glog.Fatal("初始化货币数据异常: %s %v", v.Name, err)
			return err
		}
	}
	str, _ := json.Marshal(array)
	glog.Info("初始化货币成功: %s", string(str))
	return nil
}
