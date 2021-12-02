# syntax=docker/dockerfile:1

FROM golang:1.16

WORKDIR /go/src/github.com/aljorhythm/yumseng

# tag
ARG TAG
ENV TAG=$TAG

# static assets
COPY webui webui

# local go files and packages
COPY *.go ./
COPY utils utils
COPY cheers cheers
COPY rooms rooms
COPY objectstorage objectstorage
COPY ping ping

# remote go packages
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# utils
COPY makefile ./

RUN make build

CMD [ "make", "run-binary" ]