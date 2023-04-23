# Distributed systems

## 01: Set it and (for)get it

### Running

In `kv-store`, run

```sh
$ go run server/server.go
```

In another tab/window/terminal, run

```sh
$ go run client/client.go
```

If you'd prefer to send data directly to the server without the client, start the server as above and then run

```sh
$ echo -n "get my-key" | nc localhost 8888
```

## Development

Developed with Go v1.19+

Run client:

```sh
$ go run client/main.go
```

