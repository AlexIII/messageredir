# Build stage
FROM golang:1.23.2-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o messageredir ./cmd/messageredir

# Runtime stage
FROM alpine:latest

WORKDIR /root/app
COPY --from=build /app/messageredir .
CMD ["./messageredir"]
