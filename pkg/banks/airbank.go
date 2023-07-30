package banks

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/charmap"
)

type AirBankTransaction struct {
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

func IsUTF8(content []byte) bool {
	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		return true
	}
	return false
}

func ConvertCP1250ToUTF8(filePath string) error {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
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
	err = ioutil.WriteFile(filePath, append([]byte{0xEF, 0xBB, 0xBF}, utf8Content...), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Create statement of account from AirBank CSV file (transaction history)
func CreateAirBankStatement(filePath string, accountName string) (StatementOfAccount, error) {

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

	soa := StatementOfAccount{AccountNumber: accountName, Transactions: transactions, Currency: "CZK", StartDate: transactions[len(transactions)-1].AccountingDate, EndDate: transactions[0].AccountingDate}

	return soa, nil
}

func AirBankTXConvert(mt AirBankTransaction) Transaction {
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
