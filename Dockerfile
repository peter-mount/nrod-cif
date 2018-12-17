# ============================================================
# Dockerfile used to build the nrod-cif microservice
# ============================================================

ARG arch=amd64
ARG goos=linux

# ============================================================
# Build container containing our pre-pulled libraries.
# As this changes rarely it means we can use the cache between
# building each microservice.
FROM golang:alpine as build

# The golang alpine image is missing git so ensure we have additional tools
RUN apk add --no-cache \
      curl \
      git \
      tzdata

# go-bindata
RUN go get -v github.com/kevinburke/go-bindata &&\
    go build -o /usr/local/bin/go-bindata github.com/kevinburke/go-bindata/go-bindata

# Ensure we have the libraries - docker will cache these between builds
WORKDIR /work
ADD go.mod .
RUN go mod download

#RUN go get -v \
#      github.com/lib/pq \
#      github.com/peter-mount/golib/... \
#      github.com/peter-mount/nre-feeds/util \
#      github.com/peter-mount/sortfold

# ============================================================
# source container contains the source as it exists within the
# repository.
FROM build as source
WORKDIR /work
ADD . .

# Import sql so we can build as needed
RUN go-bindata -o cifimport/sqlassets.go -pkg cifimport sql/

# ============================================================
# Compile the source.
FROM source as compiler
ARG arch
ARG goos
ARG goarch
ARG goarm
WORKDIR /work

# NB: CGO_ENABLED=0 forces a static build
RUN for bin in \
      cifimport \
      cifrest \
      cifretrieve;\
    do\
      echo "Building ${bin}";\
      CGO_ENABLED=0 \
          GOOS=${goos} \
          GOARCH=${goarch} \
          GOARM=${goarm} \
          OUT= \
      go build \
        -o /dest/bin/${bin} \
        ./${bin}/bin;\
    done

# ============================================================
# This is the final image
FROM alpine
RUN apk add --no-cache \
      curl \
      tzdata
COPY --from=compiler /dest/ /
