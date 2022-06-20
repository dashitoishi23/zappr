# syntax=docker/dockerfile:1

FROM golang:1.18-bullseye

WORKDIR /app

RUN mkdir ./pkg

COPY pkg/* /app/pkg/

COPY . ./

RUN go mod download 

RUN go build -o ./

EXPOSE 9000

CMD ["./Zappr"]