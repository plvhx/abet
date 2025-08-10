package dbase

import (
    "context"
    "database/sql"
    "log/slog"

    "github.com/jmoiron/sqlx"
)

type (
    ctxKey string

    SQLExec interface {
        sqlx.Execer
        sqlx.ExecerContext

        NamedExec(query string, arg interface{}) (sql.Result, error)
        NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
    }

    SQLQuery interface {
        sqlx.Queryer

        GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
        SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
        PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
    }

    SQLQueryExec interface {
        SQLExec
        SQLQuery
    }

    WrapTransactionFunc func(ctx context.Context) error
)

func BeginTransaction(ctx context.Context, db *sqlx.DB, fn WrapTransactionFunc, isolations ...sql.IsolationLevel) error {
    isolationLevel := sql.LevelRepeatableRead

    if len(isolations) >= 1 {
        isolationLevel = isolations[0]
    }

    var (
        tx = GetTxFromContext(ctx)
        err error
    )

    if tx == nil {
        tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: isolationLevel})

        if err != nil {
            return err
        }

        defer func() {
            if p := recover(); p != nil {
                if rollbackErr := tx.Rollback(); rollbackErr != nil {
                    slog.Warn("dbase: failed rollback on panic")
                }

                panic(p)
            } else if err != nil {
                if rollbackErr := tx.Rollback(); rollbackErr != nil {
                    slog.Warn("dbase: failed rollback on error")
                }
            } else {
                if commitErr := tx.Commit(); commitErr != nil {
                    slog.Warn("dbase: failed to commit")
                }
            }
        }()
    }

    ctx = context.WithValue(ctx, txKey, tx)
    err = fn(ctx)

    if err != nil {
        return err
    }

    return nil
}

const txKey ctxKey = "IsTransaction"

func GetTxFromContext(ctx context.Context) *sqlx.Tx {
    if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
        return tx
    }

    return nil
}
