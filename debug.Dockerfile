FROM golang:1.16.6-alpine3.14 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY graph/ graph
COPY modules/ modules
COPY utils/ utils
COPY main.go .
RUN apk add --no-cache build-base
RUN go build -gcflags="all=-N -l" -o /out/main main.go
RUN go get github.com/go-delve/delve/cmd/dlv
CMD ["/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue", "/out/main"]
