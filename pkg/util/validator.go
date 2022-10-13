package util

import (
	"context"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init()  {
	validate = validator.New()
}

func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}
