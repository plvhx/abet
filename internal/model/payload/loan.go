package payload

import (
    "mime/multipart"

    "github.com/google/uuid"
    "abet/pkg/types"
)

type CreateLoan struct {
    Name            string    `json:"name" validate:"required,max=255"`
    BorrowerID      uuid.UUID `json:"borrowerId" validate:"required,max=100"`
    BorrowerName    string    `json:"borrowerName" validate:"required,max=255"`
    BorrowerEmail   string    `json:"borrowerEmail" validate:"required,email"`
    BorrowerRate    float64   `json:"borrowerRate" validate:"required,max=100"`
    PrincipalAmount float64   `json:"principalAmount" validate:"required"`
    ROIRate         float64   `json:"roiRate" validate:"required"`
    StaffID         uuid.UUID `json:"staffID" validate:"required"`
    StaffName       string    `json:"staffName" validate:"required"`
    StaffEmail      string    `json:"staffEmail" validate:"required,email"`
}

type ApproveLoan struct {
    LoanID     uuid.UUID `json:"loanID" validate:"required"`
    StaffID    uuid.UUID `json:"staffID" validate:"required"`
    StaffEmail string    `json:"staffEmail" validate:"required,email"`
    StaffName  string    `json:"staffName" validate:"required,max=255"`
}

type DisburseLoan struct {
    LoanID                  uuid.UUID       `json:"loanID" validate:"required"`
    FieldVisitStaffID       uuid.UUID       `form:"fieldVisitStaffID" validate:"required"`
    FieldVisitStaffName     string          `form:"fieldVisitStaffName" validate:"required,max=255"`
    FieldVisitStaffMail     string          `form:"fieldVisitStaffMail" validate:"required,email"`
    StaffID                 uuid.UUID       `form:"staffID" validate:"required"`
    StaffName               string          `form:"staffName" validate:"required,max=255"`
    StaffMail               string          `form:"staffMail" validate:"required,email"`
    SignedAgreementDocument *multipart.File `form:"signedAgreementDocument"`
    DocumentExtension       string          `form:"documentExtension"`
}

type InvestLoan struct {
    LoanID           uuid.UUID `json:"loanID" validate:"required"`
    InvestorID       uuid.UUID `json:"investorID" validate:"required"`
    InvestorName     string    `json:"investorName" validate:"required,max=255"`
    InvestorEmail    string    `json:"investorEmail" validate:"required,email"`
    InvestmentAmount float64   `json:"investmentAmount" validate:"required"`
}

type LoanPaginationFilter struct {
    PaginationFilter
    Status types.LoanStatus `json:"status" query:"status"`
}
