FROM golang:1.17.4-alpine as build

WORKDIR /build
COPY . .

RUN go mod vendor

RUN go build -o /app ./cmd/server
EXPOSE 8081
CMD ["/app"]