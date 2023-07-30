package stocks

import (
	"path/filepath"
	"testing"
)

func TestCreateDegiroPortfolio(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "degiro", "txs.csv")
	_, err := CreateDegiroPortfolio(csvFile, "degiro", "EUR")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
