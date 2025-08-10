-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS loanHistory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    loanID UUID REFERENCES loans (id) NOT NULL,
    status loanStatus NOT NULL,
    staffName VARCHAR(255) NOT NULL,
    staffMail VARCHAR(255) NOT NULL,
    staffID UUID,
    extra JSONB NOT NULL DEFAULT '{}',
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS unxLoanHistories ON loanHistory (loanID, status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS loanHistory;
-- +goose StatementEnd
