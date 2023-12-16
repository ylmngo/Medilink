package mail

import (
	"fmt"
	"net/smtp"
)

// type Mailer struct {
// 	Dialer *gomail.Dialer
// 	Sender string
// }

// func NewMailer(host string, port int, username, password, sender string) *Mailer {
// 	dialer := gomail.NewDialer(host, port, username, password)
// 	return &Mailer{
// 		Dialer: dialer,
// 		Sender: sender,
// 	}
// }

// func (m Mailer) Send(recpient string, username string) error {
// 	subject := "Appointment Reminder"
// 	plainBody := fmt.Sprintf("Hello %s, Do not forget to complete your medical visit\n", username)

// 	msg := gomail.NewMessage()
// 	msg.SetHeader("To", recpient)
// 	msg.SetHeader("From", m.Sender)
// 	msg.SetHeader("Subject", subject)
// 	msg.SetBody("text/plain", plainBody)

// 	err := m.Dialer.DialAndSend(msg)
// 	return err
// }

type Mailer struct {
	sender string
	pass   string
	host   string
	auth   smtp.Auth
}

func InitMailer(sender, pass, host string) *Mailer {
	return &Mailer{
		sender: sender,
		pass:   pass,
		host:   host,
		auth:   smtp.PlainAuth("", sender, pass, host),
	}
}

func (m *Mailer) Send(to string) error {
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: Appointment Reminder\n\n", m.sender, to)
	var err error
	for i := 0; i < 3; i++ {
		err = smtp.SendMail("smtp.gmail.com:587", m.auth, m.sender, []string{to}, []byte(message))
		if err == nil {
			break
		}
	}
	return err
}
