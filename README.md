# ikuradon-salmon

![Test](https://github.com/potproject/ikuradon-salmon/workflows/Test/badge.svg) ![Go Report Card](https://goreportcard.com/badge/github.com/potproject/ikuradon-salmon)

[ikuradon](https://github.com/potproject/ikuradon) Push Server

## Architecture

Mastodon Server (WebPush API) -> ikuradon-salmon Server -> Expo Backend Server -> FCM or APNS -> Your Terminal

Mastodon Web Push API Documents: https://docs.joinmastodon.org/methods/notifications/push/

Expo Push Notifications Documents: https://docs.expo.io/guides/push-notifications/

## Security

ikuradon-salmonは、暗号化されたコンテンツを複合し、Expo サーバーに送信する仕様になっています。
そのため、サーバの管理者はコンテンツの内容を知ることが可能となります。これは、Expo Backendの仕様によるものです。
また、アクセストークンをPushサーバに送信することになるため、信頼できないikuradon-salmonサーバにアクセストークンを送ることはセキュリティ上のリスクが伴います。

ikuradon-salmon server decrypts the encrypted content and sends it to the Expo server.
Therefore, the server administrator can know the details of the content. This is by Expo Backend specification.
Sending the access token to an untrusted ikuradon-salmon server is a security risk because it will send the access token to the Push server.

## API Documents

[OpenAPI 3.0](/swagger.yaml)

## build

```
go build main.go
```

## Production Run Server (Docker)

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
