# ikuradon-salmon

ikuradon Push Server

## Architecture

Mastodon Server (WebPush API) -> ikuradon-salmon Server -> Expo Backend Server -> FCM or APNS -> Your Terminal

Mastodon Web Push API Documents: https://docs.joinmastodon.org/methods/notifications/push/

Expo Push Notifications Documents: https://docs.expo.io/guides/push-notifications/

# API Documents

[/swagger.yaml](OpenAPI 3.0)

## build
```
go build main.go
```

## Production Run Server

1. Setting `.env:` file `BASE_URL`

2. Docker run

```
docker-compose build
docker-compose up -d
```

3. SSL Portforwarding Setting

## LICENSE

AGPL v3

Copyright (C) 2020 potproject