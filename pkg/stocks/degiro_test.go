package stocks

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseDegiroPortfolio(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "degiro", "Portfolio.csv")
	csvData, err := ioutil.ReadFile(csvFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", csvFile, err)
	}
	_, err = ParseDegiroPortfolio(csvData, "degiro")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
