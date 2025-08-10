-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE loanStatus AS ENUM ('proposed', 'approved', 'invested', 'disbursed');

CREATE TABLE IF NOT EXISTS loans
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    status loanStatus NOT NULL DEFAULT 'proposed',
    borrowerId UUID NOT NULL,
    borrowerName VARCHAR(255) NOT NULL,
    borrowerEmail VARCHAR(255) NOT NULL,
    borrowerRate NUMERIC(5, 2) NOT NULL,
    borrowerAgreementUrl TEXT,
    principalAmount NUMERIC(12, 2) NOT NULL,
    totalInvestedAmount NUMERIC(12, 2) NOT NULL DEFAULT 0,
    roiRate NUMERIC(5, 2) NOT NULL,
    staffName VARCHAR(255) NOT NULL,
    staffEmail VARCHAR(255) NOT NULL,
    staffId UUID,
    staffVisitProof TEXT NOT NULL DEFAULT '',
    approvedAt TIMESTAMP,
    disbursedAt TIMESTAMP,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idxBorrowerId ON loans USING HASH (borrowerId);
CREATE INDEX IF NOT EXISTS idxLoanStatus ON loans USING HASH (status);
CREATE INDEX IF NOT EXISTS idxLoanName ON loans (name);
CREATE INDEX IF NOT EXISTS idxLoanOrderCreated ON loans (createdAt);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS loans;
DROP TYPE IF EXISTS loanStatus;
-- +goose StatementEnd
