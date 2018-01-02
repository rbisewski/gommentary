/*
 * Useful functions for reading comments and writing documentation files
 */

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ReadCommentsFromAllFilesInDirectory ... search through all files in a given directory for comments
// TODO: add logic to this file to handle the "group under" functionality
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
	count := 0
	for _, path := range listOfFilesToRead {

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
		count++

		// attach index to comments and append them
		for _, cmt := range parsed {
			cmt.Filename = path
			cmt.Index = count
			comments = append(comments, cmt)
		}
	}

	return comments, nil
}

// GetLineNumber ... obtain the current line number a comment as defined by (startIndex, endIndex) appears on
func GetLineNumber(lines [][]int, pos []int) (int, error) {
	if len(lines) == 0 || len(pos) != 2 {
		panic("Invalid string or comment indices given during line reconstruction attempt.")
	}

	end := pos[1]

	for i, line := range lines {
		if len(line) != 2 {
			return -1, fmt.Errorf("Line mismatched, file is likely corrupted.")
		}
		// check against the end of the next line
		if line[1] >= end {
			return i + 1, nil
		}
	}

	return -1, fmt.Errorf("Comment not found, file is likely corrupted.")
}

// ParseStringForComments ... obtain all comments from a given string
// TODO: functionalize and clean up parts of the regex logic used
func ParseStringForComments(contents string) ([]Comment, error) {
	if contents == "" {
		panic("A given file has unparsable contents.")
	}

	asterixComment := regexp.MustCompile("@[^@]+")
	whitespaceRegexes := []string{"\t", "\r", "\f", "\v"}
	commentStrings := make([]RawComment, 0)
	comments := make([]Comment, 0)

	// obtain newline indices, helpful for reconstructing line numbers
	newlineRegex := "\n"
	reNewline := regexp.MustCompile(newlineRegex)
	lineIndices := reNewline.FindAllStringIndex(contents, -1)

	//
	// Clean away complicated whitespace characters
	//

	// clean up unneeded whitespace characters
	contents = reNewline.ReplaceAllString(contents, " ")
	for _, str := range whitespaceRegexes {
		re := regexp.MustCompile(str)
		contents = re.ReplaceAllString(contents, "")
	}

	//
	// Handle the different comment types here
	//

	// buffer the contents with a space to allow for regexes that check for initial characters
	contents = " " + contents

	// handle the |**@keyword ;| comments
	twoAsterixAtEndsWithSemicolon := regexp.MustCompile("[^\\/]\\s{0,}\\*\\*@[a-zA-Z\\.]+ [^;]+;")
	sindices := twoAsterixAtEndsWithSemicolon.FindAllStringIndex(contents, -1)
	for _, sindex := range sindices {
		start := sindex[0] + 1
		end := sindex[1]
		// take into account very short / empty comments
		if start == end {
			continue
		}
		// get the line number the comment starts on
		num, err := GetLineNumber(lineIndices, sindex)
		if err != nil {
			return nil, err
		}
		raw := contents[start:end]
		commentStrings = append(commentStrings, RawComment{num, raw})
	}

	// handle the |** ;| comments
	twoAsterixEndsWithSemicolon := regexp.MustCompile("[^\\/]\\s{0,}\\*\\*[a-zA-Z\\.]+ [^;]+;")
	sindices = twoAsterixEndsWithSemicolon.FindAllStringIndex(contents, -1)
	for _, sindex := range sindices {
		start := sindex[0] + 1
		end := sindex[1]
		// take into account very short / empty comments
		if start == end {
			continue
		}
		// get the line number the comment starts on
		num, err := GetLineNumber(lineIndices, sindex)
		if err != nil {
			return nil, err
		}
		raw := contents[start:end]
		commentStrings = append(commentStrings, RawComment{num, raw})
	}

	// handle the |/**@ */| comments
	slashTwoAsterixAtEndsWithSlash := regexp.MustCompile("\\/\\s{0,}\\*\\*\\s{0,}@[a-zA-Z\\.]+ [^\\/]+\\*\\/")
	sindices = slashTwoAsterixAtEndsWithSlash.FindAllStringIndex(contents, -1)
	for _, sindex := range sindices {
		start := sindex[0] + 1
		end := sindex[1]
		// take into account very short / empty comments
		if start == end {
			continue
		}
		// get the line number the comment starts on
		num, err := GetLineNumber(lineIndices, sindex)
		if err != nil {
			return nil, err
		}
		raw := contents[start:end]

		// if there is only a single @ comment, just use the whole match
		pieces := asterixComment.FindAllString(raw, -1)
		if len(pieces) == 1 {
			commentStrings = append(commentStrings, RawComment{num, raw})
			continue
		}

		// handle comments of the type...
		//
		//    /**
		//     @main :title   Experiment #42
		//     @main :author  John Smith
		//     @main :org     University of Manitoba
		//     @main This experiment is designed to take into account the answer to life, the universe, and everything.
		//    */
		//
		for _, piece := range pieces {
			commentStrings = append(commentStrings, RawComment{num, piece})
		}
	}

	//
	// Convert the raw comment text into meaningful comments
	//

	// attempt to convert the above comment strings to comments
	asterixCommentWithSpace := regexp.MustCompile("@[^@\\s]+\\s")
	for _, str := range commentStrings {
		newComment := Comment{"", "", "", 0, str.LineNum, ""}

		// obtain the keyword, if any
		match := asterixCommentWithSpace.FindString(str.Text)

		// handle the comments that have keywords...
		if match != "" {
			newComment.Keyword = match
			text := asterixCommentWithSpace.Split(str.Text, -1)
			if len(text) > 1 {
				newComment.Text = text[1]
			}

		} else {
			// ... else just use the whole string as a comment
			newComment.Text = str.Text
		}

		// cleanup text
		newComment.Text = strings.TrimSpace(newComment.Text)
		newComment.Text = strings.TrimSuffix(newComment.Text, ";")
		newComment.Text = strings.TrimPrefix(newComment.Text, "**")
		newComment.Text = strings.TrimPrefix(newComment.Text, "/*")
		newComment.Text = strings.TrimSuffix(newComment.Text, "*/")
		newComment.Text = strings.TrimSpace(newComment.Text)

		comments = append(comments, newComment)
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

	markdownContents := ""

	//
	// Title comments
	//
	for _, cmt := range comments {

		// skip normal comments
		if cmt.Keyword == "" || strings.Index(cmt.Text, ":") != 0 {
			continue
		}

		// assemble title / author / organization / version information
		re := regexp.MustCompile(":[a-zA-Z\\.]+ ")
		match := re.FindString(cmt.Text)
		text := strings.Split(cmt.Text, match)
		if len(text) < 2 {
			return fmt.Errorf("Improperly formatted title comment.")
		} else if strings.ToLower(match) == ":version " {
			markdownContents += "% Version " + strings.Title(text[1]) + "\n"
		} else {
			markdownContents += "% " + strings.Title(text[1]) + "\n"
		}
	}

	//
	// Files and scripts used in project
	//
	order := 1
	keywordsMap := make(map[string]int)
	filesMap := make(map[int]string)
	markdownContents += "\n# Scripts/macros used for project\n\n"
	for _, cmt := range comments {

		// skip title comments
		if cmt.Keyword != "" && strings.Index(cmt.Text, ":") == 0 {
			continue
		}

		filesMap[cmt.Index] = cmt.Filename

		// add keyword subtitles
		trimmedKeyword := strings.TrimSpace(cmt.Keyword)
		trimmedKeyword = strings.Trim(trimmedKeyword, "@")
		if trimmedKeyword != "" && keywordsMap[trimmedKeyword] == 0 {
			keywordsMap[trimmedKeyword] = order
			order++
		}
	}
	for i := 1; i <= len(filesMap); i++ {
		indexAsString := strconv.FormatInt(int64(i), 10)
		markdownContents += "* " + indexAsString + ": " + filesMap[i] + "\n"
	}

	//
	// Normal comments
	//
	for i := 1; i <= len(keywordsMap); i++ {

		currentKeyword := ""

		for key, value := range keywordsMap {
			if i == value {
				currentKeyword = key
			}
		}

		if currentKeyword != "" {
			markdownContents += "\n# " + strings.Title(currentKeyword) + "\n\n"
		}

		counter := 1
		for _, cmt := range comments {

			// skip title comments
			if cmt.Keyword != "" && strings.Index(cmt.Text, ":") == 0 {
				continue
			}

			// skip if the comment is not associated with that keywords
			trimmedKeyword := strings.TrimSpace(cmt.Keyword)
			trimmedKeyword = strings.Trim(trimmedKeyword, "@")
			if trimmedKeyword != currentKeyword {
				continue
			}

			indexAsString := strconv.FormatInt(int64(cmt.Index), 10)
			counterAsString := strconv.FormatInt(int64(counter), 10)
			lineNumberAsString := strconv.FormatInt(int64(cmt.LineNum), 10)
			if strings.HasSuffix(cmt.Filename, ".do") {
				indexAsString = "s" + indexAsString
			}

			markdownContents += indexAsString + "." + counterAsString + ":" + lineNumberAsString + " " + cmt.Text + "\n"

			counter++
		}
	}

	// TODO: implement logic to make this write it out to a markdown file, etc
	fmt.Println(markdownContents)

	return nil
}
