package fileutils

import (
	"testing"
)

const (
	a3RecCSV = `
	v1,v2,v3
	1, 1.2, one
	2, 2.2, two
	3, 3.3, three`
)

func TestWriteToFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		data    string
		wantErr bool
	}{
		{"wrong path", "doesnotexist/file.txt", "", true},
		{"empty file", "../testing/empty.txt", "", true},
		{"3 record csv file", "../testing/3rec.csv", a3RecCSV, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteToFile(tt.path, tt.data, false); (err != nil) != tt.wantErr {
				t.Errorf("WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//WARNING: requires files created in TestWriteToFile() above
func TestReadFileIntoStringArray(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		len     int
		data    string
		wantErr bool
	}{
		{"wrong path", "doesnotexist/file.txt", 0, "", true},
		{"empty file", "../testing/empty.txt", 0, "", true},
		{"3 record csv file", "../testing/3rec.csv", 4, a3RecCSV, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strMatrix, err := ReadFileIntoStringArray(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReadFileIntoStringArray: error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(strMatrix) != tt.len {
				t.Errorf("TestReadFileIntoStringArray: wanted length of %d got %d", tt.len, len(strMatrix))
			}
		})
	}
}
