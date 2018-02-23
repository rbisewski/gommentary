# Gommentary

Markdown documentation generator using program comments from SAS and Stata
code. Since old languages such as these tend to have odd comment styles, this
program assists in obtaining critical design information from previously
written code.

## Compilation

To build this program, ensure that golang and GNU Make is installed, and
use the below command:

`make`

Currently compilation works best on POSIX compatible machines that implement
the `date` command. Other architectures will build, but will have incomplete
version information.

## Usage

An example usage of the `gommentary` program is as follows:

`./gommentary -code-dir /path/to/application/code -docs-dir /path/to/application/code/docs`

Afterwards it will output the generated documentation into the the specified
`docs-dir` location, defaulting to a folder name of `docs/` inside of the
given `code-dir` location if no `docs-dir` path is specified.

Consider running the program with the `--help` flag for additional
information regarding these flags and what options are available.

## Testing

To run the current test suite of this program, type the following command:

`make test`

If all of the tests pass and are listed as ok, then the IO functions of this
program work as expected.

# TODOs

* Consider support for additional types of SAS / Stata comments.

* Add support for other programming languages, if possible.

* Implement an argument flag to allow the user to select the desired output
  filename.
