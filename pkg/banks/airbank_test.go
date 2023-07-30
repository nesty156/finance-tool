package banks

import (
	"path/filepath"
	"testing"
)

func TestCreateAirBankStatement(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "airbank", "txs.csv")
	_, err := CreateAirBankStatement(csvFile, "airbank", "CZK")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
