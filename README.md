<div align="center">
<h1>Runner Service</h1>
</div>

## About
The Runner service allows you to ***execute*** and ***share*** Go code snippets.

Additionally, it integrates an isAdmin endpoint from [sso](https://github.com/paniccaaa/sso) repository.

Protobuf contract: [runner](https://github.com/paniccaaa/protos/blob/main/proto/runner/runner.proto) 

## Stack

- **gRPC** / gRPC-gateway (for [understanding](https://github.com/grpc-ecosystem/grpc-gateway))
- **PostgreSQL**
- **cleanenv**
- **slog**

## Installation

Before running the application, make sure Go are installed on your system.

1. **Clone** the repository
```bash
git clone https://github.com/paniccaaa/runner.git
cd runner
```
2. **Create** the .env, local.yaml and Makefile (see .example in repository)
3. **Run** the app locally
```bash
make run 
```

## Endpoints

For detailed API documentation, you can view the Swagger documentation. Copy the contents of [runner.swagger.json](https://github.com/paniccaaa/protos/blob/main/gen/golang/runner/runner.swagger.json) and paste it into [Swagger Editor](https://editor-next.swagger.io/) to explore and interact with the API.

### REST (gRPC-gateway)
- **POST /run**: Execute a Go code snippet.
- **POST /share**: Share a Go code snippet.
- **GET /shared/{id}**: Retrieve a shared Go code snippet by id.
  
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
