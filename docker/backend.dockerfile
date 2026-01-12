# ---------- Builder ----------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app
RUN apk add --no-cache ca-certificates # git

COPY go.mod go.sum ./
COPY ./backend/ ./backend/
COPY ./libs/core/ ./libs/core/
COPY ./libs/idp/ ./libs/idp/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o backend ./backend/main.go


# ---------- Runtime ----------
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/backend /app/backend

USER nonroot:nonroot
ENTRYPOINT ["/app/backend"]
