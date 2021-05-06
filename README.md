# VerifiableCredential Service

A HTTP service that wraps [vc-sdk](https://github.com/medibloc/vc-sdk)


## Features

- Sign credentials

### TODO
- Verify credentials
- Sign presentations
- Verify presentations


## Building and Running

### Without Docker
```bash
go build ./...

PORT=8888 \
go run cmd/main.go
```

### With Docker
```bash
docker build -t vc-service .
docker run -e PORT=8888 -p 8888:8888 vc-service
```

### Environment Variables

|Env Var|Desc|Default|
|-------|----|-------|
|DEBUG|Turn on debug logs|false|
|PORT|HTTP port||
|READ_TIMEOUT|HTTP read timeout|10s|
|WRITE_TIMEOUT|HTTP write timeout|10s|
|IDLE_TIMEOUT|HTTP idle timeout|60s|


## REST API

### Sign credential

```bash
curl -X POST localhost:8888/credentials/issue -d @samples/issue-credential.json
```

The sample JSON file is in the [samples/](samples) directory.

### Responses

|Status Code|Description|
|-----------|-----------|
|201|Credential successfully issued|
|400|Invalid input|
|500|Internal error|
