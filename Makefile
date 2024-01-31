#
# By default this will build the project on every non-mobile platform
# supported by the installed go environment.
#
# To limit a build to a single environment, you can force it to just a
# single platform by prefixing make with:
#
# PLATFORMS=linux:amd64: make clean all
#
# Just change the entry for your OS and CPU. These are listed in platforms.md
#
# Note: For 32 bit arm processors the 3rd parameter is important.
# e.g. use linux:arm:6 or linux:arm:7
#
# For all other processors, including arm64, leave the third field blank.
#
# For a parallel builds you can use the -j parameter to make as usual.
#
# e.g.: make -j 8 clean all
#
# Pick a value suitable to the number of cores/thread your machine has.
# This is useful for a full build of all platforms as it will build all
# of the binaries in parallel speeding up the full build.
#

.PHONY: all clean init test build

all: init test build

init:
	@echo "GO MOD   tidy";go mod tidy
	@echo "GO MOD   download";go mod download
	@echo "GENERATE build";\
	CGO_ENABLED=0 go build -o build tools/build/bin/main.go
	@./build -build Makefile.gen -build-platform "$(PLATFORMS)" -d builds -dist dist -build-archiveArtifacts "dist/*"

go-bindata:
	@if [ ! -f go-bindata ]; then echo "CURL     go-bindata";\
	curl --silent --location --output go-bindata https://github.com/kevinburke/go-bindata/releases/download/v3.25.0/go-bindata-linux-amd64;\
	chmod 755 go-bindata;\
	fi

clean: init
	@${MAKE} --no-print-directory -f Makefile.gen clean

test: init
	@${MAKE} --no-print-directory -f Makefile.gen test

build: cifimport/sqlassets.go test
	@${MAKE} --no-print-directory -f Makefile.gen all

cifimport/sqlassets.go: go-bindata
	@echo "GO GENERATE sqlassets";./go-bindata -o tools/cifimport/sqlassets.go -pkg cifimport sql/

docs: init
	@${MAKE} --no-print-directory -f Makefile.gen docs
