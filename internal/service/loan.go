package service

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "time"
    "database/sql"
    "log/slog"
    "net/http"
    "path/filepath"

    "github.com/google/uuid"
    "github.com/guregu/null"

    "golang.org/x/sync/errgroup"

    "abet/internal/model/db"
    "abet/internal/model/payload"
    "abet/pkg/dbase"
    coreError "abet/pkg/error"
    "abet/pkg/files"
    "abet/pkg/mailer"
    "abet/pkg/types"
)

const (
    goroutineLimit = 50
)

func (s *Service) CreateLoan(ctx context.Context, data payload.CreateLoan) (db.Loan, error) {
    var detail db.Loan

    err := dbase.BeginTransaction(ctx, s.Db, func(ctx context.Context) error {
        loanId, err := s.Repository.Loan.Create(ctx, data)

        if err != nil {
            slog.Error("loan service: error while creating loan", slog.String("err", err.Error()))

            return coreError.Error(
                http.StatusInternalServerError,
                "loan service: internal server error.",
            )
        }

        history := db.LoanHistory{
            LoanID: loanId,
            Status: types.StatusProposed,
            StaffID: data.StaffID,
            StaffMail: data.StaffEmail,
            StaffName: data.StaffName,
        }

        if err = s.Repository.LoanHistory.Create(ctx, history); err != nil {
            slog.Error("loan service: error while creating loan history", slog.String("err", err.Error()))

            return coreError.Error(
                http.StatusInternalServerError,
                "loan service: internal server error.",
            )
        }

        detail, err = s.Repository.Loan.GetById(ctx, loanId)

        if err != nil {
            slog.Error("loan service: error on get detail project", slog.String("err", err.Error()))

            return coreError.Error(
                http.StatusInternalServerError,
                "loan service: internal server error.",
            )
        }

        return nil
    })

    if err != nil {
        return detail, err
    }

    return detail, nil
}

func (s *Service) PaginateLoan(ctx context.Context, param payload.LoanPaginationFilter) (items []db.Loan, total int64, err error) {
    return s.Repository.Loan.Paginate(ctx, param)
}

func (s *Service) GetLoanById(ctx context.Context, id uuid.UUID) (db.Loan, error) {
    loan, err := s.Repository.Loan.GetById(ctx, id)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return db.Loan{}, coreError.Errorf(
                http.StatusNotFound,
                "Loan with ID: %s not found.",
                id.String(),
            )
        }

        slog.Error("loan service: error get loan by ID", slog.String("err", err.Error()))

        return db.Loan{}, coreError.Error(http.StatusInternalServerError, "internal server error")
    }

    if loan.BorrowerAgreementURL.Valid && loan.BorrowerAgreementURL.String != "" {
        agreementSignedURL, _ := s.Bucket.GetSignURL(ctx, loan.BorrowerAgreementURL.String)
        loan.BorrowerAgreementURL.SetValid(agreementSignedURL)
    }

    return loan, nil
}

func (s *Service) ApproveLoan(ctx context.Context, data payload.ApproveLoan) error {
    err := dbase.BeginTransaction(ctx, s.Db, func(ctx context.Context) error {
        loan, err := s.Repository.Loan.GetById(ctx, data.LoanID)

        if err != nil {
            if errors.Is(err, sql.ErrNoRows) {
                return coreError.Errorf(
                    http.StatusNotFound,
                    "Loan with ID: %s not found.",
                    data.LoanID.String(),
                )
            }

            slog.Error("loan service: error get loan by ID", slog.String("err", err.Error()))

            return coreError.Error(http.StatusInternalServerError, "internal server error")
        }

        if loan.Status != types.StatusProposed {
            return coreError.Error(
                http.StatusUnprocessableEntity,
                "loan service: only proposed project are allowed to be approved.",
            )
        }

        newLoan := db.Loan{
            Status: loan.Status.Next(),
            StaffID: data.StaffID,
            StaffEmail: data.StaffEmail,
            StaffName: data.StaffName,
            UpdatedAt: time.Now().UTC(),
            ApprovedAt: null.NewTime(time.Now().UTC(), true),
        }

        if err = s.Repository.Loan.Update(ctx, newLoan, data.LoanID); err != nil {
            slog.Error("loan service: error on update loan", slog.String("err", err.Error()))
            return err
        }

        nextHistory := db.LoanHistory{
            LoanID: data.LoanID,
            Status: newLoan.Status,
            StaffName: newLoan.StaffName,
            StaffID: newLoan.StaffID,
            StaffMail: newLoan.StaffEmail,
            UpdatedAt: time.Now().UTC(),
        }

        if err = s.Repository.LoanHistory.Create(ctx, nextHistory); err != nil {
            slog.Error("loan service: error on insert next history", slog.String("err", err.Error()))

            return coreError.Error(http.StatusInternalServerError, "internal server error")
        }

        return nil
    })

    if err != nil {
        return err
    }

    return nil
}

