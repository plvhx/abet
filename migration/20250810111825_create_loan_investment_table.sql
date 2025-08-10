-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS loanInvestment (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    loanID UUID REFERENCES loans (id) NOT NULL,
    investorID UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    agreementUrl TEXT,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idxLoanInvestmentInvestorId ON loanInvestment (investorID);
CREATE INDEX IF NOT EXISTS idxLoanInvestmentLoanId ON loanInvestment (loanID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS loanInvestment;
-- +goose StatementEnd
