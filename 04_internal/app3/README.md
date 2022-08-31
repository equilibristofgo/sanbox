# Go Microservice Template

## Start http server
```bash
go run apps/http/server.go
```

## Http swagger generation
```bash
swag init -g apps/http/server.go --parseDependency
```

## Docker build
```bash
docker build -t server-template:v1 . --no-cache
```