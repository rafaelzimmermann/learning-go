package main

import (
	"io"
	"os"
)

type FileReader struct {
	filePath string
}

type FileIterator struct {
	file      *os.File
	fileSize  int64
	offset    int64
	buffer    []byte
	firstRead bool
}

func NewFileIterator(file *os.File, fileSize int64, buffer *[]byte, offset int64) (*FileIterator, error) {
	file.Seek(offset, io.SeekStart)
	return &FileIterator{
		file:      file,
		fileSize:  fileSize,
		offset:    offset,
		buffer:    *buffer,
		firstRead: true,
	}, nil
}

var emptyByteArray = make([]byte, 0)

func (it *FileIterator) Next() ([]byte, error) {
	if it.offset >= it.fileSize {
		return emptyByteArray, io.EOF
	}
	bytes, err := it.file.ReadAt(it.buffer, it.offset)
	it.offset += int64(bytes)
	if err != nil && err != io.EOF {
		return emptyByteArray, err
	}
	if bytes == 0 {
		return emptyByteArray, io.EOF
	}
	return it.buffer[:bytes], nil
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
	buf := make([]byte, bufferSize)
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	offset, err := defineStartingOffset(file, info.Size(), &buf, n)
	if err != nil {
		return nil, err
	}
	return NewFileIterator(file, info.Size(), &buf, offset)
}

func defineStartingOffset(file *os.File, fileSize int64, buf *[]byte, n int) (int64, error) {
	if fileSize == 0 {
		return 0, nil
	}
	totalLines := 0
	resultOffset := int64(fileSize)
	isFileLastChar := true
	for totalLines <= n {
		bytes, err := file.ReadAt(*buf, max(0, resultOffset-bufferSize))
		if err != nil && err != io.EOF {
			return resultOffset, err
		}
		if bytes == 0 {
			return resultOffset, nil
		}
		for i := bytes - 1; i >= 0; i-- {
			if (*buf)[i] == '\n' {
				if !isFileLastChar {
					totalLines++
				}
			}
			if totalLines == n {
				return resultOffset, nil
			}
			resultOffset--
			isFileLastChar = false
		}
		if resultOffset <= 0 {
			return 0, nil
		}
	}
	return 0, nil
}

func (fr *FileIterator) Close() error {
	return fr.file.Close()
}
