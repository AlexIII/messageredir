# Message redir: Web Hook to Telegram

This app serves a web hook and reroutes messages to a telegram bot.

It is designed to work with [SMS to URL Forwarder](https://f-droid.org/en/packages/tech.bogomolov.incomingsmsgateway/) (Android app).

Other potential use cases include automated notifications from server monitoring systems, grafana, etc.

## How to use (docker)

```
docker build . -t messageredir
docker run -d --name messageredir-inst -e MREDIR_TG_BOT_TOKEN="YOUR_TELEGRAM_BOT_TOKEN" -v "$(pwd)/messageredir.db:/root/app/messageredir.db" -p 8089:8080 messageredir
```
