# Matches users/users/docker/openapi/Dockerfile (toolchain versions).
FROM golang:1.24.12-alpine3.23

WORKDIR /workspace

RUN apk add --no-cache nodejs npm

RUN npm install -g @redocly/cli
# Pin major v1 generator to match existing *.gen.go headers (users image used @latest).
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.16.3

COPY docker/openapi/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["generate"]
