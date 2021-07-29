FROM golang:1.16.6-alpine3.14 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY tasks/ tasks
COPY main.go .
COPY build.sh .
RUN chmod +x ./build.sh
RUN ./build.sh

FROM alpine:3.14
COPY --from=build /out /
CMD ["/main"]