package mail

import (
	"log"

	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/samber/do"

	storyW "codebase/internal/models"
)

type ServiceMail struct {
	mail *imapclient.Client
}

func NewServiceMail(container *do.Injector) (*ServiceMail, error) {
	c, err := imapclient.DialTLS(storyW.ADDRESS_IMAP, nil)
	if err != nil {
		log.Printf("failed to dial IMAP yandex: %v\n", err)
		return nil, err
	}
	return &ServiceMail{c}, nil
}

func (m *ServiceMail) Close() error {
	if m.mail != nil {
		return m.mail.Close()
	}
	return nil
}
