package mailing

import (
	"fmt"
	"go-api/internal/helpers/env"
	"gopkg.in/gomail.v2"
	"strconv"
)

func SendEmail(toEmail string, htmlBody string, subject string) error {
	from := env.GetVar("MAILING_USERNAME")
	password := env.GetVar("MAILING_PASSWORD")
	smtpHost := env.GetVar("MAILING_SMTP_SERVER")
	smtpPortStr := env.GetVar("MAILING_SMTP_PORT")

	portInt, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(smtpHost, portInt, from, password)
	d.SSL = portInt == 465

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
