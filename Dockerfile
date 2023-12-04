# syntax=docker/dockerfile:1
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /recipe-maker ./cmd/server.go

FROM golang:1.21-alpine

RUN addgroup -S nonroot && adduser -S user -G nonroot
USER user

COPY --from=build-stage /recipe-maker /recipe-maker

EXPOSE 8080

CMD ["/recipe-maker"]