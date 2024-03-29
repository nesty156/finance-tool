package bitcoin

import (
	"testing"
	"time"
)

func TestConvertToStatementOfAccount(t *testing.T) {
	// Test case 1
	btcAcc := BtcAccount{
		AccountNumber: "1234",
		Currency:      "BTC",
		Transactions: []BtcTransaction{
			{
				BlockTime: 1672003002,
				Type:      "send",
				Code:      "abcd1234",
				Amount:    "0.01",
				Fee:       "0.0001",
				Vsize:     200,
				FeeRate:   "0.5",
				Details:   "Transaction details",
			},
			{
				BlockTime: 1671934862,
				Type:      "receive",
				Code:      "efgh5678",
				Amount:    "0.1",
				Fee:       "0.0002",
				Vsize:     300,
				FeeRate:   "0.75",
				Details:   "Transaction details",
			},
		},
	}

	soa, err := btcAcc.ConvertToStatementOfAccount()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if soa.AccountNumber != "1234" {
		t.Errorf("Expected account number to be '1234', but got '%s'", soa.AccountNumber)
	}
	if soa.Currency != "BTC" {
		t.Errorf("Expected currency to be 'BTC', but got '%s'", soa.Currency)
	}
	if len(soa.Transactions) != 2 {
		t.Errorf("Expected 2 transactions, but got %d", len(soa.Transactions))
	}
	start := time.Unix(1671934862, 0)
	if !soa.StartDate.Equal(start) {
		t.Errorf("Expected start date to be %s, but got '%s'", start, soa.StartDate)
	}
	end := time.Unix(1672003002, 0)
	if !soa.EndDate.Equal(end) {
		t.Errorf("Expected end date to be %s, but got '%s'", end, soa.EndDate)
	}
	if soa.Transactions[0].Amount != 0.01 {
		t.Errorf("Expected amount of first transaction to be 0.01, but got %f", soa.Transactions[0].Amount)
	}
	if soa.Transactions[1].Amount != 0.1 {
		t.Errorf("Expected amount of second transaction to be 0.1, but got %f", soa.Transactions[1].Amount)
	}
}
