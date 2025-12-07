FROM golang:1.25.1-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /out
RUN CGO_ENABLED=0 go build -o /out/userapp ./main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget
COPY --from=builder /out/userapp /usr/local/bin/userapp
RUN adduser -D -g '' appuser
USER appuser
WORKDIR /home/appuser
EXPOSE 8082
ENV PORT=8082
HEALTHCHECK --interval=10s --timeout=2s --retries=3 \
  CMD wget -qO- http://localhost:8082/health || exit 1
ENTRYPOINT ["userapp"]
