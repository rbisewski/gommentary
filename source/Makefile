# Version
YEAR = `date +%Y`
BUILD = `date +%s`
VERSION = `date +%y.%m`

# If unable to grab the year, default to ??
ifndef YEAR
    YEAR = "??"
endif

# If unable to grab the build, default to unknown
ifndef BUILD
    BUILD = "unknown"
endif

# If unable to grab the version, default to 0.0
ifndef VERSION
    VERSION = "0.0"
endif

#
# Makefile options
#


# State the "phony" targets
.PHONY: all clean build install uninstall


all: build

build:
	@echo 'Building gommentary...'
	@go build -o gommentary -ldflags '-s -w -X main.Version='${VERSION}' -X main.Build='${BUILD}' -X main.Year='${YEAR}

clean:
	@echo 'Cleaning...'
	@rm gommentary
