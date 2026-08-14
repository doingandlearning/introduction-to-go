# Topic 6 sample code — Functional Programming

Three small programs, one module. Run everything from this directory
(`code/`).

## `cmd/closures`

`makeCounter` returns a closure holding its own private `count`. Creates
two independent counters and interleaves calls to both to prove neither
shares state with the other.

```
go run ./cmd/closures
```

## `cmd/filterdemo`

A single generic `Filter[T any]` used on `[]int` (keep evens) and then
on `[]string` (keep names longer than 3 characters) — same function,
two unrelated element types. Also a generic `Map[T, U any]` turning a
`[]Drink` into `[]float64` via `Drink.Dollars`, a method expression.

```
go run ./cmd/filterdemo
```

## `cmd/serveropts`

The full functional options example: `NewServer(opts ...ServerOption)`
built three ways — no options (defaults survive), one option, and two
stacked. Confirms defaults aren't disturbed when an option isn't
supplied.

```
go run ./cmd/serveropts
```

## A note on generics

`Filter` and `Map` both require Go 1.18+ (this module targets 1.22).
If you find older Go code solving the same problem, expect
`interface{}`, type assertions, or code generation instead — generics
weren't available until 2022.
