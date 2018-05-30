#!/bin/sh
#
# Builds all binaries
#
DEST=$1

for bin in \
  cifimport
do
  echo "Building ${bin}"
  go build -o ${DEST}/${bin} ${bin}
done
