package statement

import (
	"path/filepath"
	"testing"
)

func TestParseMonetaStatement(t *testing.T) {
	csvFile := filepath.Join("..", "..", "test-data", "moneta", "moneta-acc.csv")
	_, err := ParseMonetaStatement(csvFile, "moneta")
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
