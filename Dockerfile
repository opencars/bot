FROM golang:alpine AS build

ENV GO111MODULE=on

WORKDIR /go/src/app

LABEL maintainer="github@shanaakh.pro"

RUN apk add bash git gcc g++ libc-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /go/bin/robot ./cmd/robot/main.go

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=build /go/bin/robot ./robot
COPY templates/ templates/

ENTRYPOINT ["./robot"]