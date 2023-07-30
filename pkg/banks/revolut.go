package banks

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/gocarina/gocsv"
)

type RevolutTransaction struct {
	AccountingDate DateTime `csv:"Started Date"`
	ExecutionDate  DateTime `csv:"Completed Date"`
	Type           string   `csv:"Type"`
	Category       string   `csv:"Product"`
	Details        string   `csv:"Description"`
	Amount         float64  `csv:"Amount"`
	Fee            float64  `csv:"Fee"`
	Currency       string   `csv:"Currency"`
}

// Create statement of account from Revolut CSV file (transaction history)
func CreateRevolutStatement(filePath, accountName, currency string) (StatementOfAccount, error) {

	csvFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	revolutTXs := []*RevolutTransaction{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		return r
	})

	if err := gocsv.UnmarshalFile(csvFile, &revolutTXs); err != nil { // Load clients from file
		panic(err)
	}

	// Convert to internal format
	transactions := []Transaction{}
	for _, tx := range revolutTXs {
		transactions = append(transactions, RevolutTXConvert(*tx))
	}

	soa := StatementOfAccount{AccountNumber: accountName, Transactions: transactions, Currency: currency, StartDate: transactions[len(transactions)-1].AccountingDate, EndDate: transactions[0].AccountingDate}

	return soa, nil
}

func RevolutTXConvert(tx RevolutTransaction) Transaction {
	return Transaction{
		AccountingDate: tx.AccountingDate.Time,
		ExecutionDate:  tx.ExecutionDate.Time,
		Type:           tx.Type,
		Category:       tx.Category,
		Details:        tx.Details,
		Amount:         tx.Amount,
		Fee:            tx.Fee,
		Currency:       tx.Currency,
	}
}
