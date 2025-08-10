package types

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestCanGetLoanStatusAsString(t *testing.T) {
    assert.Equal(t, "proposed", StatusProposed.String())
}

func TestCanGetNextLoanStatusAsString(t *testing.T) {
    assert.Equal(t, "approved", StatusProposed.Next().String())
}
