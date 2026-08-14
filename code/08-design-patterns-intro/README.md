# Topic 8 sample code — Intro to Design Patterns

One module, two demos. Run everything from this directory
(`code/08-design-patterns-intro/`, or wherever this folder lives).

## `cmd/singletondemo` + `internal/config`

The Singleton pattern, done the idiomatic Go way: a package-level
`sync.Once` guards a lazily-built `*Config`.

```
go run ./cmd/singletondemo
```

Confirms two things:

1. Two separate calls to `config.GetConfig()` return the exact same
   pointer (`a == b`).
2. Fifty goroutines calling `GetConfig()` concurrently still only trigger
   the underlying loader **once** — `sync.Once` serializes the first call
   and lets everyone else through with the already-built result.

Try it: comment out `once.Do(...)` in `internal/config/config.go` and
replace it with the naive `if instance == nil { ... }` version from the
slides. Run `go run -race ./cmd/singletondemo` (if you have a Go
toolchain with the race detector available) and watch it flag the
concurrent read/write.

## `cmd/layersdemo` + `internal/user`

The Handler → Service → Repository layering, minus the HTTP handler
itself (that's Topic 10). `internal/user` defines:

- `Repository` — an interface, `FindByID(id string) (*User, error)`
- `Service` — depends on `Repository`, not a concrete type
- `NewService(repo Repository) *Service` — the constructor injection
  point (the "composition root" pattern in miniature)
- `InMemoryRepository` — a real, if minimal, `Repository` implementation

```
go run ./cmd/layersdemo
```

Builds `Service` twice: once with the real `InMemoryRepository`, once
with a `fakeRepository` defined right in `main.go` that always returns
canned data. **`Service`'s code doesn't change between the two** — only
what gets passed into `NewService` at the call site. That's the whole
point of dependency injection in Go: no framework, just an interface and
a constructor argument.

## Formatting

```
gofmt -l .    # lists files that don't match gofmt's formatting
gofmt -w .    # rewrites them in place
```
