package statement

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseMonetaStatement(t *testing.T) {
	xmlFile := filepath.Join("..", "..", "test-data", "moneta", "fake-statement04-2019.xml")
	xmlData, err := ioutil.ReadFile(xmlFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", xmlFile, err)
	}
	_, err = ParseMonetaStatement(xmlData)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
