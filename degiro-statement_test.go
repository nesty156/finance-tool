package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseDegiroPortfolio(t *testing.T) {
	csvFile := filepath.Join("degiro", "Portfolio.csv")
	csvData, err := ioutil.ReadFile(csvFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", csvFile, err)
	}
	_, err = parseDegiroPortfolio(csvData, "degiro")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
