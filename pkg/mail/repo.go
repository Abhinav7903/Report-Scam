package mail

type Repository interface {
	SendMail(to, subject, body string) error
}
