FROM golang:1.20-alpine3.17 AS builder

# build
WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/
RUN go mod download

COPY main.go /app/
COPY cmd /app/cmd
COPY pkg /app/pkg

RUN go build -o /fetch

# deploy
FROM alpine:3.17

WORKDIR /

COPY --from=builder /fetch /usr/bin/fetch

ENTRYPOINT [ "/usr/bin/fetch" ]
