# CLI Client for CRUD API

Simple command line (CLI) client to make gRPC request to [CRUD API](https://github.com/darshanman40/golang-practice-projects/tree/main/mongodb-crud-with-grpc/crud-api)

## Prequisite
- Running mongodb community server on `mongodb://localhost:27017`
- Installed dep and run `dep ensure` to download all dependancies
- Installed [protoc](https://github.com/protocolbuffers/protobuf) and need to run following command under `cli-client` directory
`protoc ./internal/proto/*.proto --go_out=plugins=grpc:.`
- Running `crud-api` server in same machine at port `50051` (Use [this](https://github.com/darshanman40/golang-practice-projects/tree/main/mongodb-crud-with-grpc/crud-api) guide)

## How to use

Create blog:

`go run main.go create -a {author name} -t {title name} -c {content}`

List all blogs:

`go run main.go list`

Read blog (Require unique id from monogo db):

`go run main.go read -i {id}`

Update blog (Require unique id from monogo db):

`go run main.go read -i {id} -a {author_name} -t {title name} -c {content}`

Delete blog (Require unique id from monogo db):

`go run main.go read -i {id}`
