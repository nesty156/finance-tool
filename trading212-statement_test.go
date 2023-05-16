package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseTrading212History(t *testing.T) {
	csvFile := filepath.Join("trading212", "from_2023-02-03_to_2023-05-16.csv")
	csvData, err := ioutil.ReadFile(csvFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", csvFile, err)
	}
	_, err = parseTrading212History(csvData)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
