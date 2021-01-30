FROM golang:1.16rc1-alpine3.12 AS build
RUN apk add --no-cache git libc-dev
WORKDIR /go/src/github.com/guicostaarantes/psi-server
COPY go.mod .
COPY go.sum .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go get github.com/99designs/gqlgen
RUN go get github.com/go-chi/chi
COPY . .
RUN chmod +x ./build.sh
RUN ./build.sh

FROM alpine:3.12
COPY --from=build /out /
EXPOSE 8082
CMD ["/main"]
