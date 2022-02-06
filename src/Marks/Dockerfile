# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

# Set destination for COPY
WORKDIR /marks

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /marks-docker

EXPOSE 8121

CMD ["/marks-docker"]