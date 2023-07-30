package stocks

import (
	"path/filepath"
	"testing"
)

func TestCreateTrading212Portfolio(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "trading212", "txs.csv")
	_, err := CreateTrading212Portfolio(csvFile, "trading212", "EUR")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
