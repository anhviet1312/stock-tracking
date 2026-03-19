package models

const (
	MailActiveSubject        = "Kích hoạt tài khoản"
	ADDRESS_IMAP      string = "imap.yandex.com:993"
)

type Email struct {
	From    string
	To      []string
	Subject string
	Body    string
}
