FROM golang:1.24-alpine AS builder

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=builder /app/bin/main .
COPY --from=builder /app/pkg ./pkg
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 3000

CMD ["./main"]