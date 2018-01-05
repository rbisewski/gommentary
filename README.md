# Gommentary

Documentation generator using program comments for SAS and Stata code.

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

* Implement logic to obtain lists of the SAS macros used.

* Add support for other programming languages, if possible.
