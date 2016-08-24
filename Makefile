
SRC=./cmd ./engine ./util

.PHONY: \
	check \
	build \
	run \
	test \
	install \
	format \
	image \
	clean

all: check test

check:
	go tool vet --all $(SRC)
	go tool vet --all *.go
	for d in $(SRC); do \
		golint $$d;\
	done

build:
	go build main.go

run:
	go run main.go start

test:
	go test $(SRC)

install: build
	[ -d /ckeeper ] || sudo mkdir /ckeeper && sudo install ckeeper /ckeeper/

fmt:
	goimports -w  $(SRC)
	goimports -w *.go
	gofmt -s -w  $(SRC)
	gofmt -s -w *.go

image:
	docker build -t yeasy/ckeeper .

clean:
	go clean ./...
