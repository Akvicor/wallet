package migrate

import (
	"encoding/json"
	"github.com/Akvicor/glog"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

// 创建默认银行
func initBank() (err error) {
	array := []*model.Bank{
		model.NewBank("中国银行", "Bank of China", "BOC", "95566"),
		model.NewBank("中国建设银行", "China Construction Bank", "CCB", "95533"),
		model.NewBank("中国农业银行", "Agricultural Bank of China", "ABC", "95599"),
		model.NewBank("中国工商银行", "Industrial and Commercial Bank of China", "ICBC", "95588"),
		model.NewBank("中国邮政储蓄银行", "Postal Savings Bank of China", "PSBC", "95580"),
		model.NewBank("交通银行", "Bank of Communications", "BCM", "95559"),
		model.NewBank("招商银行", "China Merchants Bank", "CMB", "95555"),
		model.NewBank("青岛银行", "Bank of Qingdao", "BQD", "96588"),
		model.NewBank("中国光大银行", "China Everbright Bank", "CEB", "95595"),
		model.NewBank("华夏银行", "Huaxia Bank", "HXB", "95577"),
		model.NewBank("中国民生银行", "China Minsheng Bank", "CMBC", "95568"),
		model.NewBank("上海浦东发展银行", "Shanghai Pudong Development Bank", "SPDB", "95528"),
		model.NewBank("兴业银行", "Industrial Bank", "CIB", "95561"),
		model.NewBank("中信银行", "China CITIC Bank", "CNCB", "95558"),
		model.NewBank("花旗银行", "Citibank", "CITI", "800-830-1880"),
		model.NewBank("汇丰银行", "Hongkong and Shanghai Banking Corporation", "HSBC", "800-820-8878"),
		model.NewBank("恒生银行", "Hang Seng Bank", "HSB", "800-830-4888"),
		model.NewBank("美国银行", "Bank of America", "BOA", "6678-0000"),
		model.NewBank("华美银行", "East West Bank", "EWBC", "895-5650"),
		model.NewBank("杜高斯贝银行", "Dukascopy Bank", "DSB", "-"),
		model.NewBank("Wise", "Wise", "WISE", "-"),
		model.NewBank("WildCard", "WildCard", "WDC", "-"),
		model.NewBank("币安", "Binance", "Binance", "-"),
		model.NewBank("欧易", "OKEX", "OKX", "-"),
		model.NewBank("Bybit", "Bybit", "Bybit", "-"),
		model.NewBank("Coinbase", "Coinbase Exchange", "Coinbase", "-"),
		model.NewBank("Magic Eden", "Magic Eden", "ME", "-"),
		model.NewBank("比特币", "Bitcoin", "BTC", "-"),
		model.NewBank("以太坊", "Ethereum", "ETH", "-"),
		model.NewBank("Solana", "Solana", "SOL", "-"),
		model.NewBank("小红卡", "RedotPay", "Redot", "-"),
	}
	// 初始化银行
	for _, v := range array {
		exist, err := service.Bank.Exist(v.Name, v.EnglishName)
		if err != nil {
			glog.Warning("判断银行[%s][%s]存在性失败: %v", v.Name, v.EnglishName, err)
			return err
		}
		if exist {
			continue
		}
		_, err = service.Bank.Create(v.Name, v.EnglishName, v.Abbr, v.Phone)
		if err != nil {
			glog.Fatal("初始化银行数据异常: %s %v", v.Name, err)
			return err
		}
	}
	str, _ := json.Marshal(array)
	glog.Info("初始化银行成功: %s", string(str))
	return nil
}
