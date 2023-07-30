package banks

import (
	"encoding/json"
	"math"
	"time"
)

// Create statement of account from Ceska Sporitelna JSON file (transaction history)
func CreateCSStatement(jsonData []byte) (StatementOfAccount, error) {
	var data []struct {
		Booking        string `json:"booking"`
		PartnerAccount struct {
			Number   string `json:"number"`
			BankCode string `json:"bankCode"`
		} `json:"partnerAccount"`
		Amount struct {
			Value     int    `json:"value"`
			Precision uint8  `json:"precision"`
			Currency  string `json:"currency"`
		} `json:"amount"`
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return StatementOfAccount{}, err
	}

	statement := StatementOfAccount{AccountNumber: "CeskaSporitelna", Currency: "CZK"}

	for _, record := range data {
		bookingTime, err := time.Parse("2006-01-02T15:04:05.000-0700", record.Booking)
		if err != nil {
			return StatementOfAccount{}, err
		}
		transaction := Transaction{
			AccountingDate:     bookingTime,
			ExecutionDate:      bookingTime,
			Type:               "",
			Code:               "",
			Name:               "",
			AccountOrDebitCard: record.PartnerAccount.Number + "/" + record.PartnerAccount.BankCode,
			Details:            "",
			Amount:             float64(record.Amount.Value) / math.Pow(10, float64(record.Amount.Precision)),
			Fee:                0.0,
		}

		statement.Transactions = append(statement.Transactions, transaction)
	}

	return statement, nil
}
