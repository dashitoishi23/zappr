# syntax=docker/dockerfile:1

FROM golang:1.18-bullseye

WORKDIR /app

RUN mkdir ./pkg

COPY pkg/* /app/pkg/

COPY . ./

COPY ./Caddyfile ./Caddyfile

RUN go mod download 

RUN go build -o ./

# CMD bash -c 'while !</dev/tcp/database/5432; do sleep 1; done; ./Zappr'