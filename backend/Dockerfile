# backend/Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /src
# Кэш модулей: зависимости ставятся до исходников
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/notes .

# backend/Dockerfile (фрагмент runtime-стадии)
FROM alpine:3.20 AS runtime
RUN adduser -D -H app
WORKDIR /app
COPY --from=builder /out/notes ./notes
COPY healthcheck.sh ./healthcheck.sh
RUN chmod +x /app/healthcheck.sh && chown app:app /app /app/healthcheck.sh
USER app
EXPOSE 8000
HEALTHCHECK --interval=10s --timeout=3s --retries=5 CMD ["/app/healthcheck.sh"]
CMD ["/app/notes"]
