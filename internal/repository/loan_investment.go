package repository

import (
    "context"

    sq "github.com/Masterminds/squirrel"
    "github.com/google/uuid"

    "abet/internal/model/db"
    "abet/pkg"
    "abet/pkg/dbase"
)

type LoanInvestmentRepository struct {
    *pkg.Options
}

func NewLoanInvestmentRepository(options *pkg.Options) *LoanInvestmentRepository {
    return &LoanInvestmentRepository{options}
}

func (l *LoanInvestmentRepository) GetLoanInvestors(
    ctx context.Context,
    loanId uuid.UUID,
) ([]db.LoanInvestment, error) {
    results := make([]db.LoanInvestment, 0)
    columns := []string{
        "id", "loanID", "investorID",
        "name", "email", "amount",
    }

    query, args := sq.
        Select(columns...).
        From("loanInvestment").
        Where("loanID = ?", loanId).
        PlaceholderFormat(sq.Dollar).
        MustSql()

    var conn dbase.SQLQueryExec = l.Db

    if tx := dbase.GetTxFromContext(ctx); tx != nil {
        conn = tx
    }

    if err := conn.SelectContext(ctx, &results, query, args...); err != nil {
        return results, err
    }

    return results, nil
}

func (l *LoanInvestmentRepository) Create(ctx context.Context, data db.LoanInvestment) error {
    params := sq.Eq{
        "loanID": data.LoanID,
        "investorID": data.InvestorID,
        "name": data.Name,
        "email": data.Email,
        "amount": data.Amount,
    }

    if data.AgreementURL != "" {
        params["agreementUrl"] = data.AgreementURL
    }

    query, args := sq.
        Insert("loanInvestment").
        SetMap(params).
        PlaceholderFormat(sq.Dollar).
        MustSql()

    var conn dbase.SQLQueryExec = l.Db

    if tx := dbase.GetTxFromContext(ctx); tx != nil {
        conn = tx
    }

    if _, err := conn.ExecContext(ctx, query, args...); err != nil {
        return err
    }

    return nil
}

func (l *LoanInvestmentRepository) SetAgreementURL(ctx context.Context, agreementPath string, id uuid.UUID) error {
    query, args := sq.
        Update("loanInvestment").
        SetMap(sq.Eq{
            "agreementUrl": agreementPath,
            "updatedAt": sq.Expr("NOW()"),
        }).
        Where("id = ?", id).
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
