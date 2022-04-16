package validator

import (
	"github.com/go-playground/validator"
)

func Validate(s interface{}) (map[string]map[string]string, error) {
	v := validator.New()
	m := map[string]map[string]string{}

	if err := v.Struct(s); err != nil {

		for _, v := range err.(validator.ValidationErrors) {
			m[v.StructField()] = map[string]string{}
			m[v.StructField()]["tag"] = v.Tag()
		}
	}
	return m, nil
}
