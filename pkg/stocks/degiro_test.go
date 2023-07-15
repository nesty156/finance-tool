package stocks

import (
	"path/filepath"
	"testing"
)

func TestParseDegiroPortfolio(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "degiro", "Portfolio.csv")
	_, err := ParseDegiroPortfolio(csvFile, "degiro")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
