FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /gh-proxy-go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /gh-proxy-go /usr/local/bin/gh-proxy-go
EXPOSE 8080
ENTRYPOINT ["gh-proxy-go"]
