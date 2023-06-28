package app_validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AppValidator struct {
	log       *zap.Logger
	Validator *validator.Validate
}

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return true
		}
		return false
	})

	return validate
}

func NewAppValidator(log *zap.Logger, validator *validator.Validate) *AppValidator {
	return &AppValidator{
		log:       log,
		Validator: validator,
	}
}

// ValidationErrors func for show validation errors for each invalid fields.
func (validate *AppValidator) ValidationErrors(err error) map[string]any {
	// Define fields map.
	fields := map[string]any{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		errMsg := fmt.Sprintf("validation failed on '%s' tag", err.Tag())
		param := err.Param()
		if param != "" {
			errMsg = fmt.Sprintf("%s. allow: %s", errMsg, param)
		}

		fields[strings.ToLower(err.Field())] = errMsg
	}

	return fields
}

// Validate func for validate model fields.
func (validate *AppValidator) Validate(model any) map[string]any {
	// Validate model fields.
	err := validate.Validator.Struct(model)
	if err != nil {
		return validate.ValidationErrors(err)
	}

	return nil
}

// ValidateUUID func for validate uuid.UUID fields.
func (validate *AppValidator) ValidateUUID(field string, value string) map[string]any {
	// Validate uuid.UUID fields.
	err := validate.Validator.Var(value, "uuid")
	if err != nil {
		return map[string]any{field: "invalid uuid"}
	}

	return nil
}

// ValidateEmail func for validate email fields.
func (validate *AppValidator) ValidateEmail(field string, value string) map[string]any {
	// Validate email fields.
	err := validate.Validator.Var(value, "email")
	if err != nil {
		return map[string]any{field: "invalid email"}
	}

	return nil
}

// ValidatePassword func for validate password fields.
func (validate *AppValidator) ValidatePassword(field string, value string) map[string]any {
	// Validate password fields.
	err := validate.Validator.Var(value, "password")
	if err != nil {
		return map[string]any{field: "invalid password"}
	}

	return nil
}

var ServiceModule = fx.Options(
	fx.Provide(NewValidator),
	fx.Provide(NewAppValidator),
)
