---
title: "**REST Services**"
sub_title: "Go Programming — Topic 10"
author: Kevin Cunningham
---

## Opening scenario

Your Java team's Spring controller throws a `NotFoundException`. A global
`@ExceptionHandler` catches it somewhere else entirely and turns it into
a `404` automatically — the controller method itself never mentions a
status code.

**Type in chat: what does the equivalent Go handler have to do
explicitly, and why might that actually be a good thing?**

We'll come back to this once you've written a few handlers yourself.

<!--
speaker_note: |
  Let guesses land - "check the error itself," "write an if statement,"
  "return the status by hand" are all on the right track. Don't resolve
  it yet. The "why might that be good" half usually gets fewer answers
  up front - that's fine, it's the payoff slide later.
-->

<!-- end_slide -->

## What a REST service is

A REST service exposes resources over HTTP: `GET`, `POST`, `PUT`,
`PATCH`, `DELETE`, with status codes that mean something.

Go's standard library — `net/http` plus `encoding/json` — is enough to
build one. No framework required.

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
})

http.ListenAndServe(":8080", mux)
```

<!-- pause -->

That's a complete, working endpoint. No project generator, no
`app.get(...)` wrapper — just the standard library.

<!-- end_slide -->

## No batteries-included framework

Spring Boot, Django, Express all hand you routing, serialization,
validation, and a project structure on day one.

`net/http` hands you a `ServeMux` and lets you do the rest.

<!-- incremental_lists: true -->

- Full frameworks exist — Gin, Echo, Fiber — and are genuinely fine choices
- Idiomatic Go culture leans toward fewer dependencies
- The standard library is now good enough for most real services

<!-- incremental_lists: false -->

**This is a deliberate trade-off, not a gap.** Nobody's waiting for Go to
"get a real framework" the way that sentence might land from a Java or
Node background.

<!--
speaker_note: |
  Expect "so what's the Go equivalent of Spring Boot?" as a genuine
  question, not a gotcha, from anyone with an enterprise Java background.
  The honest answer is "there isn't one you're expected to reach for by
  default" - sit with that rather than rushing to name Gin as a patch.
-->

<!-- end_slide -->

## Go 1.22 changed the router

Before February 2024, `net/http`'s router only matched exact paths — no
method matching, no path parameters. Real projects needed a third-party
router (`gorilla/mux`, `chi`) just to write `GET /items/{id}`.

<!-- pause -->

**Go 1.22 added both, directly to the standard library:**

```go
mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    // ...
})
```

**Flag this explicitly:** older tutorials and blog posts will show you a
heavier routing setup than you now strictly need. Check `go version`
before trusting one.

<!-- end_slide -->

## JSON runs on struct tags

Java's Jackson uses `@JsonProperty`. Python's Pydantic uses class-level
config. Go attaches plain string metadata directly to fields:

```go
type Item struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Quantity int    `json:"quantity"`
}
```

<!-- pause -->

A struct tag is just a string sitting on the field, doing nothing on its
own. `encoding/json` reads it via **reflection, at runtime**, to decide
which JSON key maps to which field.

<!-- end_slide -->

## Two shapes of the same job

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Streaming — HTTP already gave you a reader/writer**

```go
json.NewEncoder(w).Encode(item)

json.NewDecoder(r.Body).
    Decode(&item)
```

Reads or writes straight to/from an
`io.Reader`/`io.Writer` — no
intermediate `[]byte` in your hands.

<!-- column: 1 -->

**Marshal/Unmarshal — everywhere else**

```go
data, err := json.Marshal(item)

