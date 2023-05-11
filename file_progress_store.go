package file_progress_store

import (
	"encoding/json"
	"github.com/matt9mg/go-slice-diff"
	"github.com/pkg/errors"
	"io"
	"os"
)

// indexTypes specifies the types you can use to instantiate the FileProgressStore
type indexTypes interface {
	string | int | int8 | int16 | int32 | int64 | float32 | float64
}

type FileProgressStore[T indexTypes] struct {
	results []T
	file    *os.File
}

// NewFileProgressStore takes a storage file and its location and returns a FileProgressStore of its type or
// an error
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

// FileProgressStore appends the items to save
// an error is returned if the data cannot be saved
func (f *FileProgressStore[T]) Save(items ...T) error {
	f.results = append(f.results, items...)
	data, err := json.Marshal(f.results)

	if err != nil {
		return err
	}

	if err = f.file.Truncate(0); err != nil {
		return err
	}
	if _, err = f.file.Seek(0, 0); err != nil {
		return err
	}

	_, err = f.file.Write(data)

	return err
}

// ReturnUnprocessed checks to see if the items you wish to process have already been processed and stored within the
// current result set, this should be called be Save to check
func (f *FileProgressStore[T]) ReturnUnprocessed(toProcess []T) []T {
	return slice_diff.SliceDiff[T](f.results, toProcess)
}
