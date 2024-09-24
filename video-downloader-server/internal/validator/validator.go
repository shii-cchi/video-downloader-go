package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func Init() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("foldername", folderNameValidation)
	validate.RegisterValidation("objectid", objectIDValidation)

	return validate
}

func folderNameValidation(fl validator.FieldLevel) bool {
	folderName := fl.Field().String()
	if len(folderName) < 1 || len(folderName) > 20 {
		return false
	}

	re := regexp.MustCompile(`^[\w\s-]+$`)
	return re.MatchString(folderName)
}

func objectIDValidation(fl validator.FieldLevel) bool {
	return !fl.Field().IsZero()
}
