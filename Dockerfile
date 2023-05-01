# Build Geth into a first stage container
FROM golang:1.20-bullseye as build-gramine

RUN apt-get update && \
    apt-get install -y curl libssl-dev build-essential ca-certificates git

WORKDIR /geth
ADD Makefile geth-patches /geth/
ADD geth-patches /geth/geth-patches
RUN make PROTECT=1 TLS=1

# Pull Geth into a second stage deploy container
FROM debian:bullseye

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates

WORKDIR /geth
COPY --from=build-gramine /geth/geth ./
ADD entrypoint.sh /geth/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["/geth/entrypoint.sh"]

# Add some metadata labels to help programatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"
