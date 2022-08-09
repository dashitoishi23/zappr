# syntax=docker/dockerfile:1

FROM golang:1.18-bullseye

WORKDIR /app

RUN mkdir ./pkg

COPY pkg/* /app/pkg/

COPY . ./

COPY ./Caddyfile ./Caddyfile

ENV ZAPPR_POSTGRES_USER zappr_pg_local@zappr-local
ENV ZAPPR_POSTGRES_PASSWORD Passw0rdtechnovert
ENV ZAPPR_POSTGRES_DB zappr
ENV ZAPPR_POSTGRES_HOST zappr-local.postgres.database.azure.com

RUN go mod download 

RUN go build -o ./

# RUN apt-get update \
#  && apt-get install -y sudo

# RUN adduser --disabled-password --gecos '' docker
# RUN adduser docker sudo
# RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# USER docker

# RUN sudo apt-get update 

# RUN sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
# RUN curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
# RUN curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
# RUN sudo apt update
# RUN sudo apt install caddy

CMD ["./Zappr"]