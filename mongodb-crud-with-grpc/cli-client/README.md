# CLI Client for CRUD API

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
