# Lab 6: Functional Programming — the Go Roasters order pipeline

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

Both directories contain the same shape: a `cmd/pipeline` package (the
program entry point, already wired to walk through the first five
exercises) and an `internal/orders` package (where most of the
implementation gaps live, including the test from Exercise 6). By the
end, you'll have written a closure, a pair of generic helpers, a
functional-options constructor, and a method value — the full toolkit
from the lecture, applied to one small coffee shop — and a test that
proves one of them actually works.

---

## Exercise 1: Independent order counters

**Objective:** Confirm a closure captures its own state, not shared
global state.

**Context:** `internal/orders/orders.go` has a `TODO` on
`NewOrderCounter`. It should return a closure that starts at 0 and
increments by 1 on every call — the same shape as `makeCounter` from
the lecture.

**Tasks:**

1. Implement `NewOrderCounter() func() int` in `internal/orders/orders.go`.
2. Run `go run ./cmd/pipeline` from `starter/`. Exercise 1's output
   creates two counters, `till1` and `till2`, and interleaves calls to
   both.
3. Confirm the sequences printed for `till1` and `till2` are each their
   own independent 1, 2, 3, … — not a single shared count split between
   them.

**Key Learning:** Each call to a function that returns a closure
produces a brand-new copy of whatever variables that closure captures.
Two counters from two calls never see each other's state, even though
both were built by the exact same code.

---

## Exercise 2: One `Filter`, two element types

**Objective:** Write a generic higher-order function once and use it
across unrelated types.

**Context:** `Filter[T any]` has a `TODO` stub that currently returns
every item unfiltered.

**Tasks:**

1. Implement `Filter[T any](items []T, predicate func(T) bool) []T` in
   `internal/orders/orders.go`, using a `for` loop and `append`.
2. Run `go run ./cmd/pipeline` again. Exercise 2 calls your `Filter`
   once on a `[]int` (keeping even numbers) and once on a `[]string`
   (keeping drink names longer than 3 characters) — same function,
   different type parameters.
3. Confirm both filtered lists match the `// want` comments next to
   each call in `cmd/pipeline/main.go`.

**Key Learning:** `T any` means the compiler generates the right
version of `Filter` for whatever type you call it with. You did not
write two functions, or reach for `interface{}` and a type assertion —
one generic definition covers both cases, with full type safety.

---

## Exercise 3: `Map` over a slice of structs

**Objective:** Use a second generic helper, `Map[T, U any]`, to turn a
slice of one type into a slice of a completely different type.

**Context:** `Drink` has a `PriceCents int` field and a `Dollars()
float64` method — a derived value computed from the stored one. Both
`Dollars` and `Map` have `TODO` stubs.

**Tasks:**

1. Implement `Dollars()` on `Drink` — convert `PriceCents` to a
   `float64` and divide by 100. (Watch out for integer division: convert
   before you divide.)
2. Implement `Map[T, U any](items []T, transform func(T) U) []U` using a
   `for` loop.
3. Run `go run ./cmd/pipeline`. Exercise 3 calls `Map` with
   `orders.Drink.Dollars` — notice this passes the method itself, not a
   wrapping closure. This works because `Drink.Dollars` used this way is
   a *method expression*: Go turns it into an ordinary
   `func(Drink) float64`, exactly the shape `Map`'s `transform`
   parameter wants.
4. Confirm the printed prices match the `// want` comment.

**Key Learning:** Generic helpers compose with methods, not just
freestanding functions. `Type.Method` (no receiver instance) is a
method expression — different from the method value you'll use in
Exercise 5, which binds to one specific instance.

---

## Exercise 4: Functional options for `CoffeeOrder`

**Objective:** Build a constructor with optional settings, Go-style —
no overloading, no default parameters, no giant config struct with
every field left as a zero value by hand.

**Context:** `CoffeeOrder` has three unexported fields (`size`,
`extraShot`, `oatMilk`). `CoffeeOption`, `WithSize`, `WithExtraShot`,
`WithOatMilk`, and `NewCoffeeOrder` all have `TODO` stubs that
currently do nothing.

**Tasks:**

1. Define `CoffeeOption` as `func(*CoffeeOrder)`.
2. Implement `WithSize`, `WithExtraShot`, and `WithOatMilk`, each
   returning a `CoffeeOption` closure that sets the relevant field.
3. Implement `NewCoffeeOrder(opts ...CoffeeOption) *CoffeeOrder` — set
   defaults first (`size: "medium"`, `extraShot: false`, `oatMilk:
   false`), then loop over `opts` and apply each one.
