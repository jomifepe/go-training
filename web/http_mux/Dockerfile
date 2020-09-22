FROM golang:1.15-alpine
RUN apk update && apk add --no-cache git inotify-tools

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
COPY go.mod .
COPY go.sum .

RUN go mod download
# RUN go build -o main .
RUN go get github.com/pilu/fresh

# ENTRYPOINT sh run-server.sh
ENTRYPOINT ["fresh"]
# CMD ["/app/main"]