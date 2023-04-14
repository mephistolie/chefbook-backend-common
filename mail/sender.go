package mail

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

type Sender interface {
	Send(input Payload, attempts int) error
}

type Payload struct {
	To      string
	Subject string
	Body    string
}

func (p *Payload) SetHtmlBody(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		fmt.Printf("failed to parse file %s:%s", templateFileName, err.Error())
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	p.Body = buf.String()

	return nil
}

func (p *Payload) Validate() error {
	if p.To == "" {
		return errors.New("empty to")
	}

	if p.Subject == "" || p.Body == "" {
		return errors.New("empty subject/body")
	}

	if !isEmailValid(p.To) {
		return errors.New("invalid receiver email")
	}

	return nil
}
