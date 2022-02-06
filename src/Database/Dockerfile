# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

# Set destination for COPY
WORKDIR /database

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /database-docker

EXPOSE 8129

CMD ["/database-docker"]