FROM golang:1.16-alpine AS build

WORKDIR /go/src/app

LABEL maintainer="ashanaakh@gmail.com"

RUN apk add ca-certificates git make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN export BLDDIR=/go/bin && \
    make clean && \
    make build

FROM alpine

RUN apk update && apk upgrade

WORKDIR /app

COPY --from=build /go/bin/ ./
COPY ./config ./config
COPY ./templates ./templates

EXPOSE 8080

CMD ["./bot","-port", "8080"]
