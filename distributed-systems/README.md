# Distributed systems

## K-V Store

A mapping between Twitter handle and handle & host for AT Protocol

When adding a new Twitter handle, we go ahead and fetch the decentralized identifier (DID) from the PDS you specified as the host.

### How to

```sh
🔑 (default) tw-username
  => {{ handle: bs-handle, host: bsky.social }, did: at://did:plc:bv6ggog3tya2z3vxsub7hnal }
```

```sh
🔑 (default) set tw-username bs-handle
  => {{ handle: bs-handle, host: bsky.social }, did: at://did:plc:bv6ggog3tya2z3vxsub7hnal }
```

If you need to do so, you can set the table(s) to which you want to write:

```sh
🔑 (default) pick faves outgroup
🔑 (faves,outgroup)
```

Otherwise, you'll be working with the default table.

Using `pick` is enough to create the table; we'll switch to it if it exists, or make it and then switch to it.

If you want to remove a table, switch to another table and use `drop`:

```sh
🔑 (faves) pick default
🔑 (default) drop faves
Removed faves
```

### Running

In `kv-store`, run

```sh
$ go run lb/lb.go
```

NB This will start two instances of `server/server.go`, one on port `8889` and another on `8890`.

In another tab/window/terminal, run

```sh
$ go run client/client.go
```

If you'd prefer to send data directly to the server without the client, start the server as above and then run

```sh
$ echo -n "table get my-key" | nc localhost 8888
```

### Encoding strategy

```
Handle length (2 bytes) [00 aa]
Handle (x bytes) [01 bb ...]
Host length (2 bytes) [02 cc ...]
Host (y bytes) [03 dd ...]
DID length (2 bytes) [04 ee ...]
DID (z bytes) [05 ff ...]
```

We haven't encoded field identifiers yet, but that's coming up soon.

## Development

Developed with Go v1.19+

In your editor, open the directory&mdash;eg `kv-store`&mdash;that has the `go.work` file in order to make use of Go's workspaces.

