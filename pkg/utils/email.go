package utils

import (
	"crypto/tls"
	"log"
	"myapp/pkg/config"
	"strconv"

	"gopkg.in/gomail.v2"
)

var SendEmail = func(to, subject, body string) error {
	from := config.NewConfig().SMTPUser
	if from == "" {
		from = "weather-alert@localhost"
	}
	log.Printf("→ Sending email to %s | from=%s | subject=%q", to, from, subject)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	cfg := config.NewConfig()
	port, _ := strconv.Atoi(cfg.SMTPPort)
	d := gomail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUser, cfg.SMTPPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("✗ Error sending email to %s: %v", to, err)
		return err
	}
	log.Printf("✓ Email successfully sent to %s", to)
	return nil
}
