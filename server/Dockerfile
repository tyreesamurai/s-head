FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# adjust ./cmd/server if your main package path differs
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /src/app .
EXPOSE 8080
# optional: non-root user
RUN adduser -D -H appuser
USER appuser
CMD ["./app"]
