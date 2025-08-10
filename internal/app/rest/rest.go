package rest

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    echoFw "github.com/labstack/echo/v4"

    "abet/internal/app"
    "abet/pkg"
    coreConfig "abet/pkg/config"
)

type Server struct {
    echo *echoFw.Echo
    config *coreConfig.Config
    options *pkg.Options
}

func NewHttpServer() Server {
    var (
        appConfig = coreConfig.GetConfig()
        appContext = app.Context{Config: appConfig}
        appLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
    )

    slog.SetDefault(appLogger)

    server := Server{
        echo: echoFw.New(),
        config: &appConfig,
        options: &pkg.Options{
            Config: appConfig,
            Db: appContext.GetDB(),
            Bucket: appContext.GetS3BucketClient(),
            MailClient: appContext.GetGoMailerClient(),
        },
    }

    server.initRouter()

    return server
}

func (s *Server) Start() {
    go func() {
        err := s.echo.Start(fmt.Sprintf(":%d", s.config.AppPort))

        if err != nil && !errors.Is(err, http.ErrServerClosed) {
            s.echo.Logger.Fatal("server: fatally shutting down")
        }
    }()

    // declare buffered channel with capacity 1 to handle
    // either SIGINT or SIGTERM. no need to allocate more than
    // 1, because multiple signal never be raised in that
    // application context.
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    // tell running server goroutine to respond to raised
    // signal stack by anonymously receive value from previously
    // declared buffered channel.
    <-quit
}

func (s *Server) Stop() {
    err := s.echo.Shutdown(context.Background())

    if err != nil {
        s.echo.Logger.Fatal(err)
    }

    if s.options.Db != nil {
        _ = s.options.Db.Close()
    }

    slog.Info("server: gracefully shutting down")
}
