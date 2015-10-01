package core

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"
	"path"
	"strings"
)

func encodeRFC2047(String string) string {
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

type EmailTemplate struct {
	Files   []string
	To      mail.Address
	Subject string
	Data    interface{}
}

func NewEmailTemplate(p []string, n string, e string, s string, d interface{}) *EmailTemplate {
	return &EmailTemplate{
		Files:   p,
		To:      mail.Address{n, e},
		Subject: s,
		Data:    d,
	}
}

type EmailSender struct {
	*EmailConfig
	sender       mail.Address
	templatePath string
}

func (e *EmailSender) Parse(et *EmailTemplate) (body string, err error) {

	var doc bytes.Buffer

	ts := make([]string, 0)

	for _, v := range et.Files {
		ts = append(ts, path.Join(e.templatePath, v))
	}

	t := template.Must(template.ParseFiles(ts...))

	err = t.Execute(&doc, struct {
		From    string
		To      string
		Subject string
		Data    interface{}
	}{
		e.sender.String(),
		et.To.String(),
		et.Subject,
		et.Data,
	})

	if err != nil {
		return "", err
	}

	return doc.String(), nil
}

func (e *EmailSender) Send(et *EmailTemplate) (err error) {

	body, err := e.Parse(et)
	if err != nil {
		return err
	}

	server := fmt.Sprintf("%s:%d", e.Host, e.Port)

	auth := smtp.PlainAuth(
		"",
		e.Username,
		e.Password,
		e.Host,
	)

	go smtp.SendMail(server, auth, e.sender.Address, []string{et.To.Address}, []byte(body))

	return nil

}

func NewEmailSender(config *EmailConfig, p string) *EmailSender {
	var from mail.Address

	sender := strings.Split(config.Sender, ",")

	if len(sender) >= 2 {
		from = mail.Address{sender[0], sender[1]}
	} else {
		from = mail.Address{"", sender[1]}
	}

	return &EmailSender{config, from, p}
}
