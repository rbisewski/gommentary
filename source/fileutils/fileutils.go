package fileutils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// ReadFileIntoStringArray ... take the byte contents of a file and convert it into a string array.
func ReadFileIntoStringArray(path string) ([][]string, error) {

	path = strings.TrimSpace(path)
	if path == "" {
		panic("ReadFileIntoStringArray: passed an empty path")
	}

	// open the file and set a defer to close the file on function
	// return / panic
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	// initialize a new CSV reader instance and read in data
	csvReader := csv.NewReader(bufio.NewReader(file))
	csvReader.Comment = '#'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading file [%s]: %s",
			path, err)
	}

	if len(records) < 1 {
		return nil, fmt.Errorf("Error reading file [%s]: %s",
			path, "file is empty or does not contain valid records")
	}

	return records, nil
}

// WriteToFile ... Write string data to a file, with the option to overwrite.
func WriteToFile(path, data string, overwrite bool) error {

	path = strings.TrimSpace(path)
	if path == "" {
		panic("WriteToFile: passed an empty path")
	}
	if data == "" {
		return fmt.Errorf("WriteToFile: given data is empty")
	}

	// read-write, create if none exists or truncate existing one
	mode := os.O_RDWR | os.O_CREATE | os.O_TRUNC
	if !overwrite {
		mode |= os.O_EXCL //file must not exist
	}

	file, err := os.OpenFile(path, mode, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
