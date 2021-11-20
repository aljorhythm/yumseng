TAG ?= local

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
	docker run aljorhythm/yumseng:$(TAG)