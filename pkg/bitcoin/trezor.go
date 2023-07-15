package bitcoin

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"

	stat "github.com/nesty156/finance-tool/pkg/statement"
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

func ParseBtcAccount(filepath string) (*BtcAccount, error) {
	// read the JSON file into a byte slice
	jsonData, err := ioutil.ReadFile(filepath)
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

func (btcAcc BtcAccount) ConvertToStatementOfAccount() (*stat.StatementOfAccount, error) {
	soa := stat.StatementOfAccount{
		AccountNumber: btcAcc.AccountNumber,
		Currency:      btcAcc.Currency,
	}

	transactions := make([]stat.Transaction, len(btcAcc.Transactions))
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
