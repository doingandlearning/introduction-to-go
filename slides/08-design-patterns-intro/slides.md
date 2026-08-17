---
title: "**Intro to Design Patterns**"
sub_title: Go Programming — Topic 8
author: Kevin Cunningham
---

## Opening scenario

In a Java shop, a `UserService` class calls `new PostgresUserRepo()`
directly inside one of its methods. It works fine in production.

Now you need a unit test for `UserService` that doesn't touch a real
database.

**Type in chat: with the code as described, how do you swap that
`PostgresUserRepo` for a fake inside a test — without editing
`UserService` itself?**

We'll come back to this once we've looked at how Go handles dependencies
between types.

<!--
speaker_note: |
  Most answers will involve some form of mocking framework, reflection,
  or "you can't, not without changing the class." That's the point -
  tight coupling to a concrete type is a self-inflicted problem, and by
  the end of this session the fix will look almost too simple.
-->

<!-- end_slide -->

## What a design pattern actually is

A **design pattern** is a named, reusable solution to a recurring design
problem — not a library you import, a shape you recognize and re-implement.

The most famous catalogue is the **Gang of Four** book (1994): 23 patterns,
written against class-based OOP in C++ and Smalltalk.

<!-- pause -->

**That matters here.** Go has structs, interfaces, and no inheritance.
Some GoF patterns translate directly. Some shrink to almost nothing.
A few barely make sense at all. We'll be explicit about which is which,
every time.

<!-- end_slide -->

## A framework for judging every pattern

Before we look at Singleton, here's the lens we'll apply to **every**
pattern in this course — today and in Topic 9.

<!-- incremental_lists: true -->

- **Testability** — can you swap the real thing for a fake in a test?
- **Coupling** — does using this pattern make code more or less dependent
  on concrete types?
- **Discoverability** — can a new team member read the code and see
  what's happening, or does the pattern hide control flow?
- **Is there a simpler Go-native way?** — several GoF patterns shrink or
  disappear once you have first-class functions, interfaces, and package
  scope. Singleton is the first example.

<!-- incremental_lists: false -->

**We are not cataloguing patterns to worship them.** We're asking, every
time: does this pull its weight in Go, or is it OOP muscle memory?

<!-- end_slide -->

<!-- jump_to_middle -->

Singleton
===

<!-- end_slide -->

## Singleton: exactly one instance

**Singleton** guarantees a type has exactly one instance, globally
accessible.

```go
package config

var instance *Config

type Config struct {
	APIKey string
}

func GetConfig() *Config {
	if instance == nil {
		instance = &Config{APIKey: loadFromEnv()}
	}
	return instance
}
```

<!-- pause -->

**This is broken.** Two goroutines can both see `instance == nil` at the
same time, both allocate, and now you've built two "singletons."

<!-- end_slide -->

## The idiomatic fix: `sync.Once`

```go
var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{APIKey: loadFromEnv()}
	})
	return instance
}
```

<!-- pause -->

`once.Do(f)` guarantees `f` runs **exactly once**, no matter how many
goroutines call `GetConfig()` concurrently. Everyone else blocks until the
first call finishes, then proceeds with the result already set.

<!--
speaker_note: |
  Draw the analogy: a building's single reception desk. You get "exactly
  one instance" almost for free, the same way a building naturally has
  one front desk - nobody had to enforce that architecturally. sync.Once
  is the automatic shutter that only the first person through the door
  can trigger; everyone else in the same rush just finds the desk
  already open.
-->

<!-- end_slide -->

## Go already has a singleton scope

