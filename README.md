# Bot

[![Report](https://goreportcard.com/badge/github.com/opencars/bot)](https://goreportcard.com/report/github.com/opencars/bot)

> :blue_car: Bot will help to find a dream vehicle

## Overview

![Architecture Overview](./doc/images/architecture.svg)

List of supported environment variables

| Name             | Description        |
| ---------------- | ------------------ |
| `HOST`           | Web server host    |
| `TELEGRAM_TOKEN` | Telegram API token |
| `AUTO_RIA_TOKEN` | AutoRia API token  |
| `AUTO_RIA_TOKEN` | AutoRia API token  |

## Development

Start postgres database

```sh
docker-compose up -Vd postgres
```

Migrate database

```sh
migrate -path=migrations -database "postgres://postgres:password@localhost/bot?sslmode=disable" up
```

Prerequisites:

- [Ngrok](https://ngrok.com/).

Expose your local web server for receiving http requests

```sh
ngrok http 80
```

Export ngrok host

```sh
export HOST=<host>
```

You can run the bot with command below

```sh
go run cmd/bot/main.go --port=80
```

## License

Project released under the terms of the MIT [license](./LICENSE).
