package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	stat "github.com/nesty156/finance-tool/pkg/statement"
	"github.com/nesty156/finance-tool/pkg/stocks"
)

func SaveSoaJson(soa stat.StatementOfAccount) {
	// convert the statement of account object to a JSON byte slice
	jsonData, err := json.MarshalIndent(soa, "", "    ")
	if err != nil {
		panic(err)
	}

	name := strings.ReplaceAll(soa.AccountNumber+".json", "/", "-")

	// write the JSON byte slice to a file
	err = ioutil.WriteFile(name, jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Statement of account saved to " + name)
}

func LoadSoaJson(filename string) (*stat.StatementOfAccount, error) {
	// read the JSON file into a byte slice
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// create a new statement of account object to hold the parsed data
	var soa stat.StatementOfAccount

	// parse the JSON data into the statement of account object
	err = json.Unmarshal(jsonData, &soa)
	if err != nil {
		return nil, err
	}

	// return the statement of account object
	return &soa, nil
}

func SavePortfolioJson(portfolio stocks.Portfolio) {
	// convert the statement of account object to a JSON byte slice
	jsonData, err := json.MarshalIndent(portfolio, "", "    ")
	if err != nil {
		panic(err)
	}

	name := strings.ReplaceAll(portfolio.Name+".json", "/", "-")

	// write the JSON byte slice to a file
	err = ioutil.WriteFile(name, jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Statement of account saved to " + name)
}

func LoadPortfolioJson(filename string) (*stocks.Portfolio, error) {
	// read the JSON file into a byte slice
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// create a new portfolio object to hold the parsed data
	var portfolio stocks.Portfolio

	// parse the JSON data into the portfolio object
	err = json.Unmarshal(jsonData, &portfolio)
	if err != nil {
		return nil, err
	}

	// return the portfolio object
	return &portfolio, nil
}
