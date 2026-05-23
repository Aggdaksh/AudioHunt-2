FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/audiohunt-api ./cmd/api

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends ffmpeg ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /bin/audiohunt-api /app/audiohunt-api
COPY migrations /app/migrations

RUN mkdir -p /app/uploads/queries

ENV PORT=8080

EXPOSE 8080

CMD ["/app/audiohunt-api"]
