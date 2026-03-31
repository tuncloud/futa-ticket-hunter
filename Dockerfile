# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN set -eux; \
        GOOS_VALUE="${TARGETOS:-$(go env GOOS)}"; \
        GOARCH_VALUE="${TARGETARCH:-$(go env GOARCH)}"; \
        CGO_ENABLED=0 GOOS="$GOOS_VALUE" GOARCH="$GOARCH_VALUE" \
            go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api; \
        CGO_ENABLED=0 GOOS="$GOOS_VALUE" GOARCH="$GOARCH_VALUE" \
            go build -trimpath -ldflags="-s -w" -o /out/worker ./cmd/worker

FROM alpine:3.20 AS runtime
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY config.yaml /app/config.yaml

FROM runtime AS api
COPY --from=builder /out/api /app/api
ENTRYPOINT ["/app/api", "/app/config.yaml"]

FROM runtime AS worker
COPY --from=builder /out/worker /app/worker
ENTRYPOINT ["/app/worker", "/app/config.yaml"]
