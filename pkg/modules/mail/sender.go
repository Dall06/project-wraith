package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type Mail interface {
	Send(tmpl string, content interface{}, subject string, to []string) error
}

type mail struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewMail(from, password, smtpHost, smtpPort string) Mail {
	return &mail{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (mc *mail) Send(tmpl string, content interface{}, subject string, to []string) error {
	parsedTmpl, err := template.ParseFiles(tmpl)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	headers := "MIME-Version: 1.0\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subject, headers)))

	err = parsedTmpl.Execute(&body, content)
	if err != nil {
		return err
	}

	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", mc.smtpHost, mc.smtpPort),
		smtp.PlainAuth("", mc.from, mc.password, mc.smtpHost),
		mc.from,
		to,
		body.Bytes())
	if err != nil {
		return err
	}

	return nil

}
