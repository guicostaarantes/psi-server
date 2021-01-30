FROM golang:1.16rc1-alpine3.12 AS build
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN chmod +x ./build.sh
RUN ./build.sh

FROM alpine:3.12
COPY --from=build /out /
EXPOSE 8082
CMD ["/main"]
