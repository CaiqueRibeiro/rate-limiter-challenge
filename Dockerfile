FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server src/cmd/main.go

FROM scratch

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

ENTRYPOINT [ "./server" ]