package service

import (
	"bank-api/internal/config"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-mail/mail/v2"
)

type EmailService struct {
	dialer *mail.Dialer
	from   string
}

func NewEmailService(cfg *config.Config) *EmailService {
	dialer := mail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
	dialer.TLSConfig = &tls.Config{
		ServerName:         cfg.SMTPHost,
		InsecureSkipVerify: false,
	}

	return &EmailService{
		dialer: dialer,
		from:   cfg.SMTPUser,
	}
}

// SendPaymentNotification отправляет уведомление о платеже
func (s *EmailService) SendPaymentNotification(to string, amount float64, description string) error {
	subject := "Уведомление о платеже"
	body := fmt.Sprintf(`
        <h1>Платеж успешно выполнен</h1>
        <p>Сумма: <strong>%.2f RUB</strong></p>
        <p>Описание: %s</p>
        <small>Это автоматическое уведомление от банка</small>
    `, amount, description)

	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent to %s", to)
	return nil
}
