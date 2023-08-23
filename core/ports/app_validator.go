package ports

import "context"

type AppValidator interface {
	Validate(ctx context.Context, i interface{}) error

	ValidateUUID(ctx context.Context, field string, value string) map[string]interface{}

	ValidationErrors(ctx context.Context, err error) map[string]interface{}

	ValidateEmail(ctx context.Context, field string, value string) map[string]interface{}

	ValidatePassword(ctx context.Context, field string, value string) map[string]interface{}
}
