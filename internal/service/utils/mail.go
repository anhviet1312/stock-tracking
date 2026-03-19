package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/emersion/go-sasl"

	storyW "codebase/internal/models"
	services "codebase/internal/service"
	"codebase/pkg/mail"
)

func (service *ServiceUtils) SendMailActiveCode(ctx context.Context, email string, otpCode string) error {
	emailUsername := os.Getenv("EMAIL_USERNAME")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	auth := sasl.NewPlainClient("", emailUsername, emailPassword)

	emailNotice := storyW.Email{
		From:    emailUsername,
		To:      []string{email},
		Subject: storyW.MailActiveSubject,
		Body:    fmt.Sprintf("Mã kích hoạt tài khoản của bạn là: %s", otpCode),
	}
	err := mail.SendMail(&emailNotice, auth, services.ADDRESS_SMTP)
	if err != nil {
		return err
	}

	return err
}
