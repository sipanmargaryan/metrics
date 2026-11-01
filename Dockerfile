# syntax=docker/dockerfile:1
FROM golang:1.27.3-alpine AS base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS test
RUN go test ./...

FROM base AS build
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/metrics ./cmd/api
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.22 AS runtime
WORKDIR /app
COPY --from=build /out/metrics /app/metrics
COPY --from=build /out/migrate /app/migrate
USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/app/metrics"]
