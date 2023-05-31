package stocks

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseTrading212History(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "trading212", "vypis.csv")
	csvData, err := ioutil.ReadFile(csvFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", csvFile, err)
	}
	txs, err := ParseTrading212History(csvData)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	_ = TransactionsToPortfolio(txs, "trading212")
}
