// Package banks implements functions for working with bank statements.
// Convert bank transactions file to Statement of Account format.
package banks

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

type StatementOfAccount struct {
	AccountNumber string
	StartDate     time.Time
	EndDate       time.Time
	Currency      string
	Transactions  []Transaction
}

type Transaction struct {
	AccountingDate     time.Time
	ExecutionDate      time.Time
	Type               string
	Code               string
	Category           string
	Name               string
	AccountOrDebitCard string
	Details            string
	Currency           string
	Amount             float64
	Fee                float64
}

type CZDateTime struct {
	time.Time
}

type USDateTime struct {
	time.Time
}

type Amount struct {
	float64
}

// Convert the CSV string as internal date
func (date *CZDateTime) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		return nil
	}
	date.Time, err = time.Parse("02.01.2006", csv)
	return err
}

// Convert the CSV string as internal date
func (date *USDateTime) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		return nil
	}
	date.Time, err = time.Parse("02/01/2006", csv)
	return err
}

// Convert the CSV string to internal float64
func (f *Amount) UnmarshalCSV(csv string) (err error) {
	csv = strings.ReplaceAll(csv, " ", "")
	csv = strings.ReplaceAll(csv, ",", ".")
	f.float64, err = strconv.ParseFloat(csv, 64)
	return err
}

func SortTransactions(statement StatementOfAccount) StatementOfAccount {
	// Sort transactions by execution date
	sort.Slice(statement.Transactions, func(i, j int) bool {
		return statement.Transactions[i].ExecutionDate.Before(statement.Transactions[j].ExecutionDate)
	})

	return statement
}

func SumTransactions(statement StatementOfAccount) float64 {
	total := 0.0

	// Loop through transactions and add up the amounts
	for _, transaction := range statement.Transactions {
		total += transaction.Amount - transaction.Fee
	}

	return total
}
