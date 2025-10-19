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
	if it.firstRead {
		partBufferSize := it.offset % int64(len(it.buffer))
		it.firstRead = false
		if it.offset < bufferSize {
			it.offset += partBufferSize
			return it.buffer[it.offset-partBufferSize : it.fileSize], nil
		}
		if (it.fileSize - it.offset) < bufferSize {
			partBufferSize = it.fileSize - it.offset
			it.offset += partBufferSize
			return it.buffer[bufferSize-partBufferSize : bufferSize], nil
		}
		it.offset += partBufferSize
		return it.buffer[bufferSize-partBufferSize : bufferSize], nil
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
	offset := max(0, fileSize-bufferSize)
	resultOffset := int64(offset)
	isTail := true
	for totalLines <= n {
		file.Seek(resultOffset, io.SeekStart)
		bytes, err := file.Read(*buf)
		if err != nil && err != io.EOF {
			return offset, err
		}
		if bytes == 0 {
			return offset, nil
		}
		if isTail && (*buf)[bytes-1] == '\n' {
			bytes = bytes - 1
			isTail = false
		}
		for i := bytes - 1; i >= 0; i-- {
			if (*buf)[i] == '\n' {
				totalLines++
			}
			if totalLines == n {
				return resultOffset + int64(i) + 1, nil
			}
		}
		resultOffset = resultOffset - int64(bufferSize) - 1
		if resultOffset <= 0 {
			return 0, nil
		}
	}
	return 0, nil
}

func (fr *FileIterator) Close() error {
	return fr.file.Close()
}
