package banks

import (
	"path/filepath"
	"testing"
)

/* TODO create fake data */
func TestParseAirBankStatement(t *testing.T) {
	pdfFile := filepath.Join("..", "..", "test-data", "airbank-czk", "vypis02-2020.pdf")
	_, err := ParseAirBankStatement(pdfFile)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
