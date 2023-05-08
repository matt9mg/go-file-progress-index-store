package file_progress_store

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"os"
)

type indexTypes interface {
	string | int | int8 | int16 | int32 | int64 | float32 | float64
}
type FileProgressStore[T indexTypes] struct {
	results []T
	file    *os.File
}

func NewFileProgressStore[T indexTypes](fileStore string) (*FileProgressStore[T], error) {
	if _, err := os.Stat(fileStore); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(fileStore); err != nil {
			return nil, errors.New("unable to create file " + fileStore + " with error: " + err.Error())
		}
	}

	f, err := os.OpenFile(fileStore, os.O_RDWR, os.ModePerm)

	if err != nil {
		return nil, errors.New("unable to open file " + fileStore + " with error: " + err.Error())
	}

	fps := &FileProgressStore[T]{
		file: f,
	}

	if err = json.NewDecoder(f).Decode(&fps.results); err == io.EOF {
		return fps, nil
	}

	if err == nil {
		return fps, nil
	}

	return nil, errors.New("unable to bind results to FileProgressStore with error: " + err.Error())
}

func (f *FileProgressStore[T]) Save(items ...T) error {
	f.results = append(f.results, items...)
	data, err := json.Marshal(f.results)

	if err != nil {
		return err
	}

	if err = f.file.Truncate(0); err != nil {
		return err
	}
	_, err = f.file.Seek(0, 0)

	if err != nil {
		return err
	}

	_, err = f.file.Write(data)

	return err
}
