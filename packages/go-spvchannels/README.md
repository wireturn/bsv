# SPV Channels Client in go

This repository contains source code implementing SPV Channels Client in golang for the [SPV Channels Server](https://github.com/bitcoin-sv/spvchannels-reference)

## Setup SPV Channels server

Following the tutorial in the [SPV Channels Server](https://github.com/bitcoin-sv/spvchannels-reference), we first create the certificate using `openssl`:
```
terminal $> openssl req -x509 -out localhost.crt -keyout localhost.key -newkey rsa:2048 -nodes -sha256 -subj '/CN=localhost' -extensions EXT -config <( printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
terminal $> openssl pkcs12 -export -out devkey.pfx -inkey localhost.key -in localhost.crt # use devkey as password
```

That will create the `devkey.pfx` with password `devkey`. We then write a `docker-compose.yml` file following the tutorial. To run the [SPV Channels Server](https://github.com/bitcoin-sv/spvchannels-reference) :
```
docker-compose up -d
```

We then need to create a SPV Channels account on the server
```
docker exec spvchannels ./SPVChannels.API.Rest -createaccount spvchannels_dev dev dev
```

## Usage with swagger

The [SPV Channels Server](https://github.com/bitcoin-sv/spvchannels-reference) run by `docker-compose.yml` listen on `localhost:5010`. We can start playing with the endpoints using swagger, i.e in browser, open `https://localhost:5010/swagger/index.html`

From this page, there are a link `/swagger/v1/swagger.json` to export swagger file

## Usage with Postman

Interacting with browser might have some difficulty related to adding certificate to the system. It might be easier to use Postman to interact as Postman has a easy possibility to disable SSL certificate check to ease development propose.

From Postman, import the file `devconfig/postman.json` and set the environment config as follow

| VARIABLE    | INITIAL VALUE  |
| ----------- | -------------- |
| URL_PORT    | localhost:5010 |
| ACCOUNT     | 1              |
| USERNAME    | dev            |
| PASSWORD    | dev            |

These environment variable are used as _template_ to populate values in the `postman.json` file. There are a few more environment variable to define (look into the json file) that will depend to the endpoint and value created during the experience:

| VARIABLE     | INITIAL VALUE   |
| ------------ | --------------- |
| CHANNEL_ID   | .. to define .. |
| TOKEN_ID     | .. to define .. |
| TOKEN_VALUE  | .. to define .. |
| MSG_SEQUENCE | .. to define .. |
| NOTIFY_TOKEN | .. to define .. |

## Run tests

Run unit test
```
go clean -testcache && go test -v ./...
```

To run integration tests, make sure you have `docker-compose up -d` on your local machine, then run
```
go clean -testcache && go test -v -tags=integration ./...
```