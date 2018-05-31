#!/bin/sh

./scripts/build.sh test amd64 latest &&\
docker exec \
  postgres \
  psql \
    -U postgres \
    timetable \
    -c "delete from timetable.cif" &&\
docker run -it \
  --rm \
  -v /home/peter/area51/cif:/work:ro \
  -w /work \
  test:amd64-latest \
  cifimport \
    -d 'postgres://postgres:temppass@172.17.0.2/timetable?sslmode=disable' \
    tiplocs.cif \
    toc-update-mon.CIF \
    toc-update-tue.CIF \
    toc-update-wed.CIF
