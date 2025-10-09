package main

import (
	"os"
	"testing"
)

func TestFileReaderReturnsLastLineWhenNEqualsToOne(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "line1\nline2\nline3"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fr, err := NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create FileReader: %v", err)
	}
	defer fr.Close()

	lines, err := fr.Tail(1)
	if err != nil {
		t.Fatalf("Tail failed: %v", err)
	}

	expected := []string{"line3"}
	if len(lines) != len(expected) || lines[0] != expected[0] {
		t.Errorf("Expected %v, got %v", expected, lines)
	}
}
