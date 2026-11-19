# vector-hugo
In-memory data structure server in Go (Redis-style)

vector-hugo speaks the [RESP](https://redis.io/docs/latest/develop/reference/protocol-spec/) protocol. At the moment it supports **lists**; sets and vector sets are planned.

> **Note:** standard Redis clients such as `redis-cli` or `go-redis` are **not supported yet**. The server only accepts lowercase command names, serves one command per connection, and accepts only values made of letters, digits and `.`. For now, send raw RESP requests yourself (see below).

## Quick start

```bash
go build -o vector-hugo ./cmd
./vector-hugo -path config.yaml
```

Send a raw RESP request, for example with `nc`:

```bash
# rpush mylist hello
printf '*3\r\n$5\r\nrpush\r\n$6\r\nmylist\r\n$5\r\nhello\r\n' | nc 127.0.0.1 6379
# :1

# lrange mylist 0 1
printf '*4\r\n$6\r\nlrange\r\n$6\r\nmylist\r\n$1\r\n0\r\n$1\r\n1\r\n' | nc 127.0.0.1 6379
# *1
# $5
# hello
```

## Configuration

The server reads a YAML file passed with the `-path` flag.

```yaml
host: 127.0.0.1        # IP address to bind (required)
port: 6379             # 1-65535 (required)

timeout_unit: ss       # "ss" (seconds) or "ms" (milliseconds), default "ss"
timeout_value: 10      # per-connection timeout, default 10
max_bytes: 16000       # max bytes accepted per request, default 16000

database_options:      # optional, per-list memory settings
  mylist:
    max_memory: 20480       # arena size in bytes, default 20480
    page_replacer: true     # enable LRU eviction
    lru_pool_size: 16       # LRU candidates, default 16 (requires page_replacer)
    memory_threshold: 15360 # eviction starts above this, must be < max_memory
```

Lists without an entry in `database_options` use the built-in `standard` profile (20 KB arena, LRU enabled, pool of 16, threshold at 75%).

## Architecture

```
cmd/main.go            entry point: reads config, starts the server
internal/conf          YAML config loading and validation
internal/server        TCP server, request reading and command dispatch
internal/proto         RESP lexer, parser, binder and response builder
internal/router        maps a list name to its memory arena
internal/functions     command implementations (LPUSH, LPOP, ...)
internal/memory        mmap-backed arena and LRU replacer
```

### Server

`server.VectorHugo` listens on TCP and handles each connection in its own goroutine:

1. **Read** – data is read in 1 KB chunks into a buffer. The first line (`*<n>`) gives the number of array elements, so the server knows the request is complete after `1 + 2*n` CRLF-terminated lines (see `proto/resp_waiter.go`). Requests larger than `max_bytes` are dropped.
2. **Parse & bind** – the buffer goes through the lexer → parser → binder pipeline (below). Any error is sent back as `-ERR <stage> <message>`.
3. **Execute** – the router returns the arena for the list named in the first argument (creating it only for `LPUSH`/`RPUSH`), the matching function runs, and the result is encoded as a RESP reply.

Each connection serves one command and is then closed.

### Parser (RESP pipeline)

Requests must be RESP arrays of bulk strings, e.g. `RPUSH mylist hello`:

```
*3\r\n$5\r\nRPUSH\r\n$6\r\nmylist\r\n$5\r\nhello\r\n
```

The `internal/proto` package processes them in three stages:

| Stage  | File             | What it does |
|--------|------------------|--------------|
| Lexer  | `resp_lexer.go`  | Splits the input on CRLF and produces `(token, literal)` pairs: `*`/`$` prefixes followed by a `DIGIT`, command names (`lpush`, `lrange`, ...) as their own tokens, everything else as `TEXT`. Literals may contain only letters, digits and `.`. |
| Parser | `resp_parser.go` | Builds an `ArrayAST`: checks the `*<len>` header, then parses exactly `len` bulk strings, verifying that each `$<len>` matches the length of its text. |
| Binder | `resp_binder.go` | Turns the AST into an `ExecutableCommand{CommandName, Args}`, checking that the first element is a known command and that enough arguments were given (arity table). |

Replies are built by `resp_builder.go`: simple strings (`+OK`), errors (`-ERR ...`), integers (`:n`), bulk strings (`$n`), null (`$-1`) and arrays (`*n`).

> Command names are matched in lowercase (`lpush`, not `LPUSH`).

### Memory

Every list lives in its own `memory.Arena`: a fixed-size buffer allocated with `mmap`. Entries are stored contiguously as `|length (1 byte)|payload|`, so a single value can be at most 255 bytes.

When memory usage passes `memory_threshold`, the arena fills an LRU pool with a random sample of entries and evicts the least recently used one before each new insert. Reads via `LINDEX`/`LRANGE` refresh an entry's position in the LRU. If the arena is still full, the write fails with `Insufficient space for storing new data`.

## Commands

| Command | Syntax | Reply | Notes |
|---------|--------|-------|-------|
| `LPUSH`  | `lpush key value [value ...]` | Integer: number of values pushed | Creates the list if missing. Inserts at the head. |
| `RPUSH`  | `rpush key value [value ...]` | Integer: number of values pushed | Creates the list if missing. Inserts at the tail. |
| `LPUSHX` | `lpushx key value [value ...]` | Integer | Like `LPUSH`, but only if the list exists; otherwise returns `0`. |
| `LPOP`   | `lpop key` | Bulk string, or null if the list is missing | Removes and returns the head. |
| `RPOP`   | `rpop key` | Bulk string, or null if the list is missing | Removes and returns the tail. |
| `LLEN`   | `llen key` | Integer | Number of elements (`0` if the list is missing). |
| `LINDEX` | `lindex key index` | Bulk string | Element at `index` (0-based, no negative indexes). |
| `LRANGE` | `lrange key start end` | Array | Elements from `start` to `end` inclusive (no negative indexes). |

`BLPOP` and the set commands (`SADD`, `SREM`, `SMEMBERS`, `SISMEMBER`, `SCARD`, `SPOP`) are recognised by the lexer but not implemented yet.

## Testing

```bash
go test ./...
```