var item Item
err := json.Unmarshal(data, &item)
```

The general-purpose form: a Go
value in, a `[]byte` out (or the
reverse) — same as Python's
`json.dumps`/`loads` or JS's
`JSON.stringify`/`parse`.

<!-- reset_layout -->

<!-- pause -->

**Same struct tags, same reflection underneath — the only difference is
what's on the other end.** Reach for `Marshal`/`Unmarshal` the moment
there's no `io.Reader`/`Writer` already in scope: logging a struct,
writing a fixture to a test, building a request body in a test file
(exactly what the pre-written tests in this lab's `handler_test.go` do).

<!--
speaker_note: |
  Worth naming explicitly that they've already seen json.Marshal once
  without it being explained - it's sitting unremarked in this lab's
  own handler_test.go, building the JSON body a test POSTs in. Point
  that out directly rather than assuming they noticed and understood it
  while reading the pre-written test.
-->

<!-- end_slide -->

## The struct tag footgun

**The compiler does not check struct tags for correctness. At all.**

```go
type Item struct {
    Name string `json:"geust"` // typo — meant "guest" or "name"
}
```

<!-- pause -->

No error. No warning. Every decode just leaves `Name` empty, silently,
forever — because nothing in the incoming JSON ever matches the key
`"geust"`.

**Demo:** introduce this typo live, send a request, watch the field come
back blank with no error anywhere. Fix the tag and watch it populate.

<!--
speaker_note: |
  This is worth running live rather than just describing. The confusing
  part isn't the bug itself, it's how quiet it is - nothing crashes,
  nothing logs, the field is just always empty. That silence is exactly
  what makes it a real production footgun, not a contrived example.
-->

<!-- end_slide -->

## Real-world angle: router and struct tag

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**A router is a hotel front desk's job board.**

A `GET` for `/items/{id}` goes to the lookup clerk. A `POST` for
`/items` goes to the new-stock clerk. The board just directs traffic by
verb and path — it doesn't care what happens after.

<!-- column: 1 -->

**A struct tag is a customs sticker on a crate.**

Just text someone wrote on the outside. Unverified by the shipping
company — only the receiving customs office (`encoding/json`, via
reflection) reads it. A typo doesn't get the crate rejected. It gets
misrouted, silently.

<!-- reset_layout -->

<!-- end_slide -->

## No automatic exception-to-status translation

Spring's `@ExceptionHandler` maps a thrown exception to an HTTP response,
often without any individual controller method thinking about it.

Go has no exceptions in that sense. **Every handler decides its own
status code, explicitly, every time.**

```go
if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
    http.Error(w, "invalid JSON", http.StatusBadRequest)
    return
}
```

<!-- pause -->

More typing. Also more visible — reading one handler tells you exactly
what status comes back under exactly which condition, with nothing
happening implicitly somewhere else in the stack.

<!-- end_slide -->

## Concurrency stops being abstract here

A REST server fields concurrent requests as a matter of course — each
incoming request typically runs on its own goroutine.

<!-- pause -->

An in-memory store backing your API needs the exact same
`sync.Mutex` / `sync.RWMutex` discipline from **Topic 7** — not "this
type happens to be safe," but a real, load-bearing requirement.

**Name it out loud:** this is the goroutines-and-channels topic showing
up for real, in a service you'll actually run.

<!--
speaker_note: |
  Make this connection explicit rather than assuming students notice it
  on their own. If Topic 7 landed well, this is the slide where its
  payoff becomes concrete instead of theoretical.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Building it: Handler -> Service -> Repository
===

<!-- end_slide -->

## Same pattern, now for real

Topics 8 and 9 introduced Handler -> Service -> Repository layering and
constructor-injected interfaces conceptually.

**This topic builds the whole thing, completely, for the first time.**

<!-- incremental_lists: true -->

- **Repository** — storage. An interface plus an in-memory implementation
- **Service** — validation and business rules. No HTTP in this layer at all
- **Handler** — HTTP only. Decode, call the service, respond
- **main** — the composition root that wires all three together

<!-- incremental_lists: false -->

If Topics 8-9 landed, none of this is new information — it's application.

<!-- end_slide -->

## The Repository layer

```go
type ItemRepository interface {
    Create(item domain.Item) (domain.Item, error)
    Get(id string) (domain.Item, error)
    List() ([]domain.Item, error)
    Update(item domain.Item) (domain.Item, error)
    Delete(id string) error
}

type InMemoryRepository struct {
    mu    sync.RWMutex
    items map[string]domain.Item
}
```

<!-- pause -->

The interface is the contract. `InMemoryRepository` is one implementation
of it — a `map` guarded by the `sync.RWMutex` from Topic 7. Nothing above
this layer knows or cares that it's a map and not a database.

<!-- end_slide -->

## The Service layer

```go
type ItemService struct {
    repo repository.ItemRepository // interface, not a concrete type
}

func NewItemService(repo repository.ItemRepository) *ItemService {
    return &ItemService{repo: repo}
}

func (s *ItemService) Create(item domain.Item) (domain.Item, error) {
    if item.Quantity <= 0 {
        return domain.Item{}, fmt.Errorf("%w: quantity must be positive", ErrValidation)
    }
    return s.repo.Create(item)
}
```

<!-- pause -->

**No `http.Request`. No `http.ResponseWriter`. Anywhere.** This layer
owns business rules and nothing else — it would look identical if the
transport were gRPC instead of REST.

<!-- end_slide -->

## The Handler layer

```go
type ItemHandler struct {
    svc ItemService // interface again
}

