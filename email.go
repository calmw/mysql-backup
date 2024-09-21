package backup

import (
	"github.com/go-gomail/gomail"
)

type Email struct {
	Email    string
	Username string
	Password string
	Host     string
	Port     int
}

func NewEmail(email, username, password, host string, port int) *Email {
	return &Email{
		Email:    email,
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (e *Email) Send(to, subject, body string, attach []string) error {
	// 创建邮件消息
	msg := gomail.NewMessage()
	msg.SetHeader("From", e.Email)    // 设置发件人
	msg.SetHeader("To", to)           // 设置收件人
	msg.SetHeader("Subject", subject) // 设置邮件主题
	msg.SetBody("text/html", body)    // 设置邮件正文

	// 添加附件
	for _, f := range attach {
		msg.Attach(f) // 添加文件附件
	}

	// 设置SMTP配置
	d := gomail.NewDialer(e.Host, e.Port, e.Email, e.Password)

	// 发送邮件
	return d.DialAndSend(msg)
}
