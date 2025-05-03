package utils

import (
	"log"

	"github.com/bishalcode869/Auth-as-a-Service.git/configs"
	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(to string, subject string, body string) error {
	cfg := configs.Config
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.EmailAddress)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// setup smtp client
	d := gomail.NewDialer(cfg.SMPTServer, cfg.SMPTPort, cfg.EmailAddress, cfg.EmailPassword)
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email: ", err)
		return err
	}

	log.Println("Verfication email sent successfully to", to)
	return nil
}
