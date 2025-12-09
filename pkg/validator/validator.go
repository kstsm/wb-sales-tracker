package validator

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
)

func NewValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("rfc3339", ValidateRFC3339); err != nil {
		slog.Fatal("Failed to register rfc3339 validation", "error", err)
		os.Exit(1)
	}
	if err := validate.RegisterValidation("item_type", ValidateItemType); err != nil {
		slog.Fatal("Failed to register item_type validation", "error", err)
		os.Exit(1)
	}
	if err := validate.RegisterValidation("sort_by", ValidateSortBy); err != nil {
		slog.Fatal("Failed to register sort_by validation", "error", err)
		os.Exit(1)
	}
	if err := validate.RegisterValidation("sort_order", ValidateSortOrder); err != nil {
		slog.Fatal("Failed to register sort_order validation", "error", err)
		os.Exit(1)
	}

	return validate
}

func ValidateRFC3339(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}

func ValidateItemType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "income" || value == "expense"
}

func ValidateSortBy(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "date" || value == "amount" || value == "category"
}

func ValidateSortOrder(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "asc" || value == "desc"
}
