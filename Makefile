CORES ?= $(shell nproc)

.PHONY: all _all clean proto test

all:
	@$(MAKE) _all --no-print-directory -j$(CORES)

_all: clean proto gofigure-client gofigure-server test

clean:
	rm -f proto/*.go gofigure-client gofigure-server coverage.out

proto:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto

gofigure-server: proto
	go build -o gofigure-server server/*.go

gofigure-client: proto
	go build -o gofigure-client client/*.go

test: proto
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
