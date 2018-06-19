#!/bin/sh
#
# Builds all binaries
#
DEST=$1

BIN_DIR=${DEST}/bin/

mkdir -p ${BIN_DIR}

for bin in \
  cifimport \
  cifrest \
  cifretrieve
do
  echo "Building ${bin}"
  OUT=
  go build \
    -o ${BIN_DIR}/${bin} \
    github.com/peter-mount/nrod-cif/${bin}/bin
done
