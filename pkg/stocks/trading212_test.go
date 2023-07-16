package stocks

import (
	"path/filepath"
	"testing"
)

func TestParseTrading212History(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "trading212", "history.csv")
	txs, err := ParseTrading212History(csvFile)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	_ = TransactionsToPortfolio(txs, "trading212")
}
