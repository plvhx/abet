package custom

import (
    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
)

type EchoValidator struct {
    validator *validator.Validate
}

func (ev EchoValidator) Validate(i interface{}) error {
    if err := ev.validator.Struct(i); err != nil {
        return err
    }

    return nil
}

func NewValidator() echo.Validator {
    val := validator.New()
    return EchoValidator{validator: val}
}
