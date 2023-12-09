package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	Dialer *gomail.Dialer
	Sender string
}

func NewMailer(host string, port int, username, password, sender string) *Mailer {
	dialer := gomail.NewDialer(host, port, username, password)
	return &Mailer{
		Dialer: dialer,
		Sender: sender,
	}
}

func (m Mailer) Send(recpient string, username string) error {
	subject := "Appointment Reminder"
	plainBody := fmt.Sprintf("Hello %s, Do not forget to complete your medical visit\n", username)

	msg := gomail.NewMessage()
	msg.SetHeader("To", recpient)
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plainBody)

	err := m.Dialer.DialAndSend(msg)
	return err
}
