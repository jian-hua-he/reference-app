FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ vendor/

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY proto/ proto/

RUN GOTOOLCHAIN=auto go build -o /server ./cmd/server

FROM alpine:3.21

COPY --from=builder /server /server

ENTRYPOINT ["/server"]