func (s *Service) DisburseLoan(ctx context.Context, data payload.DisburseLoan) error {
    ext := strings.ToLower(data.DocumentExtension)

    var extWhitelist = map[string]bool{
        ".png": true,
        ".jpg": true,
        ".pdf": true,
        ".jpeg": true,
        ".webp": true,
    }

    if !extWhitelist[ext] {
        return coreError.Error(http.StatusBadRequest, "file format is not supported")
    }

    err := dbase.BeginTransaction(ctx, s.Db, func(ctx context.Context) error {
        loan, err := s.Repository.Loan.GetById(ctx, data.LoanID)

        if err != nil {
            if errors.Is(err, sql.ErrNoRows) {
                return coreError.Errorf(
                    http.StatusNotFound,
                    "Loan with id %s not found",
                    data.LoanID,
                )
            }

            slog.Error("loan service: error on get project detail", slog.String("err", err.Error()))
            return coreError.Error(http.StatusInternalServerError, "internal server error")
        }

        if loan.Status != types.StatusInvested {
            return coreError.Error(
                http.StatusUnprocessableEntity,
                "disbursement is not allowed for non invested loans",
            )
        }

        filename := filepath.Join(
            "borrower",
            "agreement",
            fmt.Sprintf("project-%s%s", loan.ID.String(), data.DocumentExtension),
        )

        _, err = s.Bucket.Upload(ctx, filename, *data.SignedAgreementDocument)

        if err != nil {
            slog.Error("loan service: error on upload to bucket", slog.String("err", err.Error()))
            return coreError.Error(http.StatusInternalServerError, "internal server error")
        }

        newLoan := db.Loan{
            Status: loan.Status.Next(),
            UpdatedAt: time.Now().UTC(),
            BorrowerAgreementURL: null.StringFrom(filename),
            DisbursedAt: null.TimeFrom(time.Now().UTC()),
            StaffID: data.StaffID,
            StaffName: data.StaffName,
            StaffEmail: data.StaffMail,
        }

        if err = s.Repository.Loan.Update(ctx, newLoan, loan.ID); err != nil {
            slog.Error("loan service: error on update project", slog.String("err", err.Error()))
            return coreError.Error(http.StatusInternalServerError, "internal server error")
        }

        history := db.LoanHistory{
            LoanID: loan.ID,
            Status: newLoan.Status,
            StaffID: newLoan.StaffID,
            StaffName: newLoan.StaffName,
            StaffMail: newLoan.StaffEmail,
            Extra: map[string]string{
                "FieldVisitStaffID": data.FieldVisitStaffID.String(),
                "FieldVisitStaffName": data.FieldVisitStaffName,
                "FieldVisitStaffMail": data.FieldVisitStaffMail,
            },
        }

        if err = s.Repository.LoanHistory.Create(ctx, history); err != nil {
            slog.Error("loan service: error on create disbursed loan history", slog.String("err", err.Error()))
            return coreError.Error(http.StatusInternalServerError, "internal server error")
        }

        return nil
    })

    if err != nil {
        return err
    }

    return nil
}

