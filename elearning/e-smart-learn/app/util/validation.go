package util

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var trans ut.Translator

// InitValidator must be called in main()
func InitValidator() error {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// use json tag to get field name in error message
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// locale english
		enLocale := en.New()
		uni := ut.New(enLocale, enLocale)

		t, _ := uni.GetTranslator("en")
		trans = t

		return enTranslations.RegisterDefaultTranslations(v, trans)

		/*
			// locale vietnamese
			viLocale := vi.New()
			uni := ut.New(viLocale, viLocale)

			t, _ := uni.GetTranslator("vi")
			trans = t

			return viTranslations.RegisterDefaultTranslations(v, trans)
		*/

	}

	return nil
}

func GetTranslator() ut.Translator {
	return trans
}

func IsValidUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
