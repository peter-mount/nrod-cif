#!/bin/sh

# Get all external libraries we need
go get -v \
      github.com/lib/pq \
      github.com/peter-mount/golib/... \
      github.com/peter-mount/nre-feeds/util \
      github.com/peter-mount/sortfold \
