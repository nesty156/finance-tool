package statement

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseCeskaSporitelnaStatement(t *testing.T) {
	jsonFile := filepath.Join("..", "..", "test-data", "ceska-sporitelna", "2018-05-01_2023-05-16.json")
	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", jsonFile, err)
	}
	_, err = ParseCeskaSporitelnaStatement(jsonData)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
