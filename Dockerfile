FROM golang:1.24-bullseye AS builder

WORKDIR /src

COPY .env ./
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM ubuntu:latest

COPY --from=builder /src/build/ozon-test-task /app/ozon-test-task
COPY --from=builder /src/.env /app/.env

EXPOSE 8080
ENTRYPOINT ["./app/ozon-test-task", "--production"]