func NewItemHandler(svc ItemService) *ItemHandler {
    return &ItemHandler{svc: svc}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
    var item domain.Item
    json.NewDecoder(r.Body).Decode(&item)

    created, err := h.svc.Create(item)
    if err != nil {
        writeServiceError(w, err) // maps ErrValidation -> 400, etc.
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(created)
}
```

<!-- pause -->

Decode JSON in. Call the service. Translate the result to a status code
and JSON out. That's the entire job of this layer — no validation logic
leaks up into it, and no storage logic leaks up into it either.

<!-- end_slide -->

## Every request carries a `context.Context`

`r.Context()` is created fresh for each incoming request, and it cancels
automatically the moment the client disconnects or the request's deadline
elapses.

```go
func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	item, err := h.svc.Get(ctx, r.PathValue("id"))
	// ...
}
```

<!-- pause -->

Passing `ctx` down isn't optional plumbing — it's how "the client gave up
waiting" reaches all the way to a slow database call, so the server stops
doing work nobody's waiting for anymore.

<!-- end_slide -->

## Threading it through Service and Repository

Extend the interfaces from a few slides back to take `ctx context.Context`
as their first parameter — the standard-library convention, and the same
one gRPC uses (Topic 11, next).

```go
type ItemRepository interface {
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	Get(ctx context.Context, id string) (domain.Item, error)
	// ...
}
```

<!-- pause -->

A real database driver (`database/sql`, `pgx`) uses `ctx` to cancel an
in-flight query the instant the parent context cancels. The in-memory
repository from this lab ignores it safely — a map has nothing to cancel
— but the parameter stays part of the interface so swapping in a real
database later doesn't touch every signature in the codebase.

<!--
speaker_note: |
  Worth naming explicitly: this changes every layer's signature versus
  what's on the earlier Repository/Service/Handler slides. Frame it as
  "here's the natural next step, once a real dependency shows up" rather
  than implying the earlier slides were incomplete - the in-memory
  version genuinely doesn't need it.
-->

<!-- end_slide -->

## One panic shouldn't take down every other request

A single panicking handler — a nil map write, an out-of-range index,
whatever — would otherwise crash the entire process mid-request, taking
every other in-flight connection down with it.

```go
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

<!-- pause -->

Wrap it the same way as `Logging`: `handler.Recover(handler.Logging(mux))`.
This is the real use of `recover` that Topic 2 filed away as "reserved for
truly exceptional situations" — a panicking handler is exactly that, and
recovering here contains it to the one goroutine serving that request
instead of the whole process.

<!--
speaker_note: |
  Good callback slide - ask the room to recall Topic 2's "no exceptions
  for ordinary failure" slide before revealing this one. The payoff is
  that recover isn't a curiosity from two topics ago, it's a real,
  small piece of production middleware they'll recognize in most Go
  service codebases.
-->

<!-- end_slide -->

## Wiring it together: the composition root

```go
func main() {
    repo := repository.NewInMemoryRepository()
    svc := service.NewItemService(repo)
    h := handler.NewItemHandler(svc)

    mux := http.NewServeMux()
    mux.HandleFunc("POST /items", h.Create)
    mux.HandleFunc("GET /items/{id}", h.Get)
    // ...

    http.ListenAndServe(":8080", handler.Logging(mux))
}
```

<!-- pause -->

**`main` is the one place in the program that knows about every concrete
type.** Every layer below it only ever sees an interface. Swap
`InMemoryRepository` for a Postgres-backed one later, and this is the
only file that changes.

<!--
speaker_note: |
  This slide is the payoff for Topics 8-9's constructor-injection
  drilling. If students can read this main function and immediately
  explain why it's called a "composition root," the earlier topics did
  their job.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Running it locally with docker-compose
===

<!-- end_slide -->

## Why containerize a Go dev service at all

<!-- incremental_lists: true -->

- A consistent environment across every machine on the team — no "works
  on my laptop"
- An easy way to add a real dependency (a database) without installing
  it locally
- It matches how the service will actually run in production and CI

<!-- incremental_lists: false -->

This isn't Topic 13's job (full deployment story, multi-stage builds in
depth) — this is "get this service running on your machine, today."

<!-- end_slide -->

## A minimal multi-stage Dockerfile

```dockerfile
# --- build stage ---
FROM golang:1.22 AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# --- final stage ---
FROM alpine:3.19
COPY --from=builder /out/server /usr/local/bin/server
EXPOSE 8080
ENTRYPOINT ["server"]
```

<!-- pause -->

**Connect this back to Topic 1:** the final image is tiny precisely
*because* Go ships one dependency-free binary. There's no JVM, no
`node_modules`, no interpreter to drag into the final image — just the
binary, copied in.

<!-- end_slide -->

## docker-compose.yml: service + Postgres

