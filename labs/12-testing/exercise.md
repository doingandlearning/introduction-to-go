# Lab 12: Testing — writing tests for the billing package

Starter code is in `starter/` — but unlike every earlier lab in this
course, the production code there is already finished. What's missing
is the tests. `starter/internal/billing/billing_test.go` has real test
function signatures with `t.Skip("TODO ...")` bodies — it compiles and
`go test` runs (reporting everything as skipped) before you touch a
line. A complete reference is in `solution/` — don't look until you've
had a go.

Both directories share the same shape: an `internal/billing` package
(the code under test) and a `cmd/billingcli` package (a runnable demo,
unaffected by how finished your tests are — `_test.go` files are never
part of a normal build).

---

## Exercise 1: A basic test, and reading a failure on purpose

**Objective:** Write the simplest kind of Go test, then learn to read
what a failure actually tells you.

**Context:** `starter/internal/billing/billing_test.go` has
`TestTenPercentOff` stubbed out with `t.Skip`.

**Tasks:**

1. Replace the `t.Skip` line with a real test: call
   `billing.TenPercentOff(100)`, compare it to `90.0`, and call
   `t.Errorf` if they don't match — following the exact shape from the
   lecture.
2. Run `go test ./...` from `starter/`. Confirm it passes.
3. Open `starter/internal/billing/billing.go` and change
   `TenPercentOff` to multiply by `0.8` instead of `0.9`. Run
   `go test ./...` again and read the failure output carefully.
4. Put `0.9` back and confirm the test passes again.

**Key Learning:** `t.Errorf` shows you exactly the values you told it
to show — the call, the "got," and the "want" — and nothing else. It
doesn't tell you *why* they differ, or point at the line in
`billing.go` that caused it. Diagnosing the cause is still on you; the
test only tells you *that* something's wrong.

---

## Exercise 2: Testing for a panic

**Objective:** Wire up `defer`/`recover` inside a test to check that a
function panics when it's supposed to.

**Context:** `Divide` in `billing.go` panics instead of returning an
error when the divisor is zero. `TestDivideByZeroPanics` is stubbed out
in the starter.

**Tasks:**

1. Replace the `t.Skip` in `TestDivideByZeroPanics` with a real
   implementation: `defer` an anonymous function that calls `recover()`
   and calls `t.Error(...)` if the result is `nil`, then call
   `Divide(10, 0)` as the last line of the test.
2. Run `go test ./...`. Confirm it passes.
3. Temporarily remove the `panic(...)` call from `Divide` in
   `billing.go` (just make it return `0` for the zero case instead).
   Run `go test ./...` again — confirm `TestDivideByZeroPanics` now
   *fails*.
4. Put the panic back. Confirm the test passes again.

**Key Learning:** There's no `assertPanics` helper anywhere in the
standard library — `defer`/`recover` are the same primitives from
Topics 2 and 5, reused here to do real work. Step 3 matters as much as
step 1: a panic test that can't fail isn't testing anything.

*Hint if stuck:* `recover()` only does something meaningful when
called directly inside a deferred function. Calling it in the main
body of the test won't catch anything — by the time a panic actually
propagates up out of `Divide`, the only code still running is whatever
was deferred.

---

## Exercise 3: Table-driven tests with `t.Run`

**Objective:** Replace a family of near-identical tests with one table
and a loop — no special syntax, just a slice of structs.

**Context:** `TierDiscount` in `billing.go` has four branches based on
price. `TestTierDiscount` is stubbed out in the starter.

**Tasks:**

1. Replace the `t.Skip` in `TestTierDiscount` with a table: a slice of
   an anonymous struct with `name`, `price`, and `want` fields. Include
   at least four cases, each with a distinct, descriptive `name`.
2. Loop over the table and wrap each case in
   `t.Run(tc.name, func(t *testing.T) {...})`, doing the same
   `got`/`want` comparison as Exercise 1 inside each subtest.
3. Run `go test -v ./...` and confirm each case is reported separately
   by name.
4. Pick one case name from your table and run only that case with:
   `go test -run 'TestTierDiscount/<case name, spaces as underscores>'`

**Key Learning:** JUnit has `@ParameterizedTest`; pytest has
`@pytest.mark.parametrize`. Go has neither — table-driven tests are
just a slice of struct literals and a `for` loop, features you already
know. `t.Run` is what turns each row into an independently
pass/fail-reported, individually targetable subtest.

