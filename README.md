# OpenCars Bot

[![Report](https://goreportcard.com/badge/github.com/shal/opencars-bot)](https://goreportcard.com/report/github.com/shal/opencars-bot)

> :blue_car: OpenCars Bot will help to find a dream vehicle

## Overview

![Architecture Overview](./doc/images/architecture.svg)

List of supported environment variables

| Name             | Description        |
|------------------|--------------------|
| `HOST`           | Web server host    |
| `PORT`           | Web server port    |
| `TELEGRAM_TOKEN` | Telegram API token |
| `AUTO_RIA_TOKEN` | AutoRia API token  |
| `RECOGNIZER_URL` | Plates recognizer  |
| `DATA_PATH`      | Path to data file  |

## Development

Prerequisites:
- [Ngrok](https://ngrok.com/).

Expose your local web server for receiving http requests

```sh
$ ngrok http 80
```

Export ngrok host

```sh
$ export HOST=<host>
```

You can run the bot with command below

```sh
$ PORT=80 go run cmd/bot/bot.go
```

## License

Project released under the terms of the MIT [license](./LICENSE).
