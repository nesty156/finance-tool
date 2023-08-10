package banks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateCSStatement(t *testing.T) {
	jsonFile := filepath.Join("..", "..", "test-data", "ceska-sporitelna", "2018-05-01_2023-05-16.json")
	jsonData, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", jsonFile, err)
	}
	_, err = CreateCSStatement(jsonData)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
