# Lab 8: Design Patterns — Singleton and Dependency Injection

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

Both directories contain the same shape: `internal/registry` (Exercises
1-3, the Singleton) and `internal/catalog` (Exercise 4, Service +
Repository with dependency injection), each with a runnable `cmd/`
package.

---

## Exercise 1: A naive Singleton

**Objective:** Build the unsafe version of Singleton and confirm it
gives you exactly one instance under normal, non-concurrent use.

**Context:** `starter/internal/registry/registry.go` has a `TODO` where
`GetSettings` should live. `loadSettings()` is already implemented and
simulates an expensive load.

**Tasks:**

1. Implement `GetSettings()`: if the package-level `instance` is `nil`,
   build one by calling `loadSettings()` and assign it; either way,
   return `instance`.
2. Run `go run ./cmd/singletonlab` from `starter/`. It calls
   `GetSettings()` twice and prints whether both calls returned the same
   pointer.
3. Confirm the output shows `same=true` and that both `%p` addresses
   match.

**Key Learning:** A package-level variable plus a "build it if nil"
check gets you a working Singleton for the easy case — one goroutine,
calling in sequence. The pointer really is shared. The problem is what
happens next.

---

## Exercise 2: Break it with concurrency

**Objective:** Reason about why the naive version is unsafe, without
needing to literally trigger the race to see the danger.

**Tasks:**

1. In `starter/cmd/singletonlab/main.go`, uncomment the commented-out
   block that fires 50 goroutines at `GetSettings()` concurrently.
2. Before running anything, answer on paper (or out loud): with the
   naive `if instance == nil { ... }` check from Exercise 1, what could
   go wrong if two goroutines both reach that `if` check before either
   one finishes assigning `instance`?
3. Run `go run ./cmd/singletonlab`. It will likely still "work" most of
   the time on your machine — that's the trap. A data race doesn't
   announce itself; it just occasionally, unpredictably, produces two
   different `*Settings` instances or a corrupted read.
4. If you have access to a full Go toolchain, run
   `go run -race ./cmd/singletonlab` and read what the race detector
   reports. If you don't, write one sentence describing what you'd
   expect it to flag and why.

**Key Learning:** "It ran fine when I tested it" is not evidence of
thread safety. Two goroutines can both observe `instance == nil` at the
same instant — the check-then-act sequence has a gap where another
goroutine can interleave.

---

## Exercise 3: Fix it with `sync.Once`

**Objective:** Replace the naive check with the idiomatic Go fix and
confirm the loader now only ever runs once.

**Tasks:**

1. In `starter/internal/registry/registry.go`, add a package-level
   `once sync.Once` alongside `instance`.
2. Rewrite `GetSettings()` to call `once.Do(func() { instance =
   loadSettings() })`, then return `instance`.
3. Re-run `go run ./cmd/singletonlab` with the concurrent block from
   Exercise 2 still active. Confirm it finishes cleanly every time you
   run it, however many times you try.
4. In one sentence, explain to a teammate what `sync.Once` guarantees
   that the naive `if instance == nil` check didn't.

**Key Learning:** `sync.Once.Do` makes "run this exactly once, safely,
under concurrent callers" a two-line fix instead of hand-rolled locking.
It's also a reminder that Singleton in Go is barely a "pattern" for the
easy case — the interesting engineering is entirely in the concurrency
guarantee, not in restricting instance count.

---

## Exercise 4: Service + Repository with dependency injection

**Objective:** Build a small Service that depends on a Repository
*interface*, then construct it twice — once with a real repository, once
with a fake — and see the same Service code work both times.

**Context:** `starter/internal/catalog/catalog.go` defines a `Book`
domain type and a `Repository` interface with TODOs. `InMemoryRepository`
(a real implementation) is already written for you.

**Tasks:**

1. Fill in the `Repository` interface: it needs one method,
   `FindByISBN(isbn string) (*Book, error)`.
2. Implement `NewService(repo Repository) *Service` — store `repo` on
   the returned `*Service`.
3. Implement `(s *Service) Describe(isbn string) (string, error)`: call
   `s.repo.FindByISBN(isbn)`, wrap any error with `fmt.Errorf` and `%w`,
   and on success return a string like `"Title (ISBN: 1234567890)"`.
4. Run `go run ./cmd/cataloglab` from `starter/`. It builds a
   `catalog.Service` twice — once with `InMemoryRepository` (real data),
   once with a hand-written `fakeRepository` (canned data) — and prints
   the result of both.
