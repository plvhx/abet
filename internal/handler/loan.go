package handler

import (
    "net/http"
    "path/filepath"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"

    "abet/internal/model/payload"
    coreError "abet/pkg/error"
)

func (h *Handler) CreateLoan(c echo.Context) error {
    var (
        ctx = c.Request().Context()
        body payload.CreateLoan
    )

    if err := c.Bind(&body); err != nil {
        return err
    }

    if err := c.Validate(&body); err != nil {
        return err
    }

    res, err := h.Service.CreateLoan(ctx, body)

    if err != nil {
        return err
    }

    return c.JSON(http.StatusCreated, payload.ResponseData{Data: res})
}

func (h *Handler) GetAllLoan(c echo.Context) error {
    var (
        ctx = c.Request().Context()
        param payload.LoanPaginationFilter
    )

    if err := c.Bind(&param); err != nil {
        return err
    }

    param.Normalize()

    items, count, err := h.Service.PaginateLoan(ctx, param)

    if err != nil {
        return err
    }

    return c.JSON(http.StatusOK, param.Paginate(items, count))
}

func (h *Handler) GetLoanById(c echo.Context) error {
    var (
        ctx = c.Request().Context()
        param = c.Param("loanId")
    )

    loanId, err := uuid.Parse(param)

    if err != nil {
        return coreError.Errorf(
            http.StatusBadRequest,
            "loan handler: \"%s\" is not valid UUID.",
            param,
        )
    }

    detail, err := h.Service.GetLoanById(ctx, loanId)

    if err != nil {
        return err
    }

    return c.JSON(http.StatusOK, payload.ResponseData{Data: detail})
}

func (h *Handler) ApproveLoan(c echo.Context) error {
    var (
        ctx = c.Request().Context()
        param = c.Param("loanId")
        body payload.ApproveLoan
    )

    loanId, err := uuid.Parse(param)

    if err != nil {
        return coreError.Errorf(
            http.StatusBadRequest,
            "loan handler: \"%s\" is not valid UUID.",
            param,
        )
    }

    body.LoanID = loanId

    if err = c.Bind(&body); err != nil {
        return err
    }

    if err = c.Validate(&body); err != nil {
        return err
    }

    if err = h.Service.ApproveLoan(ctx, body); err != nil {
        return err
    }

    return c.NoContent(http.StatusOK)
}

func (h *Handler) DisburseLoan(c echo.Context) error {
    var (
        ctx = c.Request().Context()
        param = c.Param("loanId")
        body payload.DisburseLoan
    )

    loanId, err := uuid.Parse(param)

    if err != nil {
        return coreError.Errorf(
            http.StatusBadRequest,
            "loan handler: \"%s\" is not valid UUID.",
            param,
        )
    }

    file, err := c.FormFile("signedAgreementDocument")

    if err != nil {
        return err
    }

    src, err := file.Open()

    if err != nil {
        return err
    }

    defer src.Close()

    body.LoanID = loanId
    body.SignedAgreementDocument = &src
    body.DocumentExtension = filepath.Ext(file.Filename)

    if err := c.Bind(&body); err != nil {
        return err
    }

    if err = c.Validate(&body); err != nil {
        return err
    }

    if err = h.Service.DisburseLoan(ctx, body); err != nil {
        return err
    }

    return c.NoContent(http.StatusOK)
}

func (h *Handler) InvestLoan(c echo.Context) error {
    var (
        ctx = c.Request().Context()
        param = c.Param("loanId")
        body payload.InvestLoan
    )

    loanId, err := uuid.Parse(param)

    if err != nil {
        return coreError.Errorf(
            http.StatusBadRequest,
            "loan handler: \"%s\" is not valid UUID.",
            param,
        )
    }

    if err = c.Bind(&body); err != nil {
        return err
    }

    body.LoanID = loanId

    if err = c.Validate(&body); err != nil {
        return err
    }

    if err = h.Service.InvestLoan(ctx, body); err != nil {
        return err
    }

    return c.NoContent(http.StatusOK)
}
