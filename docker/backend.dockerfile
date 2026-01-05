# ---------- Builder ----------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app
RUN apk add --no-cache ca-certificates # git

COPY go.mod go.sum ./
COPY cmd/ cmd/
COPY internal/ internal/

# Packages
COPY pkg/core/ ./pkg/core/
COPY pkg/idp/ ./pkg/idp/
COPY pkg/logger/ ./pkg/logger/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/server/main.go


# ---------- Runtime ----------
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/app /app/app

USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
