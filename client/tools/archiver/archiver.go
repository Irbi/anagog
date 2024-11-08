package tools

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

type Archiver struct {
	buf bytes.Buffer
	gz  *gzip.Writer
	err error
}

func Zip(data interface{}) (*bytes.Buffer, error) {
	buff := bytes.Buffer{}
	gz := gzip.NewWriter(&buff)

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Can't marshal", err)
		return &buff, err
	}
	_, err = gz.Write(b)
	if err != nil {
		fmt.Println("Can't gzip", err)
		return &buff, err
	}
	err = gz.Close()
	if err != nil {
		fmt.Println("Can't close gzip", err)
	}

	return &buff, nil
}
