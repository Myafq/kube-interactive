FROM golang:1.11 as builder
WORKDIR /tmp/
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -s -X main.hash=$(date +%s)" -o /tmp/interactive
FROM alpine:latest  
COPY --from=builder /tmp/interactive /bin/interactive
EXPOSE 8084
