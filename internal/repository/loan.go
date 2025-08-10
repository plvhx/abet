package repository

import (
    "context"
    "net/http"

    sq "github.com/Masterminds/squirrel"
    "github.com/google/uuid"

    "abet/internal/model/db"
    "abet/internal/model/payload"
    "abet/pkg"
    "abet/pkg/dbase"
    coreError "abet/pkg/error"
    "abet/pkg/types"
)

type LoanRepository struct {
    *pkg.Options
}

func NewLoanRepository(options *pkg.Options) *LoanRepository {
    return &LoanRepository{options}
}

func (l *LoanRepository) Paginate(
    ctx context.Context,
    param payload.LoanPaginationFilter,
) (loans []db.Loan, total int64, err error) {
    loans = make([]db.Loan, 0)
    columns := []string{
        "id", "name", "status", "borrowerId",
        "borrowerName", "borrowerRate", "borrowerAgreementUrl", "principalAmount",
        "totalInvestedAmount", "roiRate", "staffName", "staffEmail",
        "staffId", "staffVisitProof", "approvedAt", "disbursedAt",
        "createdAt", "updatedAt",
    }

    query, args := l.
        queryAll(param, columns...).
        Limit(uint64(param.Limit)).
        Offset(uint64((param.Page - 1) * param.Limit)).
        OrderBy("createdAt DESC").
        MustSql()

    queryCount, argsCount := l.queryAll(param, "COUNT(id) AS total").MustSql()

    var conn dbase.SQLQueryExec = l.Db

    if tx := dbase.GetTxFromContext(ctx); tx != nil {
        return loans, total, err
    }

    if err = conn.SelectContext(ctx, &loans, query, args...); err != nil {
        return loans, total, err
    }

    if err = conn.GetContext(ctx, &total, queryCount, argsCount...); err != nil {
        return loans, total, err
    }

    return loans, total, nil
}

func (l *LoanRepository) GetById(
    ctx context.Context,
    id uuid.UUID,
) (res db.Loan, err error) {
    columns := []string{
        "id", "name", "status", "borrowerId",
        "borrowerName", "borrowerRate", "borrowerAgreementUrl", "principalAmount",
        "totalInvestedAmount", "roiRate", "staffName", "staffEmail",
        "staffId", "staffVisitProof", "approvedAt", "disbursedAt",
        "createdAt", "updatedAt",
    }

    query, args := sq.
        Select(columns...).
        From("loans").
        Where("id = ?", id).
        PlaceholderFormat(sq.Dollar).
        MustSql()

    var conn dbase.SQLQueryExec = l.Db

    if tx := dbase.GetTxFromContext(ctx); tx != nil {
        conn = tx
    }

    if err = conn.GetContext(ctx, &res, query, args...); err != nil {
        return res, err
    }

    return res, nil
}

func (l *LoanRepository) Create(
    ctx context.Context,
    data payload.CreateLoan,
) (lastInsertedID uuid.UUID, err error) {
    query, args := sq.Insert("loans").SetMap(sq.Eq{
        "name": data.Name,
        "status": types.StatusProposed,
        "borrowerId": data.BorrowerID,
        "borrowerName": data.BorrowerName,
        "borrowerEmail": data.BorrowerEmail,
        "borrowerRate": data.BorrowerRate,
        "principalAmount": data.PrincipalAmount,
        "roiRate": data.ROIRate,
        "staffId": data.StaffID,
        "staffName": data.StaffName,
        "staffEmail": data.StaffEmail,
    }).PlaceholderFormat(sq.Dollar).Suffix("RETURNING id").MustSql()

    var conn dbase.SQLQueryExec = l.Db

    if tx := dbase.GetTxFromContext(ctx); tx != nil {
        conn = tx
    }

    if err = conn.GetContext(ctx, &lastInsertedID, query, args...); err != nil {
        return uuid.Nil, err
    }

    return lastInsertedID, err
}

func (l *LoanRepository) Update(
    ctx context.Context,
    data db.Loan,
    id uuid.UUID,
) error {
    params := sq.Eq{
        "status": data.Status,
        "updatedAt": sq.Expr("NOW()"),
    }

    if data.StaffName != "" {
        params["staffName"] = data.StaffName
    }

    if data.StaffEmail != "" {
        params["staffEmail"] = data.StaffEmail
    }

    if data.BorrowerAgreementURL.Valid && data.BorrowerAgreementURL.String != "" {
        params["borrowerAgreementUrl"] = data.BorrowerAgreementURL
    }

    if data.TotalInvestedAmount > 0 {
        params["totalInvestedAmount"] = data.TotalInvestedAmount
    }

    if data.ApprovedAt.Valid && !data.ApprovedAt.IsZero() {
        params["approvedAt"] = data.ApprovedAt
    }

    if data.DisbursedAt.Valid && !data.DisbursedAt.IsZero() {
        params["disbursedAt"] = data.DisbursedAt
    }

    query, args := sq.
        Update("loans").
        SetMap(params).
        Where("id = ?", id).
        PlaceholderFormat(sq.Dollar).
        MustSql()

    var conn dbase.SQLQueryExec = l.Db

    if tx := dbase.GetTxFromContext(ctx); tx != nil {
        conn = tx
    }

    row, err := conn.ExecContext(ctx, query, args...)

    if err != nil {
        return err
    }

    updated, err := row.RowsAffected()

    if err != nil {
        return err
    }

    if updated == 0 {
        return coreError.Error(http.StatusUnprocessableEntity, "update loan failed")
    }

    return nil
}

func (l *LoanRepository) queryAll(
    param payload.LoanPaginationFilter,
    columns ...string,
) sq.SelectBuilder {
    query := sq.
        Select(columns...).
        From("loans").
        PlaceholderFormat(sq.Dollar)

    if param.Search != "" {
        query = query.Where(sq.ILike{"name": param.Search + "%"})
    }

    if param.Status.String() != "" {
        query = query.Where("status = ?", param.Status)
    }

    return query
}
