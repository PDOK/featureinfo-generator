FROM golang:1.24 AS  build-env

WORKDIR /go/src/server

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go test -short && \
    go build -a -o /featureinfo-generator cmd/main.go

# not distroless, sometimes a shell is used with this image
FROM alpine:3 AS service

RUN apk update && apk upgrade && apk add \
      bash \
      ca-certificates \
      && \
    rm -rf /var/cache/apk/* && \
    update-ca-certificates

WORKDIR /
ENV PATH=${PATH}:/
COPY /data /data
COPY --from=build-env  /featureinfo-generator  /

ENTRYPOINT ["/featureinfo-generator"]
