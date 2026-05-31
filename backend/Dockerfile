# backend/Dockerfile
FROM golang:1.26-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o server .
EXPOSE 8090
CMD ["./server"]
