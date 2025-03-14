package function

import (
	"crypto/tls"
	"net/mail"

	"gopkg.in/gomail.v2"
)

func ServerTrace(subject, content string) error {
	from := mail.Address{Name: "vgga", Address: "vgga.dev@gmail.com"}
	to := mail.Address{Name: "vgga", Address: "vgga.dev@gmail.com"}
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from.String()) // 发送人
	mailer.SetHeader("To", to.String())     // 收件人
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", content)
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "vgga.dev@gmail.com", "radarxvtkhucvyrz")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}
	return nil
}
