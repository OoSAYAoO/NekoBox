package models

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"

	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"
)

func SendNewQuestionMail(pageID uint, question *Question) {
	user, err := GetUserByPage(pageID)
	if err != nil {
		return
	}
	page, err := GetPageByID(pageID)
	if err != nil {
		return
	}

	var mailContent bytes.Buffer
	t, _ := template.ParseFiles("views/mail/new_question.tpl")
	p := map[string]string{
		"link":     fmt.Sprintf("https://box.n3ko.co/_/%s/%d", page.Domain, question.ID),
		"question": question.Content,
	}
	_ = t.Execute(&mailContent, p)

	err = sendMail(user.Email, "【NekoBox】您有一个新的提问", mailContent.String())
	if err != nil {
		log.Println(err)
	}
}

func sendMail(to string, title string, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", beego.AppConfig.String("mail_account"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	mailPort, _ := beego.AppConfig.Int("mail_port")
	d := gomail.NewDialer(
		beego.AppConfig.String("smtp"),
		mailPort,
		beego.AppConfig.String("mail_account"),
		beego.AppConfig.String("mail_password"),
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}
