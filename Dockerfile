# syntax=docker/dockerfile:1
# Build from the bridgr-api directory (default compose context).
FROM golang:1.24-alpine AS build
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY . .
RUN go mod download \
  && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/bridgr-api ./cmd/api \
  && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/bridgr-worker ./cmd/worker

FROM alpine:3.20
RUN apk add --no-cache ca-certificates curl
COPY --from=build /out/bridgr-api /usr/local/bin/bridgr-api
COPY --from=build /out/bridgr-worker /usr/local/bin/bridgr-worker
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/bridgr-api"]
