FROM golang:1.17.0-alpine3.14 AS build
WORKDIR /app
RUN apk add --no-cache build-base
RUN go get github.com/go-delve/delve/cmd/dlv
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY graph/ graph
COPY modules/ modules
COPY utils/ utils
COPY main.go .
RUN go build -gcflags="all=-N -l" -o /out/main main.go
CMD ["/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue", "/out/main"]
