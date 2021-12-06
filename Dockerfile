# syntax=docker/dockerfile:1

FROM node:17 AS fe

WORKDIR /app

COPY react-ts-ui/public react-ts-ui/public
COPY react-ts-ui/src react-ts-ui/src
COPY react-ts-ui/makefile react-ts-ui/makefile
COPY react-ts-ui/package.json react-ts-ui/package.json
COPY react-ts-ui/package-lock.json react-ts-ui/package-lock.json

ENV NODE_OPTIONS=--openssl-legacy-provider

WORKDIR /app/react-ts-ui
RUN make build

FROM golang:1.16

WORKDIR /go/src/github.com/aljorhythm/yumseng
COPY --from=fe /app/react-ts-ui/build ./react-ts-ui/build

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