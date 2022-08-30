# syntax=docker/dockerfile:1

ARG GO_VERSION=1.18
ARG PORT=4000

# Build
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ../../.. ./

RUN go build apps/http/server.go

#CMD [ "./server" ]

# Deploy
FROM alpine:3.15 AS final

WORKDIR /

COPY --from=build app/server /server
COPY --from=build app/configs /configs

EXPOSE ${PORT}

CMD ["./server"]