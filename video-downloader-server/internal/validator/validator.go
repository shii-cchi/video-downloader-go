package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func Init() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("foldername", folderNameValidation)
	validate.RegisterValidation("videoname", videoNameValidation)
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

func videoNameValidation(fl validator.FieldLevel) bool {
	videoName := fl.Field().String()
	if len(videoName) < 1 || len(videoName) > 100 {
		return false
	}

	return true
}

func objectIDValidation(fl validator.FieldLevel) bool {
	return !fl.Field().IsZero()
}
