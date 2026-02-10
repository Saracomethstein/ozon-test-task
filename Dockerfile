FROM golang:1.24-bullseye AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags='-s -w' -o /out/service ./cmd/service

FROM scratch

COPY --from=builder /out/service /app/service

EXPOSE 8080
ENTRYPOINT ["/app/service"]
