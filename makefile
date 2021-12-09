TAG ?= local
PORT ?= 8080
HOSTNAME ?= localhost

setup:
	# githooks
	git config core.hooksPath .githooks
	go mod tidy

format: setup
	sh .format.sh

unit-test:
	go test -v $$(go list ./... | grep -v integration-tests)

run-source:
	go run .

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

docker-deploy-run: docker-stop docker-build docker-run-undetached
	echo built and run and exited server

source-integration-tests:
	sh .source-integration-tests.sh

all: docker-stop format unit-test docker-build docker-run integration-test docker-stop
	echo all done