#############      builder       #############
FROM golang:1.12.1 AS builder

WORKDIR /go/src/github.com/fanfengqiang/cert-controller
COPY . .
RUN go install ./...



############# cert controller #############
FROM ubuntu:18.04
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install socat curl git tzdata -y && \
    apt-get autoclean && \
    rm -rf /var/lib/apt/lists/* && \
    git clone --depth=1 https://github.com/Neilpang/acme.sh.git && \
    cd acme.sh && \
    ./acme.sh install --force && \
    cd .. && \
    rm -rf acme.sh && \
    ln -s /root/.acme.sh/acme.sh /usr/bin/

COPY --from=builder /go/bin/cert-controller /cert-controller

ENTRYPOINT  ["/cert-controller"]