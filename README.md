[![Go Reference](https://pkg.go.dev/badge/github.com/matt9mg/go-file-progress-index-store.svg)](https://pkg.go.dev/github.com/matt9mg/go-file-progress-index-store)

# File Progress Index Store
Stores current progress of a process index

### Installation
```
go get github.com/matt9mg/go-file-progress-index-store
```

### Examples
```go
fpis, err := file_progress_store.NewFileProgressStore[string]("/my/file/location/file.json")

if err != nil {
	log.Fatalln(err)
}

data := []string{"1", "2", "3"}

unprocessed := fpis.ReturnUnprocessed(data)

if len(unprocessed) > 0 {
	// perform some business logic
	if err := fpis.Save(unprocessed...); err != nil {
		log.Fatalln(err)
    }
}
```
