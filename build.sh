#!/bin/sh

clear
docker build -t test . || exit 1

#exit

#rm -f /home/peter/tmp/cif.db

if [ -z "$1" ]
then
  # Standalone run
  docker run -it --rm \
    --name cifserver \
    -v /home/peter/tmp/:/database \
    -p 8081:8081 \
    test \
    cifserver \
    -p 8081 \
    -d /database/cif.db
else
  # Run via traefic
  docker run -it --rm \
    --name cifserver \
    -v /home/peter/tmp/:/database \
    -l traefik.backend=nrod-cif \
    -l traefik.docker.network=bridge \
    -l traefik.frontend.rule=Host:$1 \
    -l traefik.enable=true \
    -l traefik.port=8080 \
    -l traefik.protocol=http \
    test \
    cifserver \
    -d /database/cif.db
fi
