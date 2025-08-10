package db

import (
    "time"

    "github.com/google/uuid"
    "github.com/guregu/null"
    "abet/pkg/types"
)

type Loan struct {
    ID                   uuid.UUID        `db:"id" json:"id,omitempty"`
    Name                 string           `db:"name" json:"name,omitempty"`
    Status               types.LoanStatus `db:"status" json:"status,omitempty"`
    BorrowerID           uuid.UUID        `db:"borrowerid" json:"borrowerId,omitempty"`
    BorrowerName         string           `db:"borrowername" json:"borrowerName,omitempty"`
    BorrowerEmail        string           `db:"borroweremail" json:"borrowerEmail,omitempty"`
    BorrowerRate         float64          `db:"borrowerrate" json:"borrowerRate,omitempty"`
    BorrowerAgreementURL null.String      `db:"borroweragreementurl" json:"borrowerAgreementUrl,omitempty"`
    PrincipalAmount      float64          `db:"principalamount" json:"principalAmount,omitempty"`
    TotalInvestedAmount  float64          `db:"totalinvestedamount" json:"totalInvestedAmount,omitempty"`
    ROIRate              float64          `db:"roirate" json:"roiRate,omitempty"`
    StaffName            string           `db:"staffname" json:"staffName,omitempty"`
    StaffEmail           string           `db:"staffemail" json:"staffEmail,omitempty"`
    StaffID              uuid.UUID        `db:"staffid" json:"staffId,omitempty"`
    StaffVisitProof      string           `db:"staffvisitproof" json:"staffVisitProof,omitempty"`
    ApprovedAt           null.Time        `db:"approvedat" json:"approvedAt"`
    DisbursedAt          null.Time        `db:"disbursedat" json:"disbursedAt"`
    CreatedAt            time.Time        `db:"createdat" json:"createdAt"`
    UpdatedAt            time.Time        `db:"updatedat" json:"updatedAt"`
}

type LoanHistory struct {
    ID        uuid.UUID         `db:"id" json:"id,omitempty"`
    LoanID    uuid.UUID         `db:"loanid" json:"loanID,omitempty"`
    Status    types.LoanStatus  `db:"status" json:"status,omitempty"`
    StaffName string            `db:"staffname" json:"staffName,omitempty"`
    StaffMail string            `db:"staffmail" json:"staffMail,omitempty"`
    StaffID   uuid.UUID         `db:"staffid" json:"staffID,omitempty"`
    Extra     map[string]string `db:"extra" json:"extra,omitempty"`
    CreatedAt time.Time         `db:"createdat" json:"createdAt,omitempty"`
    UpdatedAt time.Time         `db:"updatedat" json:"updatedAt,omitempty"`
}

type LoanInvestment struct {
    ID           uuid.UUID `db:"id" json:"id,omitempty"`
    LoanID       uuid.UUID `db:"loanid" json:"loanID,omitempty"`
    InvestorID   uuid.UUID `db:"investorid" json:"investorID,omitempty"`
    Name         string    `db:"name" json:"name,omitempty"`
    Email        string    `db:"email" json:"email,omitempty"`
    Amount       float64   `db:"amount" json:"amount,omitempty"`
    AgreementURL string    `db:"agreementurl" json:"amount,omitempty"`
    CreatedAt    time.Time `db:"createdat" json:"createdAt,omitempty"`
    UpdatedAt    time.Time `db:"updatedat" json:"updatedAt,omitempty"`
}
