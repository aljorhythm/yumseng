TAG ?= local
PORT ?= 8080

format:
	go fmt ./...

run-source:
	go run main.go

build:
	go build -o ./yumseng

run-binary:
	./yumseng

docker-build:
	docker build . --no-cache --tag aljorhythm/yumseng:$(TAG)

docker-run:
	docker run -e PORT=$(PORT) aljorhythm/yumseng:$(TAG)

all: format docker-build docker-run
	echo all done