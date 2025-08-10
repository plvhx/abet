package custom

import (
    "fmt"
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    "abet/internal/model/payload"
    coreError "abet/pkg/error"
)

func NewErrorHandler(err error, c echo.Context) {
    var (
        code = http.StatusInternalServerError
        message = "Internal Server Error"
        errs []coreError.FieldError
    )

    switch t := err.(type) {
    case coreError.CustomError:
        message = t.Message
        code = t.StatusCode
    case validator.ValidationErrors:
        code = http.StatusBadRequest
        message = "Validation Error"

        for _, verr := range t {
            errs = append(errs, coreError.FieldError{Field: verr.Field(), Reason: messageForTag(verr) })
        }
    case *echo.HTTPError:
        code = t.Code
        message = fmt.Sprint(t.Message)
    default:
        message = t.Error()
    }

    resp := payload.ResponseError[coreError.CustomError]{
        Error: coreError.CustomError{
            StatusCode: code,
            Message: message,
        },
    }

    if len(errs) > 0 {
        resp.Error.Errors = errs
    }

    _ = c.JSON(code, resp)
}

func messageForTag(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return "This field is required"
    case "email":
        return "Invalid email"
    case "max":
        return fmt.Sprintf("Reach maximum %s", fe.Param())
    case "min":
        return fmt.Sprintf("Does not meet minimum %s", fe.Param())
    }

    return fe.Error()
}
