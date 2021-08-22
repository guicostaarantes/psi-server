FROM golang:1.17.0-alpine3.14 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY graph/ graph
COPY modules/ modules
COPY utils/ utils
COPY main.go .
COPY build.sh .
RUN chmod +x ./build.sh
RUN ./build.sh

FROM alpine:3.14
COPY --from=build /out /
CMD ["/main"]
