package spi

type Validator interface {
	Validate(interface{}) (map[string]map[string]string, error)
}
