package sms

import (
	"fmt"
	"github.com/Akvicor/glog"
	"github.com/Akvicor/util"
	"wallet/cmd/config"
)

func Send(phone, msg string) {
	if !config.Global.SMS.Enable {
		glog.Warning("未启用短信发送功能，跳过发送短信。")
		return
	}

	res, _ := util.HttpPost(config.Global.SMS.URL, map[string]string{
		"key":     config.Global.SMS.Token,
		"sender":  config.Global.AppName,
		"phone":   phone,
		"message": msg,
	}, util.HTTPContentTypeUrlencoded, nil)
	jsonRes := util.NewJSONResult(res)
	if jsonRes == nil {
		glog.Error("短信发送失败, HTTP请求错误")
		return
	}

	result := jsonRes.Map()
	if fmt.Sprint(result["code"]) != "0" {
		glog.Error("短信发送失败, %v", result["msg"])
	}
}
