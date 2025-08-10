package rest

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"

    coreHandler "abet/internal/handler"
    coreRepository "abet/internal/repository"
    coreService "abet/internal/service"

    "abet/pkg"
    "abet/pkg/custom"
)

func (s *Server) initRouter() {
    s.echo.Validator = custom.NewValidator()
    s.echo.HTTPErrorHandler = custom.NewErrorHandler

    s.echo.Use(middleware.RequestID())
    s.echo.Use(middleware.Logger())
    s.echo.Use(middleware.Recover())
    s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowHeaders: []string{
            echo.HeaderOrigin,
            echo.HeaderContentType,
            echo.HeaderAcceptEncoding,
        },
    }))

    // register health check endpoint.
    s.echo.GET("/health", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{"message": "All OK."})
    })

    var (
        repository = coreRepository.NewRepository(
            &pkg.Options{
                Config: s.options.Config,
                Db: s.options.Db,
            },
        )

        service = coreService.NewService(
            s.options,
            repository,
        )

        handler = coreHandler.NewHandler(service)
    )

    v1 := s.echo.Group("/v1")

    route := v1.Group("/loan")

    route.POST("", handler.CreateLoan)
    route.GET("", handler.GetAllLoan)
    route.GET("/:loanId", handler.GetLoanById)
    route.PATCH("/:loanId/approval", handler.ApproveLoan)
    route.PATCH("/:loanId/disbursement", handler.DisburseLoan)
    route.PATCH("/:loanId/investment", handler.InvestLoan)
}
