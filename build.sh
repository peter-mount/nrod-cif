#!/bin/sh

clear
docker build -t test . &&\
docker run -it --rm \
  --name test \
  -v /home/peter/nrod-data:/data:ro \
  test \
  cifimport /data/toc-full.CIF
