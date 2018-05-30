#!/bin/sh

# Get all external libraries we need
go get -v \
      github.com/coreos/bbolt/... \
      github.com/lib/pq \
      github.com/peter-mount/golib/codec \
      github.com/peter-mount/golib/kernel \
      github.com/peter-mount/golib/rest \
      github.com/peter-mount/golib/statistics \
      github.com/peter-mount/golib/util \
      github.com/peter-mount/nre-feeds/util \
      gopkg.in/yaml.v2 \
      time

exit 0


gopkg.in/robfig/cron.v2 \
      github.com/gorilla/mux \
