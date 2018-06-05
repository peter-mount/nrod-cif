#!/bin/sh

# Get all external libraries we need
go get -v \
      github.com/lib/pq \
      github.com/peter-mount/golib/kernel \
      github.com/peter-mount/golib/rest \
      github.com/peter-mount/golib/util \
      github.com/peter-mount/nre-feeds/util
