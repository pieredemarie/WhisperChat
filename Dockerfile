FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app 
COPY go.mod go.sum ./ 
RUN go mod download 
COPY . . 

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o /app/main ./cmd

FROM alpine:latest 
RUN apk --no-cache add ca-certificates tzdata 
WORKDIR /app 
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
