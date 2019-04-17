# gofigure
[![GoDoc](https://godoc.org/github.com/alexhunt7/gofigure?status.svg)](https://godoc.org/github.com/alexhunt7/gofigure)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexhunt7/gofigure)](https://goreportcard.com/report/github.com/alexhunt7/gofigure)

GoFigure is a configuration management tool with remote agents written in golang, and communication using GRPC. Because it communicates over GRPC, any protocol-buffered language is supported as a client. The reference client is written in golang.

Primary design goals:
- Speed
- Testability
- Proper variable management
- Prevent user error (force best practices) 
- Flexibility
- Simplicity

GoFigure is still in development, but has most major features completed. Use caution if deploying to a production environment.

There are two main libraries in gofigure:
### Minion
Minion is the remote agent of gofigure. It listens for incoming client connections to tell it what to do.

### Master
Master is a reference client. You don't have to use it, but if you're writing your client in golang, it has several helpful functions.

If using the reference client, you can easily build a single binary with both the master and minion, and can use the built in `Bootstrap` function to start the remote minions over SSH. If using some other client, you will need to handle this yourself.


Connection example with reference client:
```go
var gofigureClient *Client
creds, err := credentials.Load(caFile, certFile, keyFile)
if err != nil {
    return gofigureClient, err
}
conn, err := master.ConnectGRPC("123.123.123.123:10000", grpc.WithTransportCredentials(creds))
if err != nil {
    return gofigureClient, err
}
return &master.Client{GofigureClient: pb.NewGofigureClient(conn)}, nil
```
