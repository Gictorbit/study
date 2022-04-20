package main

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func main() {
	err := InitValidator()
	if err != nil {
		panic(err)
	}
	myName := "viktor viktor"
	err = validate.Var(myName, "required,max=5")
	if err != nil {
		fmt.Println("name" + err.(validator.ValidationErrors).Translate(trans)[""])
	}
}

func InitValidator() error {
	if validate != nil {
		return nil
	}
	validate = validator.New()
	enlang := en.New()
	uni := ut.New(enlang, enlang)
	trans, _ = uni.GetTranslator("en")
	err := en_translations.RegisterDefaultTranslations(validate, trans)
	return err
}
