package spi

type Mailer interface {
	Send(message string) error
}
