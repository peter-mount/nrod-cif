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

# Our build scripts
ADD scripts/ /usr/local/bin/

# Ensure we have the libraries - docker will cache these between builds
RUN get.sh

# ============================================================
# source container contains the source as it exists within the
# repository.
FROM build as source
WORKDIR /go/src/github.com/peter-mount/nrod-cif
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

# NB: CGO_ENABLED=0 forces a static build
RUN CGO_ENABLED=0 \
    GOOS=${goos} \
    GOARCH=${goarch} \
    GOARM=${goarm} \
    compile.sh /dest

# ============================================================
# This is the final image
FROM area51/scratch-base:latest
COPY --from=compiler /dest/ /