4. Run `go run ./cmd/pipeline`. Exercise 4 builds three orders: no
   options, one option (`WithSize("large")`), and three stacked
   (`WithSize`, `WithExtraShot`, `WithOatMilk` together).
5. Confirm the "plain" order still shows the defaults untouched — that's
   the part that has to work for the pattern to be worth using.

**Key Learning:** This is the idiom that actually earns its keep in
production Go. It replaces both a missing feature (constructor
overloading) and another missing feature (default/named parameters)
with one small piece of machinery: a variadic list of functions that
each know how to adjust one setting.

---

## Exercise 5: A method value that remembers its receiver

**Objective:** Assign a method to a variable without calling it, and
confirm it still knows which instance it belongs to later.

**Context:** `Invoice.Total()` is already implemented in
`internal/orders/orders.go` by the time you reach this exercise (you'll
have written it as part of getting the package to compile). The gap is
in `cmd/pipeline/main.go` itself, marked with a `TODO`.

**Tasks:**

1. Open `cmd/pipeline/main.go` and find `exerciseFiveMethodValue`.
2. Where the `TODO` is, assign `inv.Total` — **no parentheses, don't
   call it** — to a variable named `getTotal`.
3. Call `getTotal()` (now with parentheses) and print the result.
4. Run `go run ./cmd/pipeline`. Confirm it prints `13.5`, and notice you
   never passed `inv` back in when you called `getTotal()` — it was
   already bound in.

**Key Learning:** A method value captures its receiver at the moment of
assignment, the same way a closure captures a variable from its
enclosing scope. `getTotal` is a plain `func() float64` from that point
on — you can pass it around, store it, or hand it to another function
that has never heard of `Invoice`.

---

## Exercise 6: Prove it with a test

**Objective:** Apply the `go test` mechanics you learned in Topic 2 to
this topic's functional-options constructor — proving `NewCoffeeOrder`'s
defaults survive when no options are passed, and that each option only
touches the field it's responsible for.

**Context:** `internal/orders/orders_test.go` has a `TestNewCoffeeOrder`
function that currently just calls `t.Skip`. You already know the shape
from Topic 2 — a `_test.go` file, a `TestX(t *testing.T)` function,
`t.Errorf`/`t.Fatalf`, no framework required. This exercise doesn't
teach anything new about *how* to write a test; it asks you to point
that tool at `NewCoffeeOrder`.

**Tasks:**

1. Open `internal/orders/orders_test.go` and remove the `t.Skip` line
   from `TestNewCoffeeOrder`.
2. Build three `CoffeeOrder`s and assert on each:
   - `NewCoffeeOrder()` with **zero options** — assert `Size() ==
     "medium"`, `ExtraShot() == false`, `OatMilk() == false`.
   - `NewCoffeeOrder(WithSize("large"))` with **one option** — assert
     `Size() == "large"`, and that the other two fields are still at
     their defaults.
   - `NewCoffeeOrder(WithSize("small"), WithExtraShot(), WithOatMilk())`
     with **options stacked** — assert all three fields reflect the
     overrides.
3. Run `go test ./...` from `starter/`. All three assertions inside
   `TestNewCoffeeOrder` should pass.
4. Now deliberately break it: in `NewCoffeeOrder`, temporarily change
   the default `size` from `"medium"` to `"large"`. Run `go test ./...`
   again and confirm your zero-options assertion fails, naming the
   mismatch. Then put the default back and confirm the test passes
   again.

**Key Learning:** A test doesn't just prove the happy path once — it
proves the defaults keep behaving correctly every time this code
changes. The "zero options" case matters most for functional options:
it's easy for a future edit to silently break the default while
everyone's attention is on the new option being added, and that's
exactly the case a test catches instantly, instead of a customer
noticing their "plain" coffee arrived as a large.

---

## Summary

By the end of this lab you should be able to:

- Write a closure that returns independent state on every call, and
  explain why two instances never collide
- Write a generic function once and use it across multiple, unrelated
  element types
- Tell a method expression (`Type.Method`) apart from a method value
  (`instance.Method`), and know which one a generic `Map` call needs
- Build a functional-options constructor with working defaults, and
  explain what problem it solves that Go has no other tool for
- Assign a method to a variable without calling it, and predict that it
  still knows its receiver when called later
- Write a test that proves a functional-options constructor's defaults
  survive when no options are passed, and confirm it fails when they
  don't
