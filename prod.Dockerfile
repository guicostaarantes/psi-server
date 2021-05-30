FROM golang:1.16rc1-alpine3.12 AS build
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

FROM alpine:3.12
COPY --from=build /out /
CMD ["/main"]
