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
	file.Seek(offset, io.SeekStart)
	return &FileIterator{
		file:   file,
		offset: offset,
		buffer: make([]byte, bufferSize),
		reader: bufio.NewReader(file),
	}, nil
}

func (it *FileIterator) Next() (string, error) {
	it.file.Seek(it.offset, io.SeekStart)
	bytes, err := it.reader.Read(it.buffer)
	it.offset += int64(bytes)
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
	offset := max(0, fileSize-bufferSize)
	resultOffset := int64(offset)
	isTail := true
	for totalLines <= n {
		file.Seek(resultOffset, io.SeekStart)
		buf := make([]byte, bufferSize)
		bytes, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return offset, err
		}
		if bytes == 0 {
			return offset, nil
		}
		bufferContent := string(buf[:bytes])
		if len(bufferContent) > 0 && isTail && bufferContent[len(bufferContent)-1] == '\n' {
			bufferContent = bufferContent[:len(bufferContent)-1]
			isTail = false
		}
		lineBreaksInBuffer := strings.Count(bufferContent, "\n")
		totalLines = totalLines + lineBreaksInBuffer
		if totalLines > n {
			linesToDrop := totalLines + 1 - n
			droppedLines := 0
			i := 0
			for droppedLines < linesToDrop {
				resultOffset = resultOffset + 1
				if bufferContent[i] == '\n' {
					droppedLines++
				}
				i++
			}
			return resultOffset, nil
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
