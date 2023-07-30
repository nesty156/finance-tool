package banks

import (
	"path/filepath"
	"testing"
)

func TestCreateMonetaStatement(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "moneta", "txs.csv")
	_, err := CreateMonetaStatement(csvFile, "moneta")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
