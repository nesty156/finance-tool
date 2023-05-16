package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseCeskaSporitelnaStatement(t *testing.T) {
	jsonFile := filepath.Join("ceska-sporitelna", "2018-05-01_2023-05-16.json")
	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", jsonFile, err)
	}
	_, err = parseCeskaSporitelnaStatement(jsonData)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
