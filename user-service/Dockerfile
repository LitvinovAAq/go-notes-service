FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/userapp ./main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /out/userapp ./userapp

EXPOSE 8082

ENTRYPOINT ["./userapp"]
