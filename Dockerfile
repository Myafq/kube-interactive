FROM golang:1.11 as builder
WORKDIR /tmp/
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.hash=$(date +%s)" -o /tmp/interactive
FROM alpine:latest  
COPY --from=builder /tmp/interactive /bin/interactive
ENTRYPOINT ["/bin/interactive"] 
