/*
 * Useful functions for reading comments and writing documentation files
 */

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// ReadCommentsFromAllFilesInDirectory ... search through all files in a given directory for comments
func ReadCommentsFromAllFilesInDirectory(codeDir string, filetypes []string) ([]Comment, error) {

	if codeDir == "" {
		panic("Code directory name is invalid")
	}

	listOfFilesToRead := make([]string, 0)
	comments := make([]Comment, 0)

	codeDirContents, err := ioutil.ReadDir(CodeDirectory)
	if err != nil {
		return nil, err
	}

	// obtain the list of files to read
	for _, file := range codeDirContents {

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

		filename := filepath.Join(codeDir, file.Name())

		listOfFilesToRead = append(listOfFilesToRead, filename)
	}

	if len(listOfFilesToRead) < 1 {
		return nil, fmt.Errorf("No parsable files were found. Exiting...")
	}

	// using the list of files, read each of them
	for i, path := range listOfFilesToRead {

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		contents := string(bytes)

		// if file is empty, skip it
		if contents == "" {
			continue
		}

		parsed, err := ParseStringForComments(contents)
		if err != nil {
			return nil, err
		}

		// if no comments, skip it
		if len(parsed) < 1 {
			continue
		}

		// attach index to comments and append them
		for _, cmt := range parsed {
			cmt.Index = i
			comments = append(comments, cmt)
		}
	}

	return comments, nil
}

// ParseStringForComments ... obtain all comments from a given string
// TODO: implement this function
func ParseStringForComments(contents string) ([]Comment, error) {
	if contents == "" {
		panic("A given file has unparsable contents.")
	}

	whitespaceRegexes := []string{"\n", "\t", "\r", "\f", "\v"}
	commentStrings := make([]string, 0)
	comments := make([]Comment, 0)

	// clean up unneeded whitespace characters
	for _, str := range whitespaceRegexes {
		re := regexp.MustCompile(str)
		contents = re.ReplaceAllString(contents, "")
	}

	// handle the |**@keyword ;| comments
	twoAsterixAtEndsWithSemicolon := regexp.MustCompile("[^\\/]\\s{0,}\\*\\*@[a-zA-Z]+ [^;]+;")
	matches := twoAsterixAtEndsWithSemicolon.FindAllString(contents, -1)
	for _, match := range matches {
		match = match[1:]
		match = strings.TrimSpace(match)
		commentStrings = append(commentStrings, match)
	}

	// TODO: check for other comments

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
