package app

import (
    "fmt"
    "log/slog"
    "time"

    _ "github.com/jackc/pgx/stdlib"
    "github.com/jmoiron/sqlx"

    "abet/pkg/config"
    "abet/pkg/files"
    "abet/pkg/mailer"
)

type Context struct {
    Config config.Config
}

func (c Context) GetDB() *sqlx.DB {
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        c.Config.DBUser,
        c.Config.DBPassword,
        c.Config.DBHost,
        c.Config.DBPort,
        c.Config.DBSchema,
    )

    db, err := sqlx.Open("pgx", dsn)

    if err != nil {
        slog.Error(
            "sqlx error: failed to open connection",
            slog.String("err", err.Error()),
            slog.String("dsn", dsn),
        )

        panic(err)
    }

    db.SetMaxIdleConns(c.Config.DBMaxIdleConn)
    db.SetMaxOpenConns(c.Config.DBMaxOpenDescriptor)
    db.SetConnMaxLifetime(time.Minute * 10)

    err = db.Ping()

    if err != nil {
        slog.Error("sqlx error: ping error", slog.String("err", err.Error()))
    }

    return db
}

func (c Context) GetS3BucketClient() files.S3Client {
    return files.NewS3Storage(
        c.Config.AWSAccessKey,
        c.Config.AWSSecretKey,
        c.Config.AWSBucketName,
        c.Config.AWSEndpoint,
        c.Config.AWSRegion,
    )
}

func (c Context) GetGoMailerClient() mailer.GoMailer {
    return mailer.NewGoMailClient(
        c.Config.SMTPUser,
        c.Config.SMTPPassword,
        c.Config.SMTPHost,
        c.Config.SMTPPort,
    )
}
