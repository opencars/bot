# Robot

[![Report](https://goreportcard.com/badge/github.com/shal/robot)](https://goreportcard.com/report/github.com/shal/robot)

> :blue_car: Robot will help to find a dream vehicle

## Overview

List of supported environment variables

| Env. Variable | Description              |
|---------------|--------------------------|
| `BOT_TOKEN`   | *Telegram API token*     |
| `RIA_API_KEY` | *AutoRia API token*      |
| `ALPR_URL`    | *URL to ALPR web server* |
| `HOST`        | *Host for telegram API webhhok* |
| `PORT`        | *Web server port* |

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
