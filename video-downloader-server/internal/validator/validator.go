package validator

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

func Init() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("foldername", folderNameValidation)
	validate.RegisterValidation("objectid", objectIDValidation)

	return validate
}

func folderNameValidation(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[\w\s-]+$`)
	return re.MatchString(fl.Field().String())
}

func objectIDValidation(fl validator.FieldLevel) bool {
	_, err := primitive.ObjectIDFromHex(fl.Field().String())
	return err == nil
}
