# syntax=docker/dockerfile:1
FROM golang:1.21

ARG PGHOST=host.docker.internal
ARG PGPORT=5432
ARG PGUSER=admin
ARG PGDATABASE=admin
ARG PGPASSWORD=admin

ENV PGHOST=$PGHOST
ENV PGPORT=$PGPORT
ENV PGUSER=$PGUSER
ENV PGDATABASE=$PGDATABASE
ENV PGPASSWORD=$PGPASSWORD

# Modify GOBIN before install to adhere to path used in makefile script.
ENV GOBIN=/root/go/bin/

RUN go install github.com/jackc/tern/v2@latest

WORKDIR /app

COPY /migrations ./migrations
COPY makefile ./

CMD ["make", "db-migrate"]
