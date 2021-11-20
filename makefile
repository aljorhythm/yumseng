TAG ?= local
PORT ?= 8080

format:
	sh .format.sh

unit-test:
	go test ./...

run-source:
	go run main.go

build:
	go build -o ./yumseng

run-binary:
	./yumseng

docker-build:
	docker build . --no-cache --tag aljorhythm/yumseng:$(TAG)

docker-run:
	docker run -d -e PORT=$(PORT) aljorhythm/yumseng:$(TAG)

docker-stop:
	docker ps -q --filter ancestor="aljorhythm/yumseng:$(TAG)" | xargs -r docker stop

all: format docker-build docker-run
	echo all done