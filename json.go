package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func saveSoaJson(soa StatementOfAccount) {
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

func loadSoaJson(filename string) (*StatementOfAccount, error) {
	// read the JSON file into a byte slice
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// create a new statement of account object to hold the parsed data
	var soa StatementOfAccount

	// parse the JSON data into the statement of account object
	err = json.Unmarshal(jsonData, &soa)
	if err != nil {
		return nil, err
	}

	// return the statement of account object
	return &soa, nil
}

func savePortfolioJson(portfolio Portfolio) {
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

func loadPortfolioJson(filename string) (*Portfolio, error) {
	// read the JSON file into a byte slice
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// create a new portfolio object to hold the parsed data
	var portfolio Portfolio

	// parse the JSON data into the portfolio object
	err = json.Unmarshal(jsonData, &portfolio)
	if err != nil {
		return nil, err
	}

	// return the portfolio object
	return &portfolio, nil
}
