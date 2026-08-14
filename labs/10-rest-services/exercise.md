# Lab 10: REST Services — an items CRUD API

Starter code is in `starter/` (TODOs to fill in, already laid out as
Handler -> Service -> Repository). A complete reference is in
`solution/` — don't look until you've had a go.

Both directories share the same shape:

```
cmd/server/main.go        composition root
internal/domain           the Item struct
internal/repository       storage: an interface + an in-memory impl
internal/service          validation and business rules
internal/handler          HTTP only
```

By the end, you'll have a full CRUD REST API running two ways: directly
with `go run`, and inside a container with `docker-compose`.

---

## Exercise 1: `GET /ping`

**Objective:** Confirm the toolchain and the layering skeleton both work
before writing anything real.

**Context:** `starter/cmd/server/main.go` already builds an empty
`http.ServeMux`. The composition root (`repo := ...`, `svc := ...`,
`h := ...`) is stubbed with `TODO`s you haven't filled in yet.

**Tasks:**

1. Add a `GET /ping` route directly in `main`, returning
   `{"status":"ok"}` as JSON, using only `net/http` and `encoding/json`.
2. Run `go run ./cmd/server` from `starter/`. Confirm with
   `curl localhost:8080/ping`.

**Key Learning:** A complete, working HTTP endpoint needs nothing beyond
the standard library — no router import, no project generator.

---

## Exercise 2: `POST /items` with validation

**Objective:** Build the first real endpoint, and see the three layers
talk to each other for the first time.

**Context:** `starter/internal/domain/item.go` has the `Item` struct
already defined with JSON tags. `starter/internal/service/service.go` has
a `TODO` where the `Quantity > 0` validation rule belongs.
`starter/internal/handler/handler.go` has a `TODO` in `Create`.

**Tasks:**

1. In `service.go`, implement the `validate` function: return an error
   (wrapping `ErrValidation`) if `Quantity <= 0` or `Name == ""`.
2. In `handler.go`, implement `Create`: decode the request body into an
   `Item`, call `h.svc.Create`, and return `201 Created` with the saved
   item, or `400 Bad Request` with a clear message on validation failure.
3. Wire `POST /items` to `h.Create` in `main.go`.
4. Test both paths with `curl` — a valid item, and one with
   `"quantity":0`.

**Key Learning:** Validation lives in the service layer, not the
handler. The handler's only job is translating the service's answer into
an HTTP response.

---

## Exercise 3: `GET /items/{id}` and `GET /items`

**Objective:** Add the in-memory repository and the read endpoints.

**Context:** `starter/internal/repository/repository.go` has
`InMemoryRepository` stubbed — the `map[string]domain.Item` field exists,
but `Get` and `List` have `TODO`s. Note the Go version you're running
(`go version`) — this exercise uses the 1.22+ `{id}` path-parameter
syntax.

**Tasks:**

1. Implement `Get` and `List` on `InMemoryRepository`, using
   `r.mu.RLock()` / `r.mu.RUnlock()` for both — they're reads.
2. Wire `GET /items/{id}` to `h.Get` and `GET /items` to `h.List` in
   `main.go`.
3. Confirm `GET /items/{id}` returns `404` for a missing ID, and `200`
   with the item for one that exists.
4. Confirm `GET /items` returns everything you've created so far.

**Key Learning:** `sync.RWMutex` lets reads run concurrently with each
other — the read lock (`RLock`/`RUnlock`) is a genuinely different tool
from the write lock, not just a naming convention.

---

## Exercise 4: `PUT /items/{id}` and `DELETE /items/{id}`

**Objective:** Complete the CRUD set with sensible status codes.

**Tasks:**

1. Implement `Update` and `Delete` on `InMemoryRepository` (both are
   writes — use `r.mu.Lock()` / `r.mu.Unlock()`), returning
   `ErrNotFound` if the ID isn't present.
2. Implement `Update` and `Delete` on the handler: `PUT` re-validates and
   returns `200` with the updated item; `DELETE` returns `204 No
   Content` on success.
3. Wire both routes in `main.go`.
4. Confirm deleting a nonexistent ID returns `404`, **not a panic**.

**Key Learning:** Nothing about "ID doesn't exist" is special-cased by
the language — it's an ordinary error value (`ErrNotFound`) that the
handler maps to a status code like any other.

---

## Exercise 5: Logging middleware

**Objective:** Wrap every handler in a Decorator (Topic 9), applied once
at the top of the stack rather than repeated per handler.

**Tasks:**

1. Implement `Logging` in `starter/internal/handler/middleware.go`: it
   takes an `http.Handler` and returns one that logs method, path, and
   duration, then delegates to the wrapped handler.
2. In `main.go`, wrap the whole `mux` in `handler.Logging(...)` when
   calling `http.ListenAndServe`.
3. Confirm every request — including `/ping` — gets logged.

**Key Learning:** One `Logging` call around the whole router covers
every route defined on it. You don't wrap each handler individually.

---

## Exercise 6: The struct tag footgun, on purpose

**Objective:** Feel the silent failure mode before you meet it by
accident in real code.

**Tasks:**

1. In `domain/item.go`, deliberately change one JSON tag — e.g.
   `Name string \`json:"geust"\`` instead of `` `json:"name"` ``.
