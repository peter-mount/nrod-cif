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
      git

# We want to build our final image under /dest
# A copy of /etc/ssl is required if we want to use https datasources
RUN mkdir -p /dest/etc &&\
    cp -rp /etc/ssl /dest/etc/

ADD scripts/ /scripts/

# Ensure we have the libraries - docker will cache these between builds
RUN /scripts/get.sh

# ============================================================
# source container contains the source as it exists within the
# repository.
FROM build as source
WORKDIR /go/src
ADD . .

# ============================================================
# Compile the source.
FROM source as compiler
ARG arch
ARG goos
ARG goarch
ARG goarm

# Build the microservice.
# NB: CGO_ENABLED=0 forces a static build
RUN CGO_ENABLED=0 \
    GOOS=${goos} \
    GOARCH=${goarch} \
    GOARM=${goarm} \
    /scripts/compile.sh /dest

# ============================================================
# Finally build the final runtime container for the specific
# microservice
#FROM scratch
FROM alpine

# The default database directory
Volume /database

# Install our built image
#COPY --from=compiler /dest/ /
COPY --from=compiler /dest/ /usr/local/bin/

#ENTRYPOINT ["/nrodcif"]
#CMD [ "-c", "/config.yaml"]
