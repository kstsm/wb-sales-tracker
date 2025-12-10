package validator

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
)

type Validate struct {
	*validator.Validate
}

func NewValidator() *Validate {
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

	return &Validate{Validate: validate}
}

func (v *Validate) FormatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var messages []string
		for _, e := range validationErrors {
			msg := e.Error()
			if idx := strings.Index(msg, "Error:"); idx != -1 {
				msg = strings.TrimSpace(msg[idx+len("Error:"):])
			}
			msg = strings.TrimPrefix(msg, "Field ")
			messages = append(messages, msg)
		}
		if len(messages) > 0 {
			return messages[0]
		}
	}
	return err.Error()
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
