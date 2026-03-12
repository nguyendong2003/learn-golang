package pkg

import (
	"fmt"
	"net/mail"
	"net/smtp"

	"elearning-api/config"
)

type EmailProvider interface {
	SendEmail(to string, subject string, body string) error
	SendPasswordReset(to string, fullName string, resetCode string) error
}

type emailProvider struct {
	emailConfig     *config.EmailConfig
	FrontendBaseUrl string
}

func NewEmailProvider(emailConfig *config.EmailConfig, frontendBaseUrl string) EmailProvider {
	return &emailProvider{
		emailConfig:     emailConfig,
		FrontendBaseUrl: frontendBaseUrl,
	}
}

func (s *emailProvider) SendEmail(to string, subject string, body string) error {

	if s.emailConfig == nil {
		return fmt.Errorf("email config is nil")
	}

	from := mail.Address{Address: s.emailConfig.SMTPUser}

	_, err := mail.ParseAddress(to)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth(
		"",
		s.emailConfig.SMTPUser,
		s.emailConfig.SMTPPass,
		s.emailConfig.SMTPHost,
	)

	smtpAddr := fmt.Sprintf("%s:%d", s.emailConfig.SMTPHost, s.emailConfig.SMTPPort)

	message := []byte(
		"From: " + from.Address + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body,
	)

	err = smtp.SendMail(
		smtpAddr,
		auth,
		s.emailConfig.SMTPUser,
		[]string{to},
		message,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (s *emailProvider) SendPasswordReset(to string, fullName string, resetCode string) error {

	subject := "Yêu cầu khôi phục mật khẩu"

	resetLink := fmt.Sprintf(
		"%s/reset-password?email=%s&code=%s",
		s.FrontendBaseUrl,
		to,
		resetCode,
	)

	body := fmt.Sprintf(`
	<div style="font-family: sans-serif; max-width: 600px; margin: auto; border: 1px solid #eee; padding: 20px;">
		<h2 style="color: #333;">Chào %s,</h2>
		<p>Chúng tôi nhận được yêu cầu khôi phục mật khẩu cho tài khoản của bạn.</p>
		<p>Vui lòng nhấn vào nút bên dưới để tiến hành đặt lại mật khẩu mới:</p>

		<div style="text-align: center; margin: 30px 0;">
			<a href="%s"
			style="background-color: #4A90E2; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; font-weight: bold;">
			Đặt lại mật khẩu ngay
			</a>
		</div>

		<p style="font-size: 13px; color: #666;">Nếu nút trên không hoạt động, hãy copy link sau:</p>
		<p style="font-size: 12px; color: #4A90E2;">%s</p>

		<hr>
		<p style="font-size: 12px; color: #999;">Link này sẽ hết hạn sau 15 phút.</p>
	</div>
	`,
		fullName,
		resetLink,
		resetLink,
	)

	return s.SendEmail(to, subject, body)
}