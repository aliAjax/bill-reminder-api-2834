FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/bill-reminder-api ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app
RUN mkdir -p /app/data && chown -R app:app /app

COPY --from=build --chown=app:app /out/bill-reminder-api /usr/local/bin/bill-reminder-api

ENV PORT=8080 \
    DATA_FILE=data/bills.json

EXPOSE 8080
USER app

ENTRYPOINT ["bill-reminder-api"]
