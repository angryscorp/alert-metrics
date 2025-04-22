package filewriter

import (
	"io"
	"os"
)

type FileWriter struct {
	fileStoragePath string
}

var _ io.Writer = (*FileWriter)(nil)

func New(fileStoragePath string) FileWriter {
	return FileWriter{fileStoragePath: fileStoragePath}
}

func (f FileWriter) Write(p []byte) (n int, err error) {
	file, err := os.Create(f.fileStoragePath)
	if err != nil {
		return 0, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	return file.Write(p)
}