```yaml
services:
  api:
    build: .
    ports: ["8080:8080"]
    environment:
      DB_HOST: db
      DB_USER: rest_service
    depends_on: [db]

  db:
    image: postgres:16
    environment:
      POSTGRES_USER: rest_service
      POSTGRES_PASSWORD: rest_service
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

<!-- pause -->

The service here still uses the in-memory repository — Postgres is
included to show what wiring in a **real dependency** looks like. That's
the natural next step once `InMemoryRepository` gets swapped out behind
the same `ItemRepository` interface.

<!-- end_slide -->

## WSL notes (read this if your dev machine is Windows)

Docker Desktop's WSL2 backend already runs containers inside a Linux VM.

<!-- incremental_lists: true -->

- Run `docker-compose` commands **from inside a WSL shell** — not
  PowerShell, not CMD — for correct performance and file permissions
- Keep your project on the **Linux filesystem** (`~/projects/...` inside
  WSL), not on a mounted Windows drive (`/mnt/c/...`)
- Cross-filesystem bind mounts between Windows and the WSL2 VM are
  measurably slower, and can make file-watching/hot-reload tools miss
  changes entirely
- A `.gitattributes` forcing `LF` line endings matters more here — a
  Windows-edited file with `CRLF` line endings, bind-mounted into a Linux
  container, can behave differently than the same file edited on Linux

<!-- incremental_lists: false -->

**This is scoped, deliberately** — these are the specific issues that
bite WSL + docker-compose setups, not a general Docker course.

<!--
speaker_note: |
  This section exists because the client specifically hit these issues
  running Docker Desktop + WSL2 for local dev. Keep it practical - if
  someone in the room isn't on Windows, this is a two-minute aside, not
  a detour.
-->

<!-- end_slide -->

## Demo: bring it up, hit it, bring it down

```
docker-compose up --build
```

```
curl localhost:8080/ping
```

```
docker-compose down
```

**Demo:** run all three from a terminal (WSL shell, if you're on
Windows). Watch the build log show the two Dockerfile stages running,
then confirm `/ping` responds from the host exactly like it did with
`go run`.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **`net/http` + `encoding/json` is enough**: no framework required, and
   that's a deliberate trade-off, not a gap
2. **Go 1.22 added real routing**: method + path params, no third-party
   router needed for most services now
3. **Struct tags drive JSON, unchecked by the compiler**: a typo fails
   silently, not loudly
4. **Every handler sets its own status code**: no automatic exception
   translation, fully visible instead
5. **Handler -> Service -> Repository, wired by constructor injection in
   `main`**: Topics 8-9's pattern, built out completely
6. **`r.Context()` carries cancellation downward**: thread it through
   Service and Repository so a client giving up actually stops the work
7. **`recover` in middleware contains a panic to one request**: the real
   production use of the `panic`/`recover` pair Topic 2 flagged as
   "reserved for truly exceptional situations"
8. **`docker-compose up` gets the whole thing running locally**,
   including a real dependency alongside it
9. **The layering pays off in tests, too**: the service is testable with
   zero HTTP involved, and the handler is testable with
   `net/http/httptest` — no running server, no real network call, ever

<!-- end_slide -->

## Back to the opening scenario

Spring's `@ExceptionHandler` translated a thrown exception into a status
code somewhere you never had to look.

**A Go handler decides its own status code, explicitly, every single
time — `writeServiceError` maps `ErrNotFound` to 404 and `ErrValidation`
to 400, in code you can read top to bottom.**

<!-- pause -->

**Type in chat: now that you've written a few of these — more typing, or
a fair trade for never having to go hunting for the exception handler
that decided your status code?**

<!--
speaker_note: |
  Don't force a winner here - the honest answer is genuinely "depends on
  team size and how much magic people are comfortable with." The point
  is that students can now articulate the trade-off specifically, not
  that Go's approach is objectively better.
-->

<!-- end_slide -->

## Bridge to Topic 11

**We've established:**

<!-- incremental_lists: true -->

- A full CRUD REST API needs nothing beyond the standard library
- Handler -> Service -> Repository, wired by constructor injection, is
  how a real Go service gets structured
- `context.Context` and a `recover` middleware are what keep one slow
  client and one bad request from taking down everyone else
- `docker-compose` gets that service — and its dependencies — running
  locally, WSL included
- That same layering means the service is testable without HTTP and the
  handler is testable without a running server, via `httptest`

<!-- incremental_lists: false -->

**Topic 11: gRPC and Protocol Buffers** — a second way to expose a
service, built for service-to-service calls rather than browsers and
`curl`, and where a strict schema replaces the JSON struct-tag footgun
you just saw.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