func (s *Service) InvestLoan(ctx context.Context, data payload.InvestLoan) error {
    var totalInvestment float64

    err := dbase.BeginTransaction(ctx, s.Db, func(ctx context.Context) error {
        loan, err := s.Repository.Loan.GetById(ctx, data.LoanID)

        if err != nil {
            if errors.Is(err, sql.ErrNoRows) {
                return coreError.Errorf(
                    http.StatusNotFound,
                    "Loan with id %s not found.",
                    data.LoanID,
                )
            }

            slog.Error("loan service: error on get detail loan", slog.String("err", err.Error()))
            return coreError.Error(http.StatusInternalServerError, "internal server error")
        }

        if data.InvestmentAmount > loan.PrincipalAmount {
            return coreError.Error(http.StatusUnprocessableEntity, "out of quota for this loan")
        }

        if loan.Status == types.StatusInvested {
            return coreError.Error(http.StatusUnprocessableEntity, "this loan is fully funded")
        }

        if loan.Status != types.StatusApproved {
            return coreError.Error(http.StatusUnprocessableEntity, "investment only for the approved project")
        }

        totalInvestment = loan.TotalInvestedAmount + data.InvestmentAmount

        if totalInvestment > loan.PrincipalAmount {
            remainingAmount := loan.PrincipalAmount - loan.TotalInvestedAmount
            msg := fmt.Sprintf("out of quota, only %.2f remaining for this loan", remainingAmount)
            return coreError.Error(http.StatusUnprocessableEntity, msg)
        }

        newLoan := db.Loan{
            Status: loan.Status,
            UpdatedAt: time.Now().UTC(),
            TotalInvestedAmount: totalInvestment,
        }

        if totalInvestment == loan.PrincipalAmount {
            newLoan.Status = loan.Status.Next()
        }

        if err = s.Repository.Loan.Update(ctx, newLoan, data.LoanID); err != nil {
            return err
        }

        investment := db.LoanInvestment{
            LoanID: data.LoanID,
            InvestorID: data.InvestorID,
            Name: data.InvestorName,
            Email: data.InvestorEmail,
            Amount: data.InvestmentAmount,
            CreatedAt: time.Now().UTC(),
            UpdatedAt: time.Now().UTC(),
        }

        if err = s.Repository.LoanInvestment.Create(ctx, investment); err != nil {
            slog.Error("loan service: error on create investment", slog.String("err", err.Error()))
            return coreError.Error(http.StatusInternalServerError, "internal server error")
        }

        if newLoan.Status == types.StatusInvested {
            invested := db.LoanHistory{
                LoanID: loan.ID,
                Status: newLoan.Status,
                StaffName: "",
                StaffMail: "",
                CreatedAt: time.Now().UTC(),
                UpdatedAt: time.Now().UTC(),
            }

            if err = s.Repository.LoanHistory.Create(ctx, invested); err != nil {
                slog.Error("loan service: error on create invested history", slog.String("err", err.Error()))
                return coreError.Error(http.StatusInternalServerError, "internal server error")
            }

            investors, err := s.Repository.LoanInvestment.GetLoanInvestors(ctx, loan.ID)

            if err != nil {
                slog.Error("loan service: error on fetch loan investors", slog.String("err", err.Error()))
                return coreError.Error(http.StatusInternalServerError, "internal server error")
            }

            eg, ctx := errgroup.WithContext(ctx)
            eg.SetLimit(goroutineLimit)

            for _, investor := range investors {
                i := investor

                eg.TryGo(func() error {
                    lenderAgreementURL, err := s.generateLenderPDF(ctx, loan, i)

                    if err != nil {
                        return err
                    }

                    if err = s.sendLendingAgreement(ctx, loan, i, lenderAgreementURL); err != nil {
                        return err
                    }

                    return nil
                })
            }

            if err := eg.Wait(); err != nil {
                return err
            }
        }

        return nil
    })

    if err != nil {
        return err
    }

    return nil
}

func (s *Service) generateLenderPDF(
    ctx context.Context,
    loan db.Loan,
    investment db.LoanInvestment,
) (signedURL string, err error) {
    content := `
Surat Perjanjian

Berikut adalah surat perjanjian investasi pada peminjaman:
- Nama Peminjam: %s
- Nama Loan: %s
- Nilai Investasi Anda: %.2f
- Bunga Bagi Hasil: %.2f %%
`

    content = fmt.Sprintf(
        content,
        investment.Name,
        loan.Name,
        investment.Amount,
        loan.ROIRate,
    )

    pdfBuff, err := files.GeneratePDFBuffer(content)

    if err != nil {
        return "", err
    }

    agreementPathKey := filepath.Join(
        "lenders",
        "loans",
        loan.ID.String(),
        "investors",
        investment.InvestorID.String(),
        "agreement.pdf",
    )

    _, err = s.Bucket.Upload(ctx, agreementPathKey, pdfBuff)

    if err != nil {
        return "", err
    }

    if err = s.Repository.LoanInvestment.SetAgreementURL(ctx, agreementPathKey, investment.ID); err != nil {
        return "", err
    }

    signedURL, err = s.Bucket.GetSignURL(ctx, agreementPathKey)

    if err != nil {
        return "", err
    }

    return signedURL, nil
}

func (s *Service) sendLendingAgreement(
    ctx context.Context,
    loan db.Loan,
    investment db.LoanInvestment,
    agreementURL string,
) error {
    var sb strings.Builder

    sb.WriteString("<p>Surat Perjanjian Investasi</p>")
    sb.WriteString(fmt.Sprintf("<p>Nama Loan: %s</p><br/>", loan.Name))
    sb.WriteString(fmt.Sprintf("<p>Nama Investor: %s</p><br/>", investment.Name))
    sb.WriteString(fmt.Sprintf("<p>Jumlah Investasi: %.2f</p><br/>", loan.TotalInvestedAmount))
    sb.WriteString(fmt.Sprintf("<p>Bunga: %.2f</p><br/>", loan.ROIRate))
    sb.WriteString(fmt.Sprintf("<p>Dokumen: %s</p><br/>", agreementURL))

    return s.MailClient.SendEmail(ctx, mailer.SendEmail{
        Subject: "Surat Perjanjian Investasi",
        ToEmail: investment.Email,
        Body: sb.String(),
    })
}
