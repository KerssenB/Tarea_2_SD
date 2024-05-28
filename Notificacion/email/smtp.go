package email

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTPConfig(host string, port int, username, password string) *SMTPConfig {
	return &SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (config *SMTPConfig) SendEmail(recipient, subject, body string) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	msg := []byte("To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	return smtp.SendMail(addr, auth, config.Username, []string{recipient}, msg)
}
