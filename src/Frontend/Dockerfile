# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

# Set destination for COPY
WORKDIR /frontend

COPY go.mod ./
COPY go.sum ./
COPY ./html ./html
RUN go mod download

COPY *.go ./

RUN go build -o /frontend-docker

EXPOSE 8120

CMD ["/frontend-docker"]