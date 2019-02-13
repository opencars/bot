# Robot

[![Build Status](https://travis-ci.org/shal/robot.svg?branch=master)](https://travis-ci.org/shal/robot)
[![Report](https://goreportcard.com/badge/github.com/shal/robot)](https://goreportcard.com/report/github.com/shal/robot)

> :blue_car: Robot will help to find a dream vehicle

## Overview

List of supported environment variables

| Name             | Description        |
|------------------|--------------------|
| `TELEGRAM_TOKEN` | Telegram API token |
| `AUTO_RIA_TOKEN` | AutoRia API token  |
| `HOST`           | Web server host    |
| `PORT`           | Web server port    |
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

You can run bot with command below

```sh
$ PORT=80 go run cmd/robot/main.go
```

## License

Project released under the terms of the MIT [license](./LICENSE).
