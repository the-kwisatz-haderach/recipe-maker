# syntax=docker/dockerfile:1
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /recipe-maker ./cmd/server.go

FROM alpine:latest

RUN apk --no-cache add curl

RUN addgroup -S nonroot && adduser -S user -G nonroot
USER user

COPY --from=build-stage /recipe-maker /recipe-maker

EXPOSE 8080

ENV PORT 8080

# set hostname to localhost
ENV HOSTNAME "0.0.0.0"

HEALTHCHECK --interval=30s --timeout=3s --start-interval=5s \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["/recipe-maker"]
