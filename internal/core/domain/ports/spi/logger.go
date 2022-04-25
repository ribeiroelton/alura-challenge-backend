package spi

type Logger interface {
	Info(...interface{})
	Error(...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}