5. Write a short paragraph (3-5 sentences) answering: what did
   dependency injection actually buy you here? Specifically — what would
   you have had to change in `Service` or `Describe` to test it against
   fake data, if `Service` had instead called `NewInMemoryRepository`
   directly inside its own constructor instead of receiving a
   `Repository` interface as a parameter?

**Key Learning:** `Service` never imports or mentions
`InMemoryRepository` — it only knows about the `Repository` interface.
That's what makes swapping in `fakeRepository` possible with zero
changes to `Service` itself: dependency injection plus an interface,
together, are what buy you a testable seam. Either one alone isn't
enough.

---

## Exercise 5: Prove it with a test

**Objective:** Use what Topic 2 taught you about `go test` to prove, in
code, the exact claim Exercise 4's paragraph made you argue in prose —
that constructor injection plus an interface makes `Service` testable
with zero framework involved. As a bonus, prove `sync.Once`'s
single-initialization guarantee the same way.

**Context:** `starter/internal/catalog/catalog.go` should now have your
completed `Service`, `Repository`, and `InMemoryRepository` from Exercise
4. A new file, `starter/internal/catalog/catalog_test.go`, has TODOs for
two tests. `starter/internal/registry/registry_test.go` has a TODO for a
third, optional test against your Exercise 3 `sync.Once` fix.

**Tasks:**

1. In `catalog_test.go`, define a `fakeRepository` type that implements
   `Repository` — a struct with a `book *Book` field and an `err error`
   field, and a `FindByISBN` method that returns `f.err` if it's set,
   otherwise `f.book`. This is the same "test wiring" idea as
   `cmd/cataloglab/main.go`'s `fakeRepository`, but built directly inside
   the test this time.
2. Implement `TestServiceDescribeSuccess`: build a `fakeRepository` seeded
   with a known `*Book`, construct a `Service` with `NewService(fake)`,
   call `Describe`, and check the returned string matches what that
   book's `Title` and `ISBN` should produce.
3. Implement `TestServiceDescribeError`: build a `fakeRepository` with
   `err` set to a sentinel error, construct a `Service` with it, call
   `Describe`, and confirm the returned error wraps yours
   (`errors.Is(gotErr, wantErr)`).
4. Run `go test ./internal/catalog/...` (or read through it by hand if
   you don't have a Go toolchain available) and confirm both tests pass.
5. Now break something on purpose: in `TestServiceDescribeSuccess`,
   change the fake's seeded `Title` without updating `want`, or make
   `FindByISBN` return the wrong book entirely. Re-run the test, read the
   failure message it produces, then put the code back the way it was.
6. Bonus: implement `TestGetSettingsSingleInstanceUnderConcurrency` in
   `registry_test.go`. Launch 50 goroutines that each call
   `GetSettings()` and store the returned pointer at their own index in a
   pre-sized slice, `wg.Wait()`, then assert every stored pointer equals
   the first one. Don't assert anything about timing or ordering — only
   the outcome, checked once every goroutine has finished.
7. In one or two sentences: what would Tasks 1-3 have looked like if
   `Service` had called `catalog.NewInMemoryRepository(...)` directly
   inside its own constructor instead of taking a `Repository` parameter?
   Could you still have written `TestServiceDescribeError` without a real
   database, a mocking library, or monkey-patching?

**Key Learning:** This is the payoff for everything Topic 8 taught about
dependency injection. Because `Service` depends on the `Repository`
interface rather than a concrete type, `catalog_test.go` can define its
own tiny `fakeRepository` and hand it to the exact same `NewService`
constructor production code uses — no mocking framework, no reflection,
no test-only build tags. The same seam that let `cmd/cataloglab` swap
`InMemoryRepository` for a fake at runtime is what lets a test swap it in
at compile time. Delete the interface, and `TestServiceDescribeError`
stops being a three-line test — you'd need a real dependency, a database
double, or a much heavier testing tool to reach the same code path. DI
isn't just an architecture preference; it's what makes `go test`, on its
own, sufficient for testing business logic in isolation.

---

## Summary

By the end of this lab you should be able to:

- Implement the naive Singleton pattern and explain why it's unsafe
  under concurrent access, even when it "seems to work" in casual
  testing
- Fix a Singleton's concurrent-init race with `sync.Once` and explain
  what guarantee it provides
- Build a Service that depends on a Repository interface rather than a
  concrete type, and inject different implementations through its
  constructor
- Explain, in your own words, why constructor injection plus an
  interface — not either one alone — is what makes code testable without
  a mocking framework
- Write a test that constructs a `Service` with a hand-written fake
  `Repository`, and explain why that's only possible because of
  dependency injection through an interface
