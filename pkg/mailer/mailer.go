package mailer

import "log"

type Mailer struct{}

func NewMailer() *Mailer {
	return &Mailer{}
}

func (m *Mailer) Send(message string) error {
	log.Println(message)
	return nil
}
