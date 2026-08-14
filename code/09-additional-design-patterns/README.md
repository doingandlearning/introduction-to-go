# Topic 9 sample code — Additional Design Patterns

Five small programs, one module. Run everything from this directory
(`code/`).

## `cmd/strategy`

`DiscountStrategy` as a plain function type — no interface, no classes,
no factory. Three interchangeable discount functions plus a map that
picks one at runtime.

```
go run ./cmd/strategy
```

Compare this against Topic 8's Singleton or the interface-based
alternative you'll build in the lab — the point is that Go only reaches
for an interface here once a "strategy" needs more than one method or
state a closure can't hold.

## `cmd/abstractfactory`

`ConnectorFactory` with two concrete families, `MySQLFactory` and
`PostgresFactory`. `describeStack` is written entirely against the
interface — it never names a concrete type.

```
go run ./cmd/abstractfactory
```

Notice how much of this ceremony exists only because we have two real,
coordinated families of types to swap. With a single real implementation
this would be premature abstraction — see the slides for where the line
is.

## `cmd/decorator`

Two HTTP middleware functions — `LoggingMiddleware` and
`HeaderMiddleware` — stacked around one base handler, invoked in-process
with `httptest` so nothing needs a real listening port.

```
go run ./cmd/decorator
```

Swap the nesting order (`HeaderMiddleware(LoggingMiddleware(...))`
instead) and rerun — the log line still fires either way here, but in a
real auth-then-log vs log-then-auth chain, ordering changes what
actually gets recorded.

## `cmd/builder`

`QueryBuilder` with chained `Select`/`From`/`Where` calls and a `Build()`
that validates before returning. The second call in `main` deliberately
omits `From` to show the validation firing.

```
go run ./cmd/builder
```

## `cmd/di`

A slightly larger dependency graph than Topic 8's single-dependency
example: `OrderHandler` depends on `OrderService`, which depends on
**two** interfaces, `Repository` and `Notifier`. Wired once with real
implementations, then wired a second time with fakes — same service
code, unchanged, either way.

```
go run ./cmd/di
```

This is manual dependency injection: every wire is a visible constructor
call in `main`, not something a framework resolves for you at runtime.

## Formatting

```
gofmt -l .    # lists files that don't match gofmt's formatting
gofmt -w .    # rewrites them in place
```
