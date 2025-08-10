FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY .. /app
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -installsuffix cgo -o app.bin cmd/rest/main.go

FROM alpine:edge
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/app.bin ./
CMD ["./app.bin"]