package bitcoin

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nesty156/finance-tool/pkg/banks"
	"github.com/nesty156/finance-tool/pkg/converter"
	"github.com/nesty156/finance-tool/pkg/util"
)

type BtcAccount struct {
	AccountNumber string           `json:"-"`
	Currency      string           `json:"coin"`
	Transactions  []BtcTransaction `json:"transactions"`
}

type BtcTransaction struct {
	BlockTime int         `json:"blockTime"`
	Type      string      `json:"type"`
	Code      string      `json:"txid"`
	Amount    string      `json:"amount"`
	Fee       string      `json:"fee"`
	Vsize     int         `json:"vsize"`
	FeeRate   string      `json:"feeRate"`
	Details   string      `json:"-"`
	Targets   []BtcTarget `json:"targets"`
}

type BtcTarget struct {
	IsAddress       bool   `json:"isAddress"`
	Amount          string `json:"amount"`
	IsAccountTarget bool   `json:"isAccountTarget"`
	Details         string `json:"metadataLabel"`
}

type TrezorStat struct {
	Name      string
	Component string
	Currency  string
	Value     float64
}

// Convert Trezor files to Statement of Account format
func ConvertTrezorToStatement(dirPath string) ([]TrezorStat, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	btcCZK, err := converter.GetBitcoinPrice("CZK")
	if err != nil {
		return nil, fmt.Errorf("error getting bitcoin price: %v", err)
	}

	stats := []TrezorStat{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		account, err := ParseBtcAccount(filePath)
		if err != nil {
			log.Printf("Error loading file %s: %v", filePath, err)
			continue
		}

		statement, err := account.ConvertToStatementOfAccount()
		if err != nil {
			log.Printf("Error converting bitcoin account to statement %s: %v", filePath, err)
			continue
		}

		value := banks.SumTransactions(*statement)

		fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value*btcCZK, "CZK")
		util.SaveSoaJson(*statement)
		stats = append(stats, TrezorStat{Name: statement.AccountNumber, Component: "trezor", Currency: "BTC", Value: value})
	}

	return stats, nil
}

func ParseBtcAccount(filepath string) (*BtcAccount, error) {
	// read the JSON file into a byte slice
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// create a new btc account object to hold the parsed data
	var btcAcc BtcAccount

	// parse the JSON data into the statement of account object
	err = json.Unmarshal(jsonData, &btcAcc)
	if err != nil {
		return nil, err
	}

	for i, transaction := range btcAcc.Transactions {
		for _, target := range transaction.Targets {
			if target.IsAccountTarget {
				btcAcc.Transactions[i].Details = target.Details
				break
			}
		}
	}

	filename := path.Base(strings.ReplaceAll(filepath, `\`, `/`))
	btcAcc.AccountNumber = strings.TrimSuffix(filename, ".json")

	return &btcAcc, nil
}

func (btcAcc BtcAccount) ConvertToStatementOfAccount() (*banks.StatementOfAccount, error) {
	soa := banks.StatementOfAccount{
		AccountNumber: btcAcc.AccountNumber,
		Currency:      btcAcc.Currency,
	}

	transactions := make([]banks.Transaction, len(btcAcc.Transactions))
	cz, _ := time.LoadLocation("Europe/Prague")
	start := time.Unix(9999999999, 0)
	end := time.Unix(0, 0)

	for i, transaction := range btcAcc.Transactions {
		timestamp := int64(transaction.BlockTime)
		date := time.Unix(timestamp, 0)
		if date.Before(start) {
			start = date
		}
		if date.After(end) {
			end = date
		}
		dateInCZ := date.In(cz)
		transactions[i].AccountingDate = dateInCZ
		transactions[i].ExecutionDate = dateInCZ

		transactions[i].Type = transaction.Type
		transactions[i].Code = transaction.Code
		transactions[i].Details = transaction.Details

		var err error

		transactions[i].Amount, err = strconv.ParseFloat(transaction.Amount, 64)
		if err != nil {
			return nil, err
		}
		if transaction.Type == "sent" {
			transactions[i].Amount = -transactions[i].Amount
		}
		transactions[i].Fee, err = strconv.ParseFloat(transaction.Fee, 64)
		if err != nil {
			return nil, err
		}
	}

	soa.StartDate = start.In(cz)
	soa.EndDate = end.In(cz)
	soa.Transactions = transactions

	return &soa, nil
}
