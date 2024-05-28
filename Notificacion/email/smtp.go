package email

import (
	"log"

	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func EnviarCorreo(config SMTPConfig, destinatario string, asunto string, cuerpo string) {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Username)
	m.SetHeader("To", destinatario)
	m.SetHeader("Subject", asunto)
	m.SetBody("text/plain", cuerpo)

	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error al enviar el correo: %v", err)
		return
	}

	log.Println("Correo enviado exitosamente")
}
