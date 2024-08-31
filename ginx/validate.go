package ginx

import (
	"errors"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	_DefaultValidate = validator.New()
)

func Validate[T any](item *T) (err error) {
	var (
		errs   validator.ValidationErrors
		fields []string
	)

	toLower := func(field string) string {
		var runes []rune

		runes = []rune(field)
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	}

	if err = _DefaultValidate.Struct(item); err != nil {
		errs = err.(validator.ValidationErrors)
		fields = make([]string, len(errs))
		// fmt.Printf("==> error: %+v, filed: %q\n", err, field)

		for i := range errs {
			fields[i] = toLower(errs[i].Field())
		}

		return errors.New(strings.Join(fields, ","))
	}

	return nil
}
