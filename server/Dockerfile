# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
COPY server.go go.mod ./
#COPY go.mod go.mod

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o server server.go

WORKDIR /
FROM alpine:latest
COPY --from=builder /workspace/server ./

EXPOSE 8080
CMD ["./server"]
