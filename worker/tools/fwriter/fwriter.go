package fwriter

import (
	"fmt"
	"os"
)

func CreateFile(fname string) (string, error) {
	fmt.Println("File path " + fname)
	newFile, newFileErr := os.Create(fname)
	if newFileErr != nil {
		fmt.Println("Error creating file")
		fmt.Println(newFileErr)
		return "", newFileErr
	}
	if err := newFile.Close(); err != nil {
		fmt.Println("Error closing file", err)
	}

	return fname, nil
}

func AppendLines(fPath, data string) {
	f, err := os.OpenFile(fPath, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(fPath)
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(f, data)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fPath + "file appended successfully")
}
