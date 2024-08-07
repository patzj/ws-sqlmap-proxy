# WebSocket SQLMap Proxy

An HTTP proxy designed to enable the use of SQLMap with WebSocket connections.

## Usage

```sh
go run main.go -addr 10.0.0.10:1337 -path /message
```

## SQLMap

```sh
sqlmap -u localhost:8000 --data='{"param":"arg"}' --method=POST -H "Content-Type: application/json" -batch
```

## Limitations

This proxy only supports WebSocket messages formatted as JSON.
