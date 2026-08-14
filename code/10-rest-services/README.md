# Topic 10 sample code — REST Services

A single `items` CRUD REST API, one Go module, structured as
Handler -> Service -> Repository (the same dependency-injection pattern
from Topics 8-9, built out completely for the first time).

## Layout

```
cmd/server/main.go        composition root — wires everything together
internal/domain           the Item struct, shared by every layer
internal/repository       storage: an interface + an in-memory impl
internal/service          validation and business rules, no HTTP
internal/handler          HTTP only: decode, call the service, respond
```

## Run it directly

```
go run ./cmd/server
```

```
curl localhost:8080/ping
curl -X POST localhost:8080/items -d '{"name":"widget","quantity":5}'
curl localhost:8080/items
curl localhost:8080/items/1
curl -X PUT localhost:8080/items/1 -d '{"name":"widget","quantity":9}'
curl -X DELETE localhost:8080/items/1
```

Every request gets logged (method, path, duration) by the `Logging`
middleware in `internal/handler/middleware.go` — the Decorator pattern
from Topic 9, wrapping the whole `mux`.

## Run it with docker-compose

```
docker-compose up --build
```

Brings up the API (built from the `Dockerfile` in this directory) on
`localhost:8080`, plus a Postgres container the app doesn't use yet —
included to show what wiring in a real dependency looks like.

```
curl localhost:8080/ping
```

```
docker-compose down
```

**If you're on Windows with WSL:** run every `docker-compose` command
*from inside your WSL shell*, not PowerShell or CMD. Docker Desktop's
WSL2 backend already runs containers inside a Linux VM, so working from
WSL avoids an extra translation layer and gives correct file permissions.
Keep this project on the Linux filesystem (e.g. `~/projects/rest-service`)
rather than under `/mnt/c/...` — bind-mounting a Windows drive into the
WSL2 VM is noticeably slower and can make file-watching tools miss
changes. See the slides for the full WSL section.

## Formatting

```
gofmt -l .    # lists files that don't match gofmt's formatting
gofmt -w .    # rewrites them in place
```
