package repository

import (
    "abet/pkg"
)

type Repository struct {
    Loan *LoanRepository
    LoanHistory *LoanHistoryRepository
    LoanInvestment *LoanInvestmentRepository
}

func NewRepository(options *pkg.Options) *Repository {
    return &Repository{
        Loan: NewLoanRepository(options),
        LoanHistory: NewLoanHistoryRepository(options),
        LoanInvestment: NewLoanInvestmentRepository(options),
    }
}
