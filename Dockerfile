# syntax=docker/dockerfile:1

FROM golang:1.18-bullseye

WORKDIR /app

RUN mkdir ./pkg

COPY pkg/* /app/pkg/

COPY . ./

ENV ZAPPR_POSTGRES_USER zappr_pg_local@zappr-local
ENV ZAPPR_POSTGRES_PASSWORD Passw0rdtechnovert
ENV ZAPPR_POSTGRES_DB zappr
ENV ZAPPR_POSTGRES_HOST zappr-local.postgres.database.azure.com
ENV REDIS_REPLICATION_MODE slave
ENV REDIS_MASTER_HOST redis-master
ENV REDIS_MASTER_PORT_NUMBER 6379
ENV REDIS_MASTER_PASSWORD Passw0rdRedisMaster
ENV REDIS_PASSWORD Passw0rdRedisReplica
ENV REDIS_REPLICATION_MODE master
ENV REDIS_PASSWORD Passw0rdRedisMaster

RUN go mod download 

RUN go build -o ./

CMD ["./Zappr"]