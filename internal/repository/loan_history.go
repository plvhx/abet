package repository

import (
    "context"
    "encoding/json"

    sq "github.com/Masterminds/squirrel"
    "abet/internal/model/db"
    "abet/pkg"
    "abet/pkg/dbase"
)

type LoanHistoryRepository struct {
    *pkg.Options
}

func NewLoanHistoryRepository(options *pkg.Options) *LoanHistoryRepository {
    return &LoanHistoryRepository{options}
}

func (l *LoanHistoryRepository) Create(
    ctx context.Context,
    data db.LoanHistory,
) error {
    params := sq.Eq{
        "loanID": data.LoanID,
        "status": data.Status,
        "staffName": data.StaffName,
        "staffMail": data.StaffMail,
    }

    if data.Extra != nil {
        buff, _ := json.Marshal(data.Extra)
        params["extra"] = buff
    }

    query, args := sq.
        Insert("loanHistory").
        SetMap(params).
        PlaceholderFormat(sq.Dollar).
        MustSql()

    var conn dbase.SQLQueryExec = l.Db

    if tx := dbase.GetTxFromContext(ctx); tx != nil {
        conn = tx
    }

    _, err := conn.ExecContext(ctx, query, args...)

    if err != nil {
        return err
    }

    return nil
}
