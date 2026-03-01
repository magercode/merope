package services

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"merope/models"
	"merope/utils"
)

type EmailService struct {
	enabled    bool
	host       string
	port       string
	username   string
	password   string
	from       string
	to         []string
	lang       *utils.LanguageManager
}

func NewEmailService(lang *utils.LanguageManager) *EmailService {
	enabled := os.Getenv("EMAIL_ENABLED") == "true"
	
	return &EmailService{
		enabled:    enabled,
		host:       os.Getenv("SMTP_HOST"),
		port:       os.Getenv("SMTP_PORT"),
		username:   os.Getenv("SMTP_USERNAME"),
		password:   os.Getenv("SMTP_PASSWORD"),
		from:       os.Getenv("SMTP_FROM"),
		to:         strings.Split(os.Getenv("SMTP_TO"), ","),
		lang:       lang,
	}
}

func (e *EmailService) Send(alert *models.Alert) error {
	if !e.enabled {
		return nil
	}

	subject := fmt.Sprintf("[%s] %s", alert.Level, alert.Title)
	body := e.formatEmailBody(alert)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", strings.Join(e.to, ","), e.from, subject, body))

	addr := fmt.Sprintf("%s:%s", e.host, e.port)
	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	err := smtp.SendMail(addr, auth, e.from, e.to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (e *EmailService) formatEmailBody(alert *models.Alert) string {
	levelColor := "#28a745" 
	switch alert.Level {
	case models.WARNING:
		levelColor = "#ffc107" 
	case models.CRITICAL:
		levelColor = "#dc3545" 
	}

	recommendationHTML := ""
	if alert.Recommendation != "" {
		recommendationHTML = fmt.Sprintf(`<p><strong>%s</p>`, alert.Recommendation)
	}

	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background-color: #f8f9fa; padding: 20px; text-align: center; }
			.content { padding: 20px; }
			.alert { border-left: 4px solid %s; padding: 15px; background-color: #f8f9fa; }
			.footer { text-align: center; padding: 20px; color: #6c757d; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h2>🔔 %s</h2>
			</div>
			<div class="content">
				<div class="alert">
					<h3>%s</h3>
					<p><strong>%s:</strong> %s</p>
					%s
					<p><strong>%s:</strong> <span style="color: %s">%s</span></p>
					<p><strong>%s:</strong> %s</p>
				</div>
			</div>
			<div class="footer">
				<p>Merope Monitoring System</p>
			</div>
		</div>
	</body>
	</html>
	`, levelColor, e.lang.GetMessage("alert_title"), alert.Title,
		e.lang.GetMessage("alert_message"), alert.Message, 
		recommendationHTML,
		e.lang.GetMessage("level"), levelColor, alert.Level,
		e.lang.GetMessage("time"), alert.Time)
}

func (e *EmailService) IsEnabled() bool {
	return e.enabled
}