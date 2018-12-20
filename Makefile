CORES ?= $(shell nproc)

.PHONY: all _all clean proto test

all:
	@$(MAKE) _all --no-print-directory -j$(CORES)

_all: clean proto gofigure test

clean:
	rm -f proto/*.go gofigure-client gofigure-server coverage.out

proto:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto

gofigure: proto
	go build -o gofigure

test: proto
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
