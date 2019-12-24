PERCENT := %

.PHONY: all clean fmt lint proto test bench upx

all: clean fmt lint proto gofigure-minion example-master test upx

clean:
	@rm -rf proto/*.go gofigure-minion example-master cov/* docs

fmt: proto
	@go fmt ./...

lint: fmt
	@golangci-lint run

proto: clean
	@protoc -I proto/ --go_out=plugins=grpc,paths=source_relative:proto proto/*.proto

gofigure-minion: proto fmt
	@CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o gofigure-minion ./cmd/minion

example-master: proto fmt
	@CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o example-master ./cmd/example

docker-sshd:
	@chmod 600 testdata/ssh_host_rsa_key*
	cd testdata && docker build -t gofigure-sshd .

test: proto fmt docker-sshd gofigure-minion
	@go test -coverprofile=cov/coverage.out --coverpkg="$$(go list ./... | grep -v /proto | paste -sd, -)" ./...
	@go tool cover -func=cov/coverage.out
	@go tool cover -html=cov/coverage.out -o cov/coverage.html

bench: proto fmt
	@go test -bench=. ./...

upx: gofigure-minion
	@upx -qq gofigure-minion
