package app_validator

import (
	"auth/core/ports"
	"auth/internal"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AppValidator struct {
	log       *zap.Logger
	config    *internal.AppConfig
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

func NewAppValidator(log *zap.Logger, config *internal.AppConfig, validator *validator.Validate) (ports.AppValidator, error) {
	return &AppValidator{
		log:       log,
		config:    config,
		Validator: validator,
	}, nil
}

// ValidationErrors func for show validation errors for each invalid fields.
func (svc *AppValidator) ValidationErrors(ctx context.Context, err error) map[string]any {
	svc.log.Info("validation errors")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.app_validator.validation_errors")
	defer span.End()

	// Define fields map.
	fields := map[string]any{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		errMsg := fmt.Sprintf("validation failed on '%s' field", err.Tag())
		param := err.Param()
		if param != "" {
			errMsg = fmt.Sprintf("%s. allow: %s", errMsg, param)
		}

		fields[strings.ToLower(err.Field())] = errMsg
	}

	return fields
}

// Validate func for validate model fields.
func (svc *AppValidator) Validate(ctx context.Context, model interface{}) error {
	svc.log.Info("validate")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.app_validator.validate")
	defer span.End()

	return svc.Validator.Struct(model)
}

// ValidateUUID func for validate uuid.UUID fields.
func (svc *AppValidator) ValidateUUID(ctx context.Context, field string, value string) map[string]any {
	svc.log.Info("uuid from string")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.app_validator.validate_uuid")
	defer span.End()

	// Validate uuid.UUID fields.
	err := svc.Validator.Var(value, "uuid")
	if err != nil {
		return map[string]any{field: "invalid uuid"}
	}

	return nil
}

// ValidateEmail func for validate email fields.
func (svc *AppValidator) ValidateEmail(ctx context.Context, field string, value string) map[string]any {
	svc.log.Info("validate email")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.app_validator.validate_email")
	defer span.End()

	// Validate email fields.
	err := svc.Validator.Var(value, "email")
	if err != nil {
		return map[string]any{field: "invalid email"}
	}

	return nil
}

// ValidatePassword func for validate password fields.
func (svc *AppValidator) ValidatePassword(ctx context.Context, field string, value string) map[string]any {
	svc.log.Info("validate password")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.validate.validate_password")
	defer span.End()

	// Validate password fields.
	err := svc.Validator.Var(value, "password")
	if err != nil {
		return map[string]any{field: "invalid password"}
	}

	return nil
}

var ServiceModule = fx.Options(
	fx.Provide(NewValidator),
	fx.Provide(NewAppValidator),
)
