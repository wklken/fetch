FROM golang:1.20-alpine3.17 AS builder

# build
WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/
RUN go mod download

COPY main.go /app/
COPY cmd /app/cmd
COPY pkg /app/pkg

RUN go build -o /httptest

# deploy
FROM alpine:3.17

WORKDIR /

COPY --from=builder /httptest /usr/bin/httptest

ENTRYPOINT [ "/usr/bin/httptest" ]
