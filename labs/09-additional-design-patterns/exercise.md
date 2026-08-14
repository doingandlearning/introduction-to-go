# Lab 9: Additional Design Patterns — TicketDesk

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

Both directories share the same shape: one `cmd/<name>` package per
exercise, each a standalone `main.go`. The scenario running through all
six is **TicketDesk**, a small event-ticketing system — discount rules,
swappable storage backends, an HTTP layer, report queries, and the
service wiring that holds it together.

---

## Exercise 1: Strategy, two ways

**Objective:** Build the same three ticket discounts as a function type,
then again as an interface + structs, and compare the ceremony.

**Context:** `cmd/discounts/main.go` has TODOs for both parts. Part A
mirrors the lecture's `DiscountStrategy` function type. Part B is the
same three discounts (no discount, member, early bird) built as a
`Discounter` interface with one struct per implementation.

**Tasks:**

1. Implement Part A: `TicketDiscount`, the three discount functions, and
   `ApplyTicketDiscount`.
2. Implement Part B: the `Discounter` interface and its three structs.
3. Run `go run ./cmd/discounts` from `starter/`. Confirm both parts
   print the same three numbers.
4. Count the lines each part needed. Which one would you rather extend
   with a fourth discount?

**Key Learning:** Go only reaches for an interface-based Strategy once
an implementation needs more than one method or state a closure can't
hold — for three stateless, single-purpose rules, the function-type
version does the same job with less ceremony.

---

## Exercise 2: Abstract Factory for storage backends

**Objective:** Swap an entire coordinated family of types (a store and
its factory) behind one interface.

**Context:** `cmd/connectors/main.go` has TODOs for a simplified,
single-method `StoreFactory` with two concrete factories: `LocalFactory`
and `CloudFactory`.

**Tasks:**

1. Implement `Store`, `StoreFactory`, the two `Store` types, and the two
   factories.
2. Implement `describeBackend(f StoreFactory)` — it must never name a
   concrete type.
3. Run `go run ./cmd/connectors` from `starter/`, calling
   `describeBackend` once with each factory. Confirm each prints its own
   backend's description.

**Key Learning:** The entire value of Abstract Factory is that
`describeBackend` never changes when the backend does — only the one
line in `main` that picks which factory to pass in.

---

## Exercise 3: Decorator via HTTP middleware

**Objective:** Confirm that decorator stacking order is something you
can observe, not just reason about.

**Context:** `cmd/middleware/main.go` has TODOs for `LoggingMiddleware`
and a second middleware, `HeaderMiddleware`, wrapped around a base
handler.

**Tasks:**

1. Implement `LoggingMiddleware` and `HeaderMiddleware`.
2. Chain both around `baseHandler()` and invoke the chain with
   `httptest.NewRequest` / `httptest.NewRecorder`. Print status code,
   the `X-Ticket-System` header, and the body.
3. Run `go run ./cmd/middleware` from `starter/`, then swap the nesting
   order and rerun.