---

## Exercise 4: Coverage — finding and closing a gap

**Objective:** Use `go test -cover` and the HTML coverage report to
find an untested line, then write a test for it.

**Context:** Depending on which four cases you chose in Exercise 3,
it's likely your table doesn't exercise every branch of
`TierDiscount`.

**Tasks:**

1. Run `go test -cover ./...` and note the percentage.
2. Run `go test -coverprofile=cover.out ./...` followed by
   `go tool cover -html=cover.out`. This opens a browser view of
   `billing.go` with covered lines in green and uncovered lines in red.
3. Find at least one red (uncovered) line inside `TierDiscount`. Write
   a new test function (or add a case to your table) that exercises it.
4. Re-run the coverage report and confirm the line is now green.

**Key Learning:** `go test -cover` and `go tool cover -html=...` ship
with the toolchain — there's no equivalent of installing JaCoCo or
`coverage.py` as a separate dependency. Coverage tells you what code
*ran* during your tests, not whether your assertions were meaningful —
100% coverage with weak assertions still passes.

---

## Exercise 5: Benchmarking

**Objective:** Write a benchmark, run it, and read what it reports.

**Context:** `BenchmarkTenPercentOff` is stubbed out in the starter
with `b.Skip`.

**Tasks:**

1. Replace the `b.Skip` line with a loop:
   `for i := 0; i < b.N; i++ { billing.TenPercentOff(100) }` (adjust
   the package qualifier if needed — this file is already in package
   `billing`).
2. Run `go test -bench=. -benchmem ./...`. Note the reported `ns/op`
   (nanoseconds per call) and the allocation counts (`B/op`,
   `allocs/op`).
3. Make a small, deliberate change to `TenPercentOff` — for example,
   temporarily return `price - (price * 0.1)` instead of
   `price * 0.9` — and re-run the benchmark. Do the numbers move? Put
   the original implementation back afterward.

**Key Learning:** `b.N` isn't a number you pick — `go test` runs your
benchmark function repeatedly, adjusting `b.N` between runs until the
timing is stable, similar to a stopwatch operator deciding to time "a
second's worth of widgets" rather than insisting on exactly 100 no
matter how fast the line runs. Benchmarking lives in the same
`testing` package and the same `go test` command as everything else in
this lab — no separate JMH-style setup.

---

## Exercise 6 (optional/stretch): `testify`

**Objective:** Compare vanilla `testing.T` against a popular
third-party assertion library, and form your own opinion.

**Context:** `github.com/stretchr/testify`'s `assert` package is close
to a de facto standard in real Go codebases, but it is a third-party
*addition* — not part of the standard library. There is no built-in
`assertEquals` in Go; if you go looking for one, this is confirmation
it genuinely doesn't exist.

**Tasks:**

1. Look at `solution/internal/billing/testify_example.go.txt` — a
   version of `TestTenPercentOff` rewritten with `assert.Equal` instead
   of a manual `if`/`t.Errorf`. (It's named `.go.txt`, not `_test.go`,
   so the lab compiles offline without needing to fetch the dependency
   — see the comment at the top of that file for how to actually run
   it if you have network access.)
2. Compare it side by side with your Exercise 1 answer.
3. Decide: would you reach for `testify` in your own projects? Would
   your answer be different if you were writing training material for
   people seeing Go for the first time, versus writing tests for a
   production codebase?

**Key Learning:** Go deliberately ships a minimal `testing` package and
leaves assertion style as a community/team decision, not a language
one. That's a real trade-off, not a gap — plenty of experienced Go
teams standardize on `testify`, and plenty deliberately stick to
vanilla `if`/`t.Errorf` for exactly the transparency reason in the
task above.

---

## Summary

By the end of this lab you should be able to:

- Write a basic Go test using nothing but `if` and `t.Errorf`, and read
  a failure message for exactly what it does and doesn't tell you
- Test that a function panics using `defer`/`recover`, and prove the
  test is meaningful by confirming it fails when the panic is removed
- Write a table-driven test with `t.Run`, and target a single case by
  name with `go test -run`
- Use `go test -cover` and the HTML coverage report to find and close
  an untested branch
- Write and run a benchmark, and explain why `b.N` isn't a number you
  choose
- Explain what `testify` adds on top of `testing.T`, and that it is
  not part of the Go standard library
