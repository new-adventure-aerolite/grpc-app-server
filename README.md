# app server

App server handles rest request from frontend, communicates with the auth server to do the validation, and forwarders them to fight server in grpc.

## Overview

![rpg-game](./images/rpg-game.png)

## installation

```sh
$ git clone https://github.com/new-adventure-aerolite/grpc-app-server.git

$ go mod tidy && go mod vendor

$ go build -o app-server ./main.go
```

help info:

```sh
$ ./app-server -h            
Usage of ./app-server:
  -addr string
        fight svc addr (default "127.0.0.1:8001")
  -auth-server-addr string
        auth svc addr (default "127.0.0.1:6666")
  -port string
        listen port (default "8000")
  -tls-cert string
        tls cert
  -tls-key string
        tls key
```

## LICENSE

[MIT](./LICENSE)