2. Restart the server, `POST` an item with a correctly-spelled `"name"`
   field, and confirm the response comes back with `Name` empty — no
   error anywhere.
3. Fix the tag and confirm the field populates correctly.

**Key Learning:** `encoding/json` matches keys by exact string via
reflection, with zero compile-time checking. A typo doesn't fail loudly
— it just quietly stops working.

---

## Exercise 7: Run it with docker-compose

**Objective:** Bring the whole service up in a container, the same way
it would run on a teammate's machine or in CI.

**Context:** `code/Dockerfile` and `code/docker-compose.yml` (from the
topic's main `code/` directory) are a working reference if you want to
copy them into your lab directory rather than write your own from
scratch.

**Tasks:**

1. From your lab directory (with a `Dockerfile` and `docker-compose.yml`
   in place), run `docker-compose up --build`.
2. From your host machine, confirm `curl localhost:8080/ping` responds —
   the container is serving the same endpoint you tested with `go run`
   earlier.
3. Run `docker-compose down` and confirm the containers stop.

*If you're on Windows with WSL:* do all of this from inside a WSL shell,
with the project directory living on the Linux filesystem, not under
`/mnt/c/...`. See the slides for why.

**Key Learning:** The same binary that ran with `go run` on your laptop
is what ends up inside the container — Docker isn't a separate build
process, it's packaging the same `go build` output.

---

## Exercise 8: Prove it with a test

**Objective:** Put Topic 2's `go test` mechanics to work on the layered
architecture you've just built, and see why the layering was worth it —
each layer can be tested completely on its own.

**Context:** No test files exist yet for this lab.
`starter/internal/service/service_test.go` and
`starter/internal/handler/handler_test.go` are already laid out with the
real test function names and `t.Skip(...)` bodies — fill the bodies in,
don't rename the functions.

**Tasks:**

1. **Service layer, no HTTP involved at all.** In `service_test.go`,
   implement `TestItemService_Create_StoresValidItem`: construct a real
   `repository.NewInMemoryRepository()`, wrap it in
   `service.NewItemService(...)`, call `Create` with a valid item, and
   assert the returned item has a generated `ID` and the fields you sent.
   Confirm it's actually stored by calling `Get` with the returned ID.
2. Implement `TestItemService_Create_RejectsInvalidQuantity`: call
   `Create` with `Quantity: 0` and assert the error satisfies
   `errors.Is(err, ErrValidation)`.
3. Run `go test ./internal/service/...` and confirm both pass. Notice
   there's no server, no port, and no `net/http` import anywhere in this
   file — the service knowing nothing about HTTP is what makes this
   possible.
4. **Handler layer, with `net/http/httptest`.** In `handler_test.go`,
   implement `TestItemHandler_Create_ReturnsCreated`: build a real
   `service.NewItemService` over a real in-memory repository, wrap it in
   `NewItemHandler(...)`, build a request with
   `httptest.NewRequest(http.MethodPost, "/items", ...)` and a JSON body,
   record the response with `httptest.NewRecorder()`, call `h.Create(rec,
   req)` directly — no `mux`, no `ListenAndServe`, no open port — and
   assert `rec.Code == http.StatusCreated` and the decoded body carries
   the fields you sent.
5. Implement `TestItemHandler_Create_ReturnsBadRequestOnInvalidQuantity`
   the same way with `Quantity: 0`, asserting `rec.Code ==
   http.StatusBadRequest`.
6. Run `go test ./internal/handler/...` and confirm both pass.
7. **Break something on purpose.** Temporarily edit `validate` in
   `service.go` so it no longer rejects `Quantity <= 0` (comment out that
   check). Re-run `go test ./internal/service/...` and watch
   `TestItemService_Create_RejectsInvalidQuantity` fail with a clear
   message — not a panic, not a silent pass. Put the check back and
   confirm the test passes again.
8. From your lab directory, run `go test ./...` and confirm every
   package passes.

**Key Learning:** This is the payoff of Handler -> Service ->
Repository: because each layer depends on an interface rather than a
concrete neighbor, each layer is independently testable. The service
test proves the validation rule with no HTTP request in sight; the
handler test proves the status code and JSON shape with no running
server and no real network call — `httptest.NewRequest` and
`httptest.NewRecorder` simulate the entire HTTP transaction in-process.
Bolting tests on only at the end, after the whole stack exists, would
mean testing all three layers at once through a real server, with no way
to isolate which layer actually broke.

---

## Summary

By the end of this lab you should be able to:

- Build a complete CRUD REST API using only `net/http` and
  `encoding/json`
- Structure a real service as Handler -> Service -> Repository, with
  each layer depending on an interface, wired together by constructor
  injection in a composition root
- Explain why validation belongs in the service layer and HTTP status
  mapping belongs in the handler layer
- Use `sync.RWMutex` correctly for a read-heavy in-memory store
- Recognize the struct-tag JSON footgun on sight, from having caused it
  once deliberately
- Bring the same service up locally with `docker-compose`, and explain
  what changes (and what doesn't) about running it that way
- Test the service layer directly with no HTTP involved, and test the
  handler layer with `net/http/httptest` instead of a running server —
  each layer proven in isolation
