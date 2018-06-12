#!/bin/sh
#
# Builds all binaries
#
DEST=$1

for bin in \
  cifimport \
  cifrest \
  cifretrieve
do
  echo "Building ${bin}"
  go build -o ${DEST}/${bin} github.com/peter-mount/nrod-cif/${bin}/bin
done
