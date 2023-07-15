package statement

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dslipak/pdf"
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
	Name               string
	AccountOrDebitCard string
	Details            string
	Amount             float64
	Fee                float64
}

/* From row of statement and creates transaction with all data */
func createTransaction(row pdf.Row) (Transaction, int, error) {
	offset := 0
	const numFields = 6
	layout := "02.01.2006"

	parseDate := func(date string) (time.Time, error) {
		return time.Parse(layout, strings.TrimSpace(date))
	}

	parseFloat := func(value string) (float64, error) {
		value = strings.ReplaceAll(value, " ", "")
		value = strings.ReplaceAll(value, ",", ".")
		return strconv.ParseFloat(value, 64)
	}

	/* Automatically adds offset */
	parseField := func(index int) string {
		return strings.TrimSpace(row.Content[index+offset].S)
	}

	var err error
	temp := Transaction{}
	temp.AccountingDate, err = parseDate(parseField(0))
	if err != nil {
		return Transaction{}, 0, err
	}

	temp.ExecutionDate, err = parseDate(parseField(2))
	if err != nil {
		return Transaction{}, 0, err
	}

	if row.Content[10].S == " " {
		temp.Type = parseField(6)
		temp.Code = parseField(8)
	} else {
		temp.Type = parseField(6) + " " + parseField(8)
		temp.Code = parseField(10)
		offset += 2
	}

	if row.Content[10+offset].S == " " && row.Content[12+offset].S == " " {
		offset -= 4
	} else if row.Content[10+offset].S == " " && row.Content[14+offset].S == " " {
		temp.AccountOrDebitCard = parseField(12)
		offset -= 2
	} else {
		temp.Name = parseField(12)
		temp.AccountOrDebitCard = parseField(14)
		if row.Content[16+offset].S != " " {
			temp.AccountOrDebitCard += "\n" + parseField(16)
			offset += 2
		}
	}

	if row.Content[18+offset].S != " " {
		temp.Details = parseField(18)
		offset += 2
		for row.Content[18+offset].S != " " {
			temp.Details += "\n" + parseField(18)
			offset += 2
		}
	}

	temp.Amount, err = parseFloat(parseField(20))
	if err != nil {
		return Transaction{}, 0, err
	}

	temp.Fee, err = parseFloat(parseField(24))
	if err != nil {
		return Transaction{}, 0, err
	}

	count := 0
	for i := range row.Content {
		if row.Content[i].S == " " {
			count++
			if count == numFields {
				return temp, i + 2, nil
			}
		}
	}

	return temp, 4, nil
}

func mergeTwoStatements(first, second StatementOfAccount) (StatementOfAccount, error) {
	if first.AccountNumber != second.AccountNumber {
		return StatementOfAccount{}, fmt.Errorf("Account numbers must be the same to merge")
	}
	if first.StartDate.After(second.StartDate) {
		tmp := first
		first = second
		second = tmp
	}
	first.EndDate = second.EndDate
	first.Transactions = append(first.Transactions, second.Transactions...)
	return first, nil
}

func MergeStatements(statements []StatementOfAccount) (StatementOfAccount, error) {
	if len(statements) == 0 {
		return StatementOfAccount{}, fmt.Errorf("Cannot merge an empty list of statements")
	}
	baseStatement := statements[0]
	for _, statement := range statements[1:] {
		tmp, err := mergeTwoStatements(baseStatement, statement)
		if err != nil {
			continue
		}
		baseStatement = tmp
	}
	return baseStatement, nil
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
