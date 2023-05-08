package file_progress_store

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

type indexTypes interface {
	string | int | int8 | int16 | int32 | int64 | float32 | float64
}

type FileProgressStore[T indexTypes] struct {
	Results []T
}

func NewFileProgressStore[T indexTypes](fileStore string) (*FileProgressStore[T], error) {
	if _, err := os.Stat(fileStore); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(fileStore); err != nil {
			return nil, errors.New("unable to create file " + fileStore + " with error: " + err.Error())
		}
	}

	f, err := os.Open(fileStore)

	if err != nil {
		return nil, errors.New("unable to open file " + fileStore + " with error: " + err.Error())
	}

	fps := &FileProgressStore[T]{}

	if err = json.NewDecoder(f).Decode(&fps.Results); err != nil {
		return fps, errors.New("unable to bind results to FileProgressStore with error: " + err.Error())
	}

	return fps, nil
}
