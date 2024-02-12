# Go Project Template
## Usage
### Create migration
```shell
go run main.go migrate create -c config.yaml --path="migrations" --name="some name"
```
### Create apply all migrations
```shell
go run main.go migrate up -c config.yaml --path="migrations"
```
### Rollback single migration
```shell
go run main.go migrate down -c config.yaml --path="migrations"
```

### Start gRPC server with HTTP REST API
```shell
go run main.go serve -c config.yaml
```
