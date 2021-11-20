# syntax=docker/dockerfile:1

FROM golang:1.16

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

COPY makefile ./

RUN make build

EXPOSE 8080

CMD [ "make", "run-binary" ]