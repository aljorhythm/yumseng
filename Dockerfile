# syntax=docker/dockerfile:1

FROM golang:1.16

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

COPY makefile ./

RUN make build

CMD [ "make", "run-binary", "${PATH}" ]