package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func saveToJson(soa StatementOfAccount) {
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

func loadFromJson(filename string) (*StatementOfAccount, error) {
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
