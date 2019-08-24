package validator

import (
	"github.com/asaskevich/govalidator"
	"regexp"
)

func init() {
	govalidator.CustomTypeTagMap.Set("phone", func(i interface{}, o interface{}) bool {
		r := regexp.MustCompile(`^\+\d{5,20}$`)
		switch v := i.(type) {
		case string:
			return r.MatchString(v)
		}
		return false
	})
}