4. Write a one-sentence comment on what changed (or didn't) between the
   two orderings, and why.

**Key Learning:** With two well-behaved middlewares, output is identical
either way — the moment one middleware can short-circuit the chain (an
auth check, for example), ordering stops being cosmetic.

---

## Exercise 4: Builder with validation

**Objective:** Implement a Builder whose `Build()` step refuses to
produce a result until required fields are set.

**Context:** `cmd/querybuilder/main.go` has TODOs for a `ReportBuilder`
with chained `Columns`, `Table`, and `OrderBy` methods.

**Tasks:**

1. Implement `ReportBuilder` and its chained methods.
2. Implement `Build() (string, error)`, returning an error if `Columns`
   or `Table` was never called.
3. Run `go run ./cmd/querybuilder` from `starter/`. Confirm a fully
   chained call succeeds, and a call missing `Table()` returns a
   non-nil error.

**Key Learning:** `Build()` gives you a natural place to enforce
required steps — a plain struct literal has no equivalent hook.

---

## Exercise 5: The same builder, as functional options

**Objective:** Feel the actual tradeoff between Builder and functional
options, rather than being told about it.

**Context:** `cmd/queryoptions/main.go` has TODOs to rebuild
`ReportBuilder` from Exercise 4 as `ReportOption` functions plus
`NewReport(opts ...ReportOption)`.

**Tasks:**

1. Implement `ReportConfig`, `ReportOption`, `WithColumns`, `WithTable`,
   `WithOrderBy`, and `NewReport`.
2. Try to replicate the same validation: `NewReport` should still return
   an error if columns or table were never supplied.
3. Run `go run ./cmd/queryoptions` from `starter/` and confirm both the
   success and error cases behave the same as Exercise 4.
4. Write one sentence: did enforcing "Table is required" feel natural
   with functional options, or did it feel forced?

**Key Learning:** Functional options and Builder can produce the same
final validation, but Builder makes "you must set this before Build
succeeds" a property of the chain itself; functional options only catch
it after the fact, inside the constructor. Neither is universally
better — pick based on whether your options are genuinely independent.

---

## Exercise 6: Manual DI with two dependencies

**Objective:** Wire a service that depends on two interfaces, then
rewire it with fakes, and confirm the service code never changes.

**Context:** `cmd/orderwiring/main.go` has TODOs for `TicketService`,
which depends on an `Inventory` interface and an `Emailer` interface —
plus real implementations (`WarehouseInventory`, `SMTPEmailer`) and fake
ones (`FakeInventory`, `FakeEmailer`).

**Tasks:**

1. Implement `Inventory`, `Emailer`, `TicketService`,
   `NewTicketService`, and `Book`.
2. Implement the two real implementations and the two fakes.
3. In `main`, wire a `TicketService` with the real implementations and
   call `Book` once.
4. Wire a **second** `TicketService` with the fakes and call `Book`
   again. Print what the fakes recorded.
5. Confirm `TicketService`, `NewTicketService`, and `Book` did not
   change between the two wirings — only the constructor arguments did.

**Key Learning:** This is manual dependency injection: the entire graph
is visible as plain constructor calls in `main`, and substituting test
doubles requires touching nothing but the wiring — no framework, no
reflection, no mocking library.

---

## Exercise 7: Prove it with a test

**Objective:** Write tests for three of today's patterns — the Strategy
functions in `discounts`, the Builder validation in `querybuilder`, and
the two-dependency wiring in `orderwiring` — confirming each is as
testable in practice as the lecture claimed.

**Context:** Since Topic 2 you've used `_test.go` files, `TestX(t
*testing.T)`, and `t.Errorf`/`t.Fatalf` in every lab. Today's targets
are `package main`, exactly like `errdemo_test.go` back in Topic 2 —
you're testing unexported functions and types directly in the same
package, no exported API required.

**Tasks:**

1. In `starter/cmd/discounts`, create `discounts_test.go`. Write
   `TestDiscountFunctions` as a table-driven test (`t.Run` per case)
   that calls `NoDiscount`, `MemberDiscount`, and `EarlyBirdDiscount`
   directly against a fixed price and asserts the exact result each
   returns. Write a second test, `TestApplyTicketDiscount`, that calls
   `ApplyTicketDiscount` with two different `TicketDiscount` values and
   confirms the result actually changes when the strategy does.
2. In `starter/cmd/querybuilder`, create `querybuilder_test.go`. Write
   `TestBuildSuccess` covering the happy path — `Columns`, `Table`, and
   `OrderBy` all set — and assert the exact SQL string `Build()`
   returns. Write `TestBuildMissingTable` that omits `Table()` and
   asserts `Build()` returns a non-nil error.
3. In `starter/cmd/orderwiring`, create `orderwiring_test.go`. Write
   `TestBookWithFakes` that constructs a `TicketService` with the
   `FakeInventory` and `FakeEmailer` you already implemented for
   Exercise 6, calls `Book`, and asserts what each fake recorded.
4. Run `go test ./...` from `starter/`. All three test files should
   currently report `SKIP` (they're stubbed with `t.Skip`) — replace
   each `t.Skip` with a real implementation and get them passing.
5. Deliberately break something and confirm your tests catch it: change
   `EarlyBirdDiscount`'s multiplier from `0.70` to `0.75`, or delete the
   `b.columns == "" || b.table == ""` check in `Build()`, or comment out
   the `s.mailer.Send` call in `Book`. Rerun `go test ./...`, confirm
   the right test fails with a useful message, then revert the change.

**Key Learning:** The payoff scales with how the pattern was built. A
Strategy function is a pure function — testing it costs nothing beyond
calling it with an input and checking the output, no setup at all. A
Builder's validation is exactly the kind of rule you want caught in CI
before a caller ever sees a malformed query. And a two-dependency DI
graph is no harder to fake than the one-dependency version from Topic
8 — `TicketService` only ever sees `Inventory` and `Emailer`, so it has
no way to tell that `main` (or a test) swapped in fakes.

---

## Summary

By the end of this lab you should be able to:

- Choose between a function-type and an interface-based Strategy based
  on whether an implementation needs more than one method or held state
- Write an Abstract Factory that lets a caller stay ignorant of which
  concrete family it's using
- Explain when Decorator ordering is purely cosmetic and when it changes
  behavior
- Implement a Builder whose `Build()` enforces required steps, and
  explain why functional options can't naturally do the same
- Wire a multi-dependency service by hand in `main`, and substitute test
  doubles without touching the service's own code
- Write tests for a Strategy function, a Builder's validation path, and
  a two-dependency DI wiring — the same `_test.go` habit from Topic 2,
  now applied to today's patterns
