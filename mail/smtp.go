package mail

import (
	"errors"
	"github.com/go-gomail/gomail"
	"time"
)

type SmtpSender struct {
	email   string
	pass    string
	host    string
	port    int
	timeout time.Duration
}

func NewSmtpSender(email, pass, host string, port int, timeout time.Duration) (*SmtpSender, error) {
	if !isEmailValid(email) {
		return nil, errors.New("invalid email email")
	}

	return &SmtpSender{email: email, pass: pass, host: host, port: port, timeout: timeout}, nil
}

func (s *SmtpSender) Send(payload Payload, attempts int) error {
	if err := payload.Validate(); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.email)
	msg.SetHeader("To", payload.To)
	msg.SetHeader("Subject", payload.Subject)
	msg.SetBody("text/html", payload.Body)

	dialer := gomail.NewDialer(s.host, s.port, s.email, s.pass)

	var err error = nil
	for i := 0; i < attempts; i++ {
		if err = dialer.DialAndSend(msg); err == nil {
			break
		} else {
			err = errors.New("failed to sent email via smtp")
			time.Sleep(s.timeout)
		}
	}

	return err
}
