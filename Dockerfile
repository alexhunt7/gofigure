FROM golang:alpine
RUN apk add --no-cache git protobuf gcc musl-dev

ENV GO111MODULE=on GOOS=linux GOARCH=amd64 MAKEFLAGS="-j2"
WORKDIR /go/src/github.com/alexhunt7/gofigure

COPY ./go.mod ./go.sum /go/src/github.com/alexhunt7/gofigure/
RUN go mod download
# I have no idea how to properly get the same protoc-gen-go as my protobuf version
RUN go get github.com/golang/protobuf/protoc-gen-go@$(awk '/github.com\/golang\/protobuf/ {print $2}' go.mod)

COPY ./proto/ /go/src/github.com/alexhunt7/gofigure/proto/
RUN protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto

COPY ./client/ /go/src/github.com/alexhunt7/gofigure/client/
RUN go build -o gofigure-client client/*.go

COPY ./server/ /go/src/github.com/alexhunt7/gofigure/server/
RUN go build -o gofigure-server server/*.go
