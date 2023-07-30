package banks

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/gocarina/gocsv"
)

type MonetaTransaction struct {
	AccountingDate DateTime `csv:"Splatnost"`
	ExecutionDate  DateTime `csv:"Odesláno"`
	Type           string   `csv:"Typ transakce"`
	Code           string   `csv:"-"`
	Name           string   `csv:"Název účtu příjemce"`
	AccountNumber  string   `csv:"Číslo protiúčtu"`
	AccountBank    string   `csv:"Banka protiúčtu"`
	Details        string   `csv:"Zpráva pro příjemce"`
	Amount         Amount   `csv:"Částka"`
	Fee            float64  `csv:"-"`
}

// Create statement of account from Moneta CSV file (transaction history)
func CreateMonetaStatement(fileName string, accountName string) (StatementOfAccount, error) {

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	monetaTransactions := []*MonetaTransaction{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma = ';' // Use semicolon separator
		return r      // Allows use quotes in CSV
	})

	if err := gocsv.UnmarshalFile(csvFile, &monetaTransactions); err != nil { // Load clients from file
		panic(err)
	}

	// Convert to internal format
	transactions := []Transaction{}
	for _, mt := range monetaTransactions {
		transactions = append(transactions, MonetaTXConvert(*mt))
	}

	soa := StatementOfAccount{AccountNumber: accountName, Transactions: transactions, Currency: "CZK", StartDate: transactions[len(transactions)-1].AccountingDate, EndDate: transactions[0].AccountingDate}

	return soa, nil
}

func MonetaTXConvert(mt MonetaTransaction) Transaction {
	return Transaction{
		AccountingDate:     mt.AccountingDate.Time,
		ExecutionDate:      mt.ExecutionDate.Time,
		Type:               mt.Type,
		Code:               mt.Code,
		Name:               mt.Name,
		AccountOrDebitCard: mt.AccountNumber + "/" + mt.AccountBank,
		Details:            mt.Details,
		Amount:             float64(mt.Amount.float64),
		Fee:                mt.Fee,
	}
}