Here's the twist: in Java, you reach for Singleton specifically because
the language lets **anyone** write `new Config()`. Go doesn't have that
problem the same way.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// This IS a singleton.
// Zero ceremony.
var Cfg = loadFromEnv()
```

Package-level vars exist **exactly once**, initialized before `main`
runs.

<!-- column: 1 -->

```go
// This is for a narrower case:
// expensive/side-effecting init
// you want to skip entirely if
// unused.
once.Do(...)
```

<!-- reset_layout -->

**Singleton is close to a non-pattern in Go for the common case.** A
package already behaves like a built-in singleton scope.

<!-- end_slide -->

## No real "private constructor" trick

Go has no constructors and no access modifier finer than the package
boundary.

An unexported struct type (`type config struct{}`) blocks **other**
packages from constructing one directly, forcing them through
`GetConfig()`.

<!-- pause -->

**But code in the same package can still build a second one.** This is a
convention enforced by the package boundary — not an airtight guarantee
like Java's `private` constructor plus a static factory.

<!-- end_slide -->

## Singleton: pros and cons

<!-- column_layout: [1, 1] --> 

<!-- column: 0 -->

# Pros

- Guarantees a single instance without ceremony, via package scope 
- `sync.Once` solves the concurrent-init race correctly in a few lines 
- Cheap to reach for — often just a package-level `var` 

<!-- pause -->

<!-- column: 1 -->

# Cons

- Global mutable state is harder to test — no easy way to swap in a fake config without touching global state or adding an interface anyway 
- Hides a dependency inside function bodies — a function calling `GetConfig()` internally isn't honest about depending on config, versus taking it as a parameter 
- Can make tests interfere with each other if the singleton holds mutable state across test runs 

<!-- pause -->

<!-- reset_layout -->

**Type in chat: for a CLI tool that runs once and exits, is the
global-state downside even real? What changes for a long-running service
with hundreds of tests?**

<!--
speaker_note: |
  Push for a real answer here, not just "it depends." A one-shot CLI
  tool has no test-interference problem because there's no second test
  run sharing process state. A service with a growing test suite is
  exactly where the hidden-dependency cost compounds - every new test
  either shares state or has to route around it.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Handler, Service, Repository
===

<!-- end_slide -->

## Structuring a Go web service

The client calls this "handler/controller/service." Go doesn't have a
fixed vocabulary the way Spring MVC does — naming varies by shop. Here's
the version we'll standardize on for this course:

<!-- incremental_lists: true -->

- **Handler** — owns the HTTP concern only: parse the request, call the
  service, translate the result into a status code and body
- **Service** — owns business logic: validation, orchestration, domain
  rules. Knows nothing about HTTP
- **Repository** — owns data access: a database, an in-memory map, an
  external API. No business rules of its own

<!-- incremental_lists: false -->

**This is a preview.** Topic 10 builds this out with a real REST service
end to end.

<!-- end_slide -->

## Handler owns HTTP, nothing else

A service method should be callable from a CLI, a gRPC handler
(Topic 11), or a test — with **no HTTP involved at all**.

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.Service.FindUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
```

**Nothing in here validates business rules.** It parses, delegates, and
translates the result into an HTTP response.

<!-- end_slide -->

## The layering, in code

```go
type Handler struct {
	Service UserService // interface, not a concrete type
}

type UserService struct {
	Repo UserRepository // interface, not a concrete type
}

type UserRepository interface {
	FindByID(id string) (*User, error)
}
```

<!-- pause -->

Each layer only knows about an **interface** one level down, never the
concrete type. That's not an accident — it's what makes each layer
independently testable, and it's where dependency injection enters. More
on that next.

<!-- end_slide -->

## Handler/Service/Repository: pros and cons


<!-- column_layout: [1,1] -->
<!-- column: 0 -->
# Pros

- Each layer is independently testable — fake the service in a handler test, fake the repository in a service test 
- Business logic isn't tangled with HTTP parsing or SQL 
- Clear place for new code to go as the service grows 
<!-- column: 1 -->

<!-- pause -->

# Cons

- For a genuinely small service, three layers can be more ceremony than the problem needs — a 40-line CRUD toy doesn't need this 
- Over-abstracting early is a real anti-pattern — an interface with exactly one implementation, forever, buys you nothing 
- Go culture leans "write the concrete version first, extract an interface when you have two implementations or a real testing need" — a different default from shops that interface everything up front 

<!-- pause -->

<!-- reset_layout -->
**Type in chat: for a service with one endpoint and one table, do you
reach for all three layers on day one, or add them when the second
endpoint shows up?**

<!--
speaker_note: |
  There's no single right answer, but push people to name the trigger
  that would make them introduce the extra layer - "when I add a second
  repository implementation" or "when the handler test starts needing a
  real database" are good, concrete answers. "Always" and "never" are
  both worth pushing back on.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Dependency Injection
===

<!-- end_slide -->

## No framework, no container

Go has no DI framework/container culture the way Spring (`@Autowired`) or
Angular do. No reflection-based container resolves your object graph for
you.

<!-- pause -->

Tools like `google/wire` and `uber-go/fx` exist — worth a one-line
mention. They **exist, are not required, and most teams don't reach for
them.**

The idiomatic Go approach is simpler: **manual constructor injection via
interfaces.**

<!-- end_slide -->

## Constructor injection

