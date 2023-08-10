package banks

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/charmap"
)

type AirBankTransaction struct {
	ExecutionDate USDateTime `csv:"Datum provedení"`
	Type          string     `csv:"Typ úhrady"`
	Name          string     `csv:"Název protistrany"`
	Category      string     `csv:"Kategorie plateb"`
	AccountNumber string     `csv:"Číslo účtu protistrany"`
	Details       string     `csv:"Zpráva pro příjemce"`
	Amount        Amount     `csv:"Částka v měně účtu"`
	Fee           float64    `csv:"Poplatek v měně účtu"`
	Currency      string     `csv:"Měna účtu"`
}

func IsUTF8(content []byte) bool {
	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		return true
	}
	return false
}

func ConvertCP1250ToUTF8(filePath string) error {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Check if the file is already UTF-8 encoded
	if IsUTF8(content) {
		return nil // No need to convert, already UTF-8
	}

	// Convert CP1250 to UTF-8
	utf8Content, err := charmap.Windows1250.NewDecoder().Bytes(content)
	if err != nil {
		return err
	}

	// Write the UTF-8 content back to the file
	err = os.WriteFile(filePath, append([]byte{0xEF, 0xBB, 0xBF}, utf8Content...), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Create statement of account from AirBank CSV file (transaction history)
func CreateAirBankStatement(filePath, accountName, currency string) (StatementOfAccount, error) {

	// Convert CP1250 to UTF-8
	err := ConvertCP1250ToUTF8(filePath)
	if err != nil {
		panic(err)
	}

	csvFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	airbankTXs := []*AirBankTransaction{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma = ';' // Use semicolon separator
		return r      // Allows use quotes in CSV
	})

	if err := gocsv.UnmarshalFile(csvFile, &airbankTXs); err != nil { // Load clients from file
		panic(err)
	}

	// Convert to internal format
	transactions := []Transaction{}
	for _, tx := range airbankTXs {
		transactions = append(transactions, AirBankTXConvert(*tx))
	}

	soa := StatementOfAccount{AccountNumber: accountName, Transactions: transactions, Currency: currency, StartDate: transactions[len(transactions)-1].AccountingDate, EndDate: transactions[0].AccountingDate}

	return soa, nil
}

func AirBankTXConvert(tx AirBankTransaction) Transaction {
	return Transaction{
		AccountingDate:     tx.ExecutionDate.Time,
		ExecutionDate:      tx.ExecutionDate.Time,
		Type:               tx.Type,
		Name:               tx.Name,
		Category:           tx.Category,
		AccountOrDebitCard: tx.AccountNumber,
		Details:            tx.Details,
		Amount:             float64(tx.Amount.float64),
		Fee:                tx.Fee,
		Currency:           tx.Currency,
	}
}
