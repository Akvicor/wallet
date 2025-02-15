package mail

import (
	"errors"
	"github.com/Akvicor/glog"
	"github.com/wneessen/go-mail"
	"wallet/cmd/config"
)

const (
	Plain = mail.TypeTextPlain
	Html  = mail.TypeTextHTML
)

func Send(to, subject string, contentType mail.ContentType, body string) error {
	if !config.Global.Mail.Enable {
		glog.Warning("未启用邮件发送功能，跳过发送邮件。")
		return errors.New("未启用邮件发送功能")
	}

	var err error

	msg := mail.NewMsg()
	err = msg.From(config.Global.Mail.From)
	if err != nil {
		glog.Warning("邮件发送失败: %v", err)
		return errors.New("邮件发送失败")
	}
	err = msg.To(to)
	if err != nil {
		glog.Warning("邮件发送失败: %v", err)
		return errors.New("邮件发送失败")
	}
	msg.Subject(subject)
	msg.SetBodyString(contentType, body)

	client, err := mail.NewClient(config.Global.Mail.SmtpHost, mail.WithPort(config.Global.Mail.SmtpPort), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(config.Global.Mail.Username), mail.WithPassword(config.Global.Mail.Password), mail.WithSSL())
	if err != nil {
		glog.Warning("邮件发送失败: %v", err)
		return errors.New("邮件发送失败")
	}

	err = client.DialAndSend(msg)
	if err != nil {
		glog.Error("邮件发送失败, %v", err)
		return errors.New("邮件发送失败")
	}
	return nil
}
