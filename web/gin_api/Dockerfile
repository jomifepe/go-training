# Building
FROM golang:alpine AS builder

RUN mkdir -p $GOPATH/src/github.com/jomifepe/gin_api
WORKDIR $GOPATH/src/github.com/jomifepe/gin_api
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /go/bin/gin_api

# Running
FROM alpine
RUN apk add --update ca-certificates
COPY --from=builder /go/bin/gin_api /go/bin/gin_api
ENTRYPOINT ["/go/bin/gin_api"]