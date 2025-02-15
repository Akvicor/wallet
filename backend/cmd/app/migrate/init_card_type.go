package migrate

import (
	"encoding/json"
	"github.com/Akvicor/glog"
	"wallet/cmd/app/server/model"
	"wallet/cmd/app/server/service"
)

// 创建默认银行卡类型
func initCardType() (err error) {
	array := []*model.CardType{
		model.NewCardType("Union Pay"),
		model.NewCardType("Visa"),
		model.NewCardType("MasterCard"),
		model.NewCardType("JCB"),
		model.NewCardType("American Express"),
		model.NewCardType("Discover"),
		model.NewCardType("Diners Club"),
		model.NewCardType("Crypto"),
	}
	// 初始化银行卡类型
	for _, v := range array {
		exist, err := service.CardType.Exist(v.Name)
		if err != nil {
			glog.Warning("判断银行卡类型[%s]存在性失败: %v", v.Name, err)
			return err
		}
		if exist {
			continue
		}
		_, err = service.CardType.Create(v.Name)
		if err != nil {
			glog.Fatal("初始化银行卡类型数据异常: %s %v", v.Name, err)
			return err
		}
	}
	str, _ := json.Marshal(array)
	glog.Info("初始化银行卡类型成功: %s", string(str))
	return nil
}
