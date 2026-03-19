package mail

import (
	"fmt"
	"strings"
	"time"

	"codebase/internal/models"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

func SendMail(email *models.Email, auth sasl.Client, addressSmtp string) error {
	// Set up header
	headers := make(map[string]string)
	headers["From"] = email.From
	headers["To"] = strings.Join(email.To, ",")
	headers["Subject"] = email.Subject
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + email.Body)

	err := smtp.SendMail(addressSmtp, auth, email.From, email.To, strings.NewReader(msg.String()))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
