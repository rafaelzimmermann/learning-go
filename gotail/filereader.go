package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type FileReader struct {
	file   *os.File
	reader *bufio.Reader
}

const bufferSize = 4096

func NewFileReader(filePath string) (*FileReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &FileReader{file: file, reader: bufio.NewReader(file)}, nil
}

func (fr *FileReader) Tail(n int) ([]string, error) {
	info, err := fr.file.Stat()
	if err != nil {
		return nil, err
	}
	offset := -min(bufferSize, int64(info.Size()))
	fr.file.Seek(offset, io.SeekEnd)

	buf := make([]byte, bufferSize)
	bytes, err := fr.reader.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	lines := strings.Split(string(buf[:bytes]), "\n")
	start := 0
	totalLines := len(lines)
	if totalLines > n {
		start = totalLines - n
	}
	return lines[start:], nil
}

func (fr *FileReader) Close() error {
	return fr.file.Close()
}

// At this point, this reads once the file using a buffer of 4096 bytes
// and returns the last n lines. So next step would be to keep navigating
// up in the file until we have enough lines.
// We should only count the lines, till find the correct position, to avoid
// running out of memory for large files and to be as efficient as possible.

// As the end goal is to read multple files at once, I think a printer will be
// provided. This printer will be shared between multiple filereaders. Then I
// will get into how to read it concurrently.

// probably the printer will have some form of queue, which readers will try to
// write to. The printer will be the only one writing to stdout.
