package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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

	iterator, err := fr.Tail(1)
	if err != nil {
		t.Fatalf("Tail failed: %v", err)
	}
	defer iterator.Close()

	expected := []string{"line3"}
	bufferContent, err := iterator.Next()
	if err != nil && len(bufferContent) == 0 {
		t.Fatalf("Iterator Next failed: %v", err)
	}
	result := strings.Split(bufferContent, "\n")
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFileReaderIgnoresLastEmptyLine(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "line1\nline2\nline3\n"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fr, err := NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create FileReader: %v", err)
	}

	iterator, err := fr.Tail(1)
	if err != nil {
		t.Fatalf("Tail failed: %v", err)
	}
	defer iterator.Close()

	expected := string("line3\n")
	bufferContent, err := iterator.Next()
	if err != nil {
		t.Fatalf("Iterator Next failed: %v", err)
	}
	result := string(bufferContent)
	if len(result) != len(expected) || result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFileReaderReturnsAllLinesWhenNGreaterThanTotalLines(t *testing.T) {
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

	iterator, err := fr.Tail(10)
	if err != nil {
		t.Fatalf("Tail failed: %v", err)
	}
	defer iterator.Close()

	expected := content
	result, err := iterator.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("Iterator Next failed: %v", err)
	}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFileReaderReturnsLastNLinesWhenFileBiggerThanBuffer(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := strings.Repeat("file_content\n", 400) + "\nline1\nline2\nline3"
	if len(content) < 4096 {
		t.Fatalf("Test content is not larger than buffer size")
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fr, err := NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create FileReader: %v", err)
	}

	iterator, err := fr.Tail(1)
	if err != nil {
		t.Fatalf("Tail failed: %v", err)
	}
	defer iterator.Close()

	expected := []string{"line3"}
	bufferContent, err := iterator.Next()
	if err != nil && len(bufferContent) == 0 {
		t.Fatalf("Iterator Next failed: %v", err)
	}
	result := strings.Split(bufferContent, "\n")
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFileReaderReturnsLast1000LinesFrom4000LinesFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := ""
	for i := range 4000 {
		content += fmt.Sprintf("line%v\n", i+1)
	}

	if len(content) < 4096 {
		t.Fatalf("Test content is not larger than buffer size")
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fr, err := NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create FileReader: %v", err)
	}

	iterator, err := fr.Tail(1000)
	if err != nil {
		t.Fatalf("Tail failed: %v", err)
	}
	defer iterator.Close()

	bufferContent, err := iterator.Next()
	output := string(bufferContent)
	for err == nil && err != io.EOF && len(bufferContent) > 0 {

		bufferContent, err = iterator.Next()
		output += string(bufferContent)
	}
	result := strings.Split(strings.Trim(output, "\n"), "\n")
	if len(result) != 1000 {
		t.Errorf("Expected %v, got %v", 1000, len(result))
	}
}
