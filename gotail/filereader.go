package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type FileReader struct {
	filePath string
}

type FileIterator struct {
	file   *os.File
	reader *bufio.Reader
	offset int64
	buffer []byte
}

func NewFileIterator(file *os.File, bufferSize int, offset int64) (*FileIterator, error) {
	file.Seek(offset, io.SeekEnd)
	return &FileIterator{
		file:   file,
		offset: offset,
		buffer: make([]byte, bufferSize),
		reader: bufio.NewReader(file),
	}, nil
}

func (it *FileIterator) Next() (string, error) {
	bytes, err := it.reader.Read(it.buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	if bytes == 0 {
		return "", io.EOF
	}
	return string(it.buffer[:bytes]), nil
}

const bufferSize = 4096

func NewFileReader(filePath string) (*FileReader, error) {
	return &FileReader{filePath: filePath}, nil
}

func (fr *FileReader) Tail(n int) (*FileIterator, error) {
	file, err := os.Open(fr.filePath)
	if err != nil {
		return nil, err
	}
	offset, err := defineStartingOffset(file, n)
	if err != nil {
		return nil, err
	}
	return NewFileIterator(file, bufferSize, offset)
}

func defineStartingOffset(file *os.File, n int) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	fileSize := int64(info.Size())
	if fileSize == 0 {
		return 0, nil
	}
	totalLines := 0
	offset := -min(bufferSize, fileSize)
	resultOffset := int64(0)
	for totalLines <= n {
		file.Seek(offset, io.SeekEnd)
		resultOffset = offset
		buf := make([]byte, bufferSize)
		bytes, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return -fileSize, err
		}
		if bytes == 0 {
			return -fileSize, nil
		}
		lines := strings.Split(string(buf[:bytes]), "\n")
		totalLines = totalLines + len(lines)
		if totalLines > n {
			linesToDrop := totalLines - n
			for i := 0; i < linesToDrop; i++ {
				resultOffset = resultOffset + int64(len(lines[i])) + 1
			}
			return resultOffset, nil
		}
		offset = offset + int64(bufferSize)
		if offset > fileSize {
			return -fileSize, nil
		}
	}
	return -fileSize, nil
}

func (fr *FileIterator) Close() error {
	return fr.file.Close()
}
