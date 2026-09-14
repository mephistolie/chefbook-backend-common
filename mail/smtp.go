package mail

import (
	"errors"
	"github.com/go-gomail/gomail"
	"time"
)

type SmtpSender struct {
	email    string
	username string
	pass     string
	host     string
	port     int
	timeout  time.Duration
}

func NewSmtpSender(email, pass, host string, port int, timeout time.Duration, username ...string) (*SmtpSender, error) {
	if !isEmailValid(email) {
		return nil, errors.New("invalid email email")
	}

	login := email
	if len(username) > 0 && username[0] != "" {
		login = username[0]
	}
	return &SmtpSender{email: email, username: login, pass: pass, host: host, port: port, timeout: timeout}, nil
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

	dialer := gomail.NewDialer(s.host, s.port, s.username, s.pass)

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
