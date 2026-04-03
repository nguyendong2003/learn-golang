package pkg

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"time"

	"elearning-api/config"
)

type EmailProvider interface {
	SendEmail(to string, subject string, body string) error
	SendPasswordReset(to string, fullName string, resetCode string) error
	SendEventReminder(to string, eventName string, studentName string, startTime time.Time, roomToken string) error
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
	subject := "Reset Your Password"
	resetLink := fmt.Sprintf(
		"%s/reset-password?email=%s&code=%s",
		s.FrontendBaseUrl,
		to,
		resetCode,
	)
	body := fmt.Sprintf(`
	<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f5f7fa; padding: 30px;">
	<div style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; border: 1px solid #e6e6e6;">

		<!-- Header -->
		<div style="background: #1a73e8; color: white; padding: 20px 30px;">
		<h2 style="margin: 0; font-size: 20px;">🔐 Password Reset</h2>
		</div>

		<!-- Body -->
		<div style="padding: 30px;">
		<p style="margin-top: 0;">Hi <strong>%s</strong>,</p>

		<p>
			We received a request to reset the password for your account.
			Click the button below to create a new password.
		</p>

		<!-- CTA -->
		<div style="text-align: center; margin: 30px 0;">
			<a href="%s"
			style="background-color: #1a73e8; color: #ffffff; padding: 14px 28px; border-radius: 6px; text-decoration: none; font-weight: 600; display: inline-block;">
			Reset Password
			</a>
		</div>

		<!-- Expiry -->
		<div style="background: #f1f3f4; border-radius: 6px; padding: 14px; font-size: 14px; color: #5f6368; margin-bottom: 20px;">
			⏳ This link will expire in <strong>15 minutes</strong> for security reasons.
		</div>

		<!-- Warning -->
		<p style="font-size: 14px; color: #d93025;">
			If you did NOT request this, please ignore this email. Your account is still secure.
		</p>

		<!-- Fallback link -->
		<p style="font-size: 13px; color: #5f6368; margin-top: 20px;">
			If the button doesn't work, copy and paste this link into your browser:
		</p>
		<p style="word-break: break-all; font-size: 13px; color: #1a73e8;">
			%s
		</p>
		</div>

		<!-- Footer -->
		<div style="background: #fafafa; padding: 20px; text-align: center; font-size: 12px; color: #9aa0a6;">
		<p style="margin: 0;">For security reasons, do not share this email or link with anyone.</p>
		<p style="margin: 5px 0 0;">© 2026 Your Company. All rights reserved.</p>
		</div>

	</div>
	</div>
	`,
		fullName,
		resetLink,
		resetLink,
	)
	return s.SendEmail(to, subject, body)
}

func (s *emailProvider) SendEventReminder(to string, eventName string, studentName string, startTime time.Time, roomToken string) error {
	subject := fmt.Sprintf("Reminder: %s starts soon!", eventName)
	body := fmt.Sprintf(`
		<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f5f7fa; padding: 30px;">
		<div style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; border: 1px solid #e6e6e6;">
			
			<!-- Header -->
			<div style="background: #1a73e8; color: white; padding: 20px 30px;">
			<h2 style="margin: 0; font-size: 20px;">📅 Event Reminder</h2>
			</div>

			<!-- Body -->
			<div style="padding: 30px;">
			<p style="margin-top: 0;">Hi <strong>%s</strong>,</p>

			<p>
				This is a reminder that your event 
				<strong style="color: #1a73e8;">"%s"</strong> will start soon.
			</p>

			<!-- Time Box -->
			<div style="background: #f1f3f4; border-radius: 6px; padding: 16px; margin: 20px 0; text-align: center;">
				<div style="font-size: 14px; color: #5f6368;">Start Time</div>
				<div style="font-size: 20px; font-weight: 600; color: #202124;">
				%s
				</div>
			</div>

			<p>Please make sure you are ready before the session begins.</p>

			<!-- CTA -->
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s"
				style="background-color: #34a853; color: white; padding: 14px 28px; border-radius: 6px; text-decoration: none; font-weight: 600; display: inline-block;">
				Join Event
				</a>
			</div>

			<!-- Fallback link -->
			<p style="font-size: 13px; color: #5f6368;">
				If the button doesn't work, copy and paste this link into your browser:
			</p>
			<p style="word-break: break-all; font-size: 13px; color: #1a73e8;">
				%s
			</p>
			</div>

			<!-- Footer -->
			<div style="background: #fafafa; padding: 20px; text-align: center; font-size: 12px; color: #9aa0a6;">
			<p style="margin: 0;">You're receiving this reminder because you're registered for this event.</p>
			<p style="margin: 5px 0 0;">© 2026 Your Company. All rights reserved.</p>
			</div>

		</div>
		</div>
		`,
		studentName,
		eventName,
		startTime.Local().Format("Mon, 02 Jan 2006 15:04 MST"),
		roomToken,
		roomToken,
	)

	return s.SendEmail(to, subject, body)
}
