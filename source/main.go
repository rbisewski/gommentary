/*~

# Summary
Documentation generator for SAS and Stata code, written in golang.

# Usage
Use the -h flag when running this program for basic usage information, and
consider reading the included doc.go for further details.

# Author
Robert Bisewski <robert.bisewski@umanitoba.ca>

Copyright (c) 2017 - Vaccine & Drug Evaluation Centre, Winnipeg, MB, Canada.
All rights reserved.

~*/

package main

import (
	"flag"
	"fmt"
	"os"
)

//
// Globals
//
var (
	// The version and build number are appended via Makefile, else default to
	// the below in the event that fails.
	Year    = "??"
	Version = "0.0"
	Build   = "unknown"

	// Whether or not to print the version + build information
	PrintVersionArgument = false

	// Code directory
	CodeDirectory = ""

	// Documentation directory
	DocumentationDirectory = "docs/"

	// File types with parsable comments
	ValidFiletypes = []string{".sas", ".do"}
)

//
// Program Main
//
func main() {

	err := setupArguments()
	if err != nil {
		fatal(err)
	}

	// if the version flag has been set to true, print the version
	// information and quit
	if PrintVersionArgument {
		fmt.Printf("Copyright %s, Vaccine & Drug Evaluation Centre, VDEC.ca\nDocumentation Generator for SAS and Stata code v%s, Build: %s\n",
			Year, Version, Build)
		os.Exit(0)
	}

	if err := validArgument(); err != nil {
		fmt.Println(usageMessage)
		fatal(err)
	}

	// default to storing generated docs in a "docs/" folder
	if DocumentationDirectory == "" {
		DocumentationDirectory = "docs/"
	}

	// attempt to read the contents of the code directory
	extractedComments, err := ReadCommentsFromAllFilesInDirectory(CodeDirectory, ValidFiletypes)
	if err != nil {
		fatal(err)
	}

	// create the docs directory; if it already exists nothing will
	// happen and the program will continue regardless
	err = os.MkdirAll(DocumentationDirectory, 0644)
	if err != nil {
		fatal(err)
	}

	// write the documentation to the docs directory
	err = WriteDocumentation(DocumentationDirectory, extractedComments)
	if err != nil {
		fatal(err)
	}

	os.Exit(0)
}

const (
	redColor = "\x1b[31m"
)

// Fatal prints error message in red and exits to shell with code 1
func fatal(err error) {
	fmt.Fprintf(os.Stderr, redColor+"%s\n", err)
	os.Exit(1)
}

// Setup the program arguments
func setupArguments() error {

	flag.Usage = func() {
		fmt.Println(usageMessage)
	}

	flag.StringVar(&CodeDirectory, "code-dir", "", "")
	flag.StringVar(&DocumentationDirectory, "docs-dir", "docs/", "")
	flag.BoolVar(&PrintVersionArgument, "version", false, "")

	flag.Parse()

	return nil
}

//validArgument returns an error if a necessary argument is missing
func validArgument() error {
	if CodeDirectory == "" {
		return fmt.Errorf("Invalid code directory path. Please enter a valid path and file.")
	}
	return nil
}
