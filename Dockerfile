# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine AS build

# Setup ENV
WORKDIR /app

# Download prereqs
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy sources
COPY . .

# Build
RUN go build -o /probesvc

##
## Deploy
##
FROM alpine:latest

WORKDIR /

COPY --from=build /probesvc /probesvc
COPY ./service/ffprobe /ffprobe

RUN chmod +x /probesvc
RUN chmod +x /ffprobe

EXPOSE 8080

USER nobody:nogroup

ENTRYPOINT ["/probesvc"]
