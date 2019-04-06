FROM golang:alpine
RUN apk add --no-cache git protobuf gcc musl-dev make upx
RUN wget -O - https://github.com/golangci/golangci-lint/releases/download/v1.16.0/golangci-lint-1.16.0-linux-amd64.tar.gz | tar -xz --strip-components 1 -C /usr/bin golangci-lint-1.16.0-linux-amd64/golangci-lint

ENV GO111MODULE=on GOOS=linux GOARCH=amd64 MAKEFLAGS="-j2"
WORKDIR /go/src/github.com/alexhunt7/gofigure

COPY ./go.mod ./go.sum /go/src/github.com/alexhunt7/gofigure/
RUN go mod download
# I have no idea how to properly get the same protoc-gen-go as my protobuf version
RUN go get github.com/golang/protobuf/protoc-gen-go@$(awk '/github.com\/golang\/protobuf/ {print $2}' go.mod)

COPY ./ /go/src/github.com/alexhunt7/gofigure/
RUN make lint upx
