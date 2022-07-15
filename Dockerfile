# syntax=docker/dockerfile:1

FROM golang:1.18-bullseye

WORKDIR /app

RUN mkdir ./pkg

COPY pkg/* /app/pkg/

COPY . ./

# ENV POSTGRES_USER="zappr_pg_local@zappr-local"
# ENV POSTGRES_PASSWORD="p@$$w0rd"
# ENV POSTGRES_DB="zappr"
# ENV POSTGRES_HOST="zappr-local.postgres.database.azure.com"

RUN go mod download 

RUN go build -o ./

EXPOSE 9000

CMD ["./Zappr"]