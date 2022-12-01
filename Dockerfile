# syntax=docker/dockerfile:1

FROM golang:1.18-bullseye

WORKDIR /app

RUN mkdir ./pkg

COPY pkg/* /app/pkg/

COPY . ./

COPY ./Caddyfile ./Caddyfile

ENV ZAPPR_POSTGRES_USER postgres
ENV ZAPPR_POSTGRES_PASSWORD postgres
ENV ZAPPR_POSTGRES_DB zappr
ENV ZAPPR_POSTGRES_HOST 0.0.0.0:5432

RUN go mod download 

RUN go build -o ./

CMD ["./Zappr"]