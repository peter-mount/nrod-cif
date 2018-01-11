#!/bin/sh

clear
docker build -t test . || exit 1

#exit

#rm -f /home/peter/tmp/cif.db

docker run -it --rm \
  --name cifserver \
  -v /home/peter/tmp/:/database \
  -p 8081:8081 \
  test \
  cifserver \
  -p 8081 \
  -d /database/cif.db
