# Runner service
The Runner service allows you to execute and share Go code snippets.

Protobuf contract: [runner](https://github.com/paniccaaa/protos/blob/main/proto/runner/runner.proto) 

## Stack

- **gRPC** / gRPC-gateway
- **PostgreSQL**
- **cleanenv**
- **slog**

## Installation

Before running the application, make sure Go are installed on your system.

1. **Clone** the repository
```bash
git clone https://github.com/paniccaaa/runner.git
```
2. **Create** the .env, local.yaml and Makefile (see .example in repository)
3. **Run** the app locally
```bash
make run 
```

## Endpoints

### REST (gRPC-gateway)
- **POST /run**: Retrieve a list of all posts.
- **POST /share**: Create a new post.
- **GET /shared/{id}**: Retrieve a post by its identifier.
  
### gRPC
- **RunCode**
```bash
$ grpcurl -plaintext -d '{"code": "package main\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"Hello RunCode!\")\n}"}' \
 localhost:44000 runner.Runner/RunCode

{
  "code": "package main\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"Hello RunCode!\")\n}",
  "output": "Hello RunCode!\n"
}
```
- **ShareCode**
```bash
$ grpcurl -plaintext -d '{"code": "package main\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"Hello world! Test ShareCode\")\n}"}' \
 localhost:44000 runner.Runner/ShareCode

{
  "id": "2"
}
```
- **GetCodeByID**
```bash
$ grpcurl -plaintext -d '{"id": "2"}' localhost:44000 runner.Runner/GetCodeByID

{
  "code": "package main\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"Hello world! Test ShareCode\")\n}",
  "output": "Hello world! Test ShareCode\n"
}
```
## Database migration
See Makefile.example and make sure [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) on your system.
```bash
make mig_up
# and
make mig_down
```
