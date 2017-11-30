/*
 * Contains a list of structure definitions
 */

package main

// Comment object definition
type Comment struct {

	// text value of keyword, as defined via the @ symbol; blank means not a keyword comment
	Keyword string

	// which keyword this comment should be grouped under
	GroupUnder string

	// path to the file the comment was found in, mostly useful for debugging purposes
	Filename string

	// index position representing the order the file was read in
	Index int

	// line number that the comment was obtained on
	LineNum int

	// ascii content of the given comment
	Text string
}
