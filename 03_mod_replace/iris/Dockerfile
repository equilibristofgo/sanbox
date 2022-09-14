FROM golang:1.17 AS builder
RUN apt-get update
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /go/src/demoapp
COPY . .
RUN go build -mod vendor -o /go/bin/demoapp
#RUN go install -mod vendor

FROM scratch
COPY --from=builder /go/bin/demoapp .
#COPY --from=builder /go/bin/iris .
ENTRYPOINT ["./demoapp"]
#ENTRYPOINT ["./iris"]
