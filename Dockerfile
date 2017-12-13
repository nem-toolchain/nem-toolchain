FROM golang:1.9.2-alpine3.7

RUN mkdir -p /go/src/github.com/nem-toolchain/nem-toolchain
ADD . /go/src/github.com/nem-toolchain/nem-toolchain
WORKDIR /go/src/github.com/nem-toolchain/nem-toolchain

RUN apk add --update make git
RUN ls
RUN make setup
RUN make build && make install
