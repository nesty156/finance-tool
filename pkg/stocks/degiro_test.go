package stocks

import (
	"path/filepath"
	"testing"
)

func TestCreateDegiroPortfolio(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "degiro", "Portfolio.csv")
	_, err := CreateDegiroPortfolio(csvFile, "degiro")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
