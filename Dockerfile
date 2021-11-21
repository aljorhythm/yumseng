# syntax=docker/dockerfile:1

FROM golang:1.16

WORKDIR /go/src/github.com/aljorhythm/yumseng

# static assets
COPY webui webui

# local go files and packages
COPY .tag ./
COPY *.go ./
COPY utils utils
COPY cheers cheers
COPY ping ping

# remote go packages
COPY go.mod ./
RUN go mod download

# utils
COPY makefile ./

RUN make build

CMD [ "make", "run-binary" ]