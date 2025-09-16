package services

import (
	"backend-golang/internal/helpers"
)

type EmailService interface {
	SendVerificationEmail(email, username, link string) error
	SendResetPasswordEmail(email, username, link string) error
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}

func (s *emailService) SendVerificationEmail(email, username, link string) error {
	return helpers.SendEmail(email, username, link, "verification_email", "Verifikasi Email Anda")
}

func (s *emailService) SendResetPasswordEmail(email, username, link string) error {
	return helpers.SendEmail(email, username, link, "reset_password_email", "Reset Password Anda")
}
