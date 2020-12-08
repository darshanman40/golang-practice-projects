# CRUD API

API helps to execute Create/Read/Update/Delete operations on Mongo DB

## Technologies

- [gRPC](https://grpc.io/)
- [protobuf](https://developers.google.com/protocol-buffers)
- [mongodb client](https://github.com/mongodb/mongo-go-driver)

### TODO

- [ ] Mock services and unit test cases (part 1 added for blog-services)
- [ ] Write code for cpu and memory profiling
- [ ] Docker file for containerization

## Run unit Test

### Pre-requisite

- Running mongodb community server on
`mongodb://localhost:27017`
- Installed [dep](https://github.com/golang/dep) and run `dep ensure --vendor-only` to download all dependancies
- Installed [protoc](https://github.com/protocolbuffers/protobuf) and need to run following command under `crud-api` directory

`protoc ./internal/proto/*.proto --go_out=plugins=grpc:.`

To run tests using terminal, open `crud-api` directory (not `cmd/crud-api`) and run,

```go test -timeout 10s ./...```

To run tests and generate cover profile (`cp.out` file),

```go test -coverprofile cp.out -timeout 10s ./...```
