FROM golang:1.17 AS builder

WORKDIR $GOPATH/src/github.com/kedahttp/http-add-on

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPROXY="https://proxy.golang.org"
RUN go build -o /bin/megaboom .

FROM ubuntu:14.04

WORKDIR /
COPY --from=builder /bin/megaboom .

EXPOSE 8080

CMD ./megaboom
