package mail

import (
	"bytes"
	"html/template"
	"indetail/conf"
	"log"
	"net/smtp"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	auth := smtp.PlainAuth("", conf.EMAIL_ACCOUNT, conf.PASSWORD_ACCOUNT, conf.HOST_NAME)
	if err := smtp.SendMail(conf.HOST_NAME + conf.PORT, auth, conf.EMAIL_FROM, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *Request) Send(templateName string, items interface{}) bool {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Println(err)
	}
	if ok := r.sendMail(); ok {
		return true
	} else {
		return false
	}
}
