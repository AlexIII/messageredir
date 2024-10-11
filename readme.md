# Message redir: redirect messages via web hook to Telegram

This app serves a web hook and reroutes messages to a telegram bot.

Designed to work with [SMS to URL Forwarder](https://f-droid.org/en/packages/tech.bogomolov.incomingsmsgateway/) (Android app).

Other potential use cases include automated notifications from server monitoring systems, grafana, etc.

## How to use (docker)

```sh
docker build . -t messageredir
docker run -d --name messageredir-inst -e MREDIR_TG_BOT_TOKEN="YOUR_TELEGRAM_BOT_TOKEN" -v "$(pwd)/messageredir.db:/root/app/messageredir.db" -p 8089:8080 messageredir
```

## Config

| YAML Name        | Environment Variable | Type   | Description          |
|------------------|----------------------|--------|----------------------|
| `dbFileName`     | `DB_FILE_NAME`       | string | Database file name   |
| `tgBotToken`     | `TG_BOT_TOKEN`       | string | Telegram bot token   |
| `userTokenLength`| `USER_TOKEN_LENGTH`  | int    | User token length    |
| `logUserMessages`| `LOG_USER_MESSAGES`  | bool   | Log user messages    |
| `restApiPort`    | `REST_API_PORT`      | int    | REST API port        |
| `tlsCertFile`    | `TLS_CERT_FILE`      | string | TLS certificate file |
| `tlsKeyFile`     | `TLS_KEY_FILE`       | string | TLS key file         |
| `logFileName`    | `LOG_FILE_NAME`      | string | Log file name        |

## How to post message to hook

`POST http://localhost:8089/<TOKEN_THE_BOT_ISSUED_FOR_YOU>/smstourlforwarder`

Body:

```json
{
     "from": "%from%",
     "text": "%text%",
     "sentStamp": %sentStampMs%,
     "receivedStamp": %receivedStampMs%,
     "sim": "%sim%"
}
```
