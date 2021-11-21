TAG ?= local
PORT ?= 8080
HOSTNAME ?= localhost

setup:
	# githooks
	git config core.hooksPath .githooks

format: setup
	sh .format.sh

unit-test:
	go test $(go list ./... | grep -v integration-tests)

run-source:
	go run main.go

build:
	go build -o ./yumseng

run-binary:
	./yumseng

docker-build: setup
	docker build . --build-arg TAG=$(TAG) --no-cache --tag aljorhythm/yumseng:$(TAG)

docker-run:
	docker run -d -p $(PORT):$(PORT) -e PORT=$(PORT) aljorhythm/yumseng:$(TAG)

integration-test:
	HOST=$(HOSTNAME):$(PORT) go test -v $$(go list ./... | grep integration-tests)

docker-run-undetached:
	docker run --expose=$(PORT) -p $(PORT):$(PORT) -e PORT=$(PORT) aljorhythm/yumseng:$(TAG)

docker-stop:
	docker ps -q --filter ancestor="aljorhythm/yumseng:$(TAG)" | xargs -r docker stop

all: docker-stop format unit-test docker-build docker-run integration-test docker-stop
	echo all done