```go
type UserRepository interface {
	FindByID(id string) (*User, error)
}

type UserService struct {
	repo UserRepository // depends on the interface, not a concrete DB type
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}
```

<!-- pause -->

`UserService` never mentions Postgres, an in-memory map, or anything
concrete. It depends on a **shape**, not a type — this is Topic 5's
implicit interface satisfaction, doing real work.

<!-- end_slide -->

## Composition root

Somewhere — usually `main()`, or a small setup function — concrete types
get constructed and wired together. The standard term for this is the
**composition root**.

```go
func main() {
	repo := postgres.NewUserRepo(db)      // concrete, real
	svc := NewUserService(repo)           // wired here, once
	handler := &UserHandler{Service: svc}

	http.ListenAndServe(":8080", handler)
}
```

<!-- pause -->

In a test, the composition root is the test function itself:

```go
svc := NewUserService(fakeRepo{users: canned})
```

**Same `UserService` code. Zero framework involved.**

<!-- end_slide -->

## Dependency injection: pros and cons

<!-- column_layout: [1,1] -->

<!-- column: 0 -->
# Pros

- Explicit — read a constructor signature and you see exactly what depends on what, nothing injected by magic or reflection 
- Fast — no container startup, no reflection overhead 
- No framework lock-in — it's just passing a value that satisfies an interface 
<!-- column: 1 -->

<!-- pause -->

# Cons


- At real scale, wiring a large dependency graph by hand in `main` gets verbose 
- That verbosity is exactly the itch tools like `wire` scratch — code-generating the wiring, still with no runtime reflection 
- Manual DI is a real tradeoff, not strictly better in every case — it's a decision, not a free lunch 

<!-- pause -->

<!-- reset_layout -->
**Type in chat: for a service with 3 dependencies vs. one with 30, does
manual wiring in `main` still feel "explicit," or does it start feeling
like boilerplate?**

<!--
speaker_note: |
  Three dependencies wired by hand reads as clear documentation. Thirty
  starts to look like the exact problem wire.go generation exists to
  solve. Let the room reach that conclusion themselves before naming
  wire explicitly again.
-->

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Every pattern gets judged on testability, coupling, discoverability,
   and "is there a simpler Go-native way"** — not adopted by default
2. **Singleton is close to a non-pattern in Go**: package scope already
   gives you one instance; `sync.Once` only earns its keep for
   expensive/side-effecting init
3. **Handler → Service → Repository**: HTTP concerns, business logic, and
   data access as three independently testable layers — but don't add
   layers a small service doesn't need yet
4. **Dependency injection in Go is manual constructor injection via
   interfaces** — no container, no framework, wired once at the
   composition root
5. **Interfaces + constructors are what make code testable** — this ties
   straight back to Topic 5's implicit interface satisfaction
6. **DI's real payoff shows up in `go test`**: swap in a hand-written
   fake dependency in a test, with zero mocking framework involved — the
   lab has you write exactly that test against today's `Service`

<!-- end_slide -->

## Back to the opening scenario

The Java `UserService` called `new PostgresUserRepo()` directly — no way
to substitute anything without editing the class.

**In Go, `UserService` never mentions `PostgresUserRepo` at all.** It
takes a `UserRepository` interface in its constructor. Production code
passes the real repo. A test passes a fake one. Same `UserService`, zero
changes.

<!-- pause -->

**Type in chat: what made that swap possible — was it dependency
injection itself, or the interface it was injected as?**

<!--
speaker_note: |
  The honest answer is both, and neither is sufficient alone. Passing a
  concrete PostgresUserRepo as a constructor argument (DI without an
  interface) still couples you to Postgres. An interface with no
  constructor injection (still calling a global GetRepo()) still hides
  the dependency. You need both together.
-->

<!-- end_slide -->

## Bridge to Topic 9

**We've established:**

<!-- incremental_lists: true -->

- A pros/cons framework: testability, coupling, discoverability, "is
  there a simpler Go-native way"
- Singleton, and why Go often needs less ceremony than the GoF version
- Handler → Service → Repository as this course's layering standard
- Dependency injection as manual constructor injection via interfaces
- Today's lab: a real `go test` against `Service`, wired up with a fake
  `Repository` instead of the real one — proof, not just theory

<!-- incremental_lists: false -->

**Topic 9: Additional Design Patterns** — goes deeper on dependency
injection at scale, and covers **Strategy, Abstract Factory, Decorator,
and Builder**, each run through the same pros/cons framework you just
learned.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
