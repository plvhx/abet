package types

type LoanStatus string

func (l LoanStatus) String() string {
    return string(l)
}

func (l LoanStatus) Next() LoanStatus {
    stateMap := map[string]LoanStatus{
        StatusProposed.String(): StatusApproved,
        StatusApproved.String(): StatusInvested,
        StatusInvested.String(): StatusDisbursed,
    }

    return stateMap[string(l)]
}

const (
    StatusProposed LoanStatus = "proposed"
    StatusApproved LoanStatus = "approved"
    StatusInvested LoanStatus = "invested"
    StatusDisbursed LoanStatus = "disbursed"
)
