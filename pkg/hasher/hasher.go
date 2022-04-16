package hasher

import (
	"crypto/sha256"
	"fmt"
)

func Hash(s interface{}) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(fmt.Sprintf("%v", s)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
