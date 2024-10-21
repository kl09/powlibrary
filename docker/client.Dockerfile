FROM golang:1.22-alpine AS base

ENV GOPATH=/opt/service/.go

WORKDIR /opt/service/
COPY . .

RUN go mod download
RUN go build -o /bin/main ./cmd/client

ENTRYPOINT ["/bin/main"]
