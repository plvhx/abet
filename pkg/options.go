package pkg

import (
    "github.com/jmoiron/sqlx"

    "abet/pkg/config"
    "abet/pkg/files"
    "abet/pkg/mailer"
)

type Options struct {
    Config     config.Config
    Db         *sqlx.DB
    Bucket     files.Bucket
    MailClient mailer.EmailClient
}
