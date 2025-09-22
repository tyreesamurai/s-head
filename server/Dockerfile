# Use a Go image that satisfies the toolchain requirement
FROM golang:1.24.1-alpine AS build
WORKDIR /src

# Cache modules (optional but faster with BuildKit)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
  CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

FROM alpine:3.20
WORKDIR /app
# trust roots for HTTPS calls
RUN apk add --no-cache ca-certificates && adduser -D -H appuser
COPY --from=build /src/app .
USER appuser
EXPOSE 8080
CMD ["./app"]
