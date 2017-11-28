/*
 * Useful functions for reading comments and writing documentation files
 */

package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// ReadCommentsFromAllFilesInDirectory ... search through all files in a given directory for comments
// TODO: implement this function
func ReadCommentsFromAllFilesInDirectory(codeDir string, filetypes []string) ([]Comment, error) {

	if codeDir == "" {
		panic("Code directory name is invalid")
	}

	comments := make([]Comment, 0)

	codeDirContents, err := ioutil.ReadDir(CodeDirectory)
	if err != nil {
		return nil, err
	}

	for i, file := range codeDirContents {

		// skip directory pointers
		if file.Name() == "." || file.Name() == ".." {
			continue
		}

		// check if file is a valid type
		isAcceptedFiletype := false
		for _, t := range filetypes {
			if strings.HasSuffix(file.Name(), t) {
				isAcceptedFiletype = true
				break
			}
		}

		// skip files that are non-accepted file types
		if !isAcceptedFiletype {
			continue
		}

		// TODO: insert further logic here
		i = i
	}

	return comments, nil
}

// WriteDocumentation ... generate documentation using the comments and write it out to file
// TODO: implement this function
func WriteDocumentation(docsDir string, comments []Comment) error {

	if docsDir == "" {
		panic("Docs directory name is invalid")
	}
	if len(comments) < 1 {
		return fmt.Errorf("No comments were present in the files. Exiting...")
	}

	return nil
}
