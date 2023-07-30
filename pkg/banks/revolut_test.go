package banks

import (
	"path/filepath"
	"testing"
)

func TestCreateRevolutStatement(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "revolut", "txs.csv")
	_, err := CreateRevolutStatement(csvFile, "revolut", "CZK")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
