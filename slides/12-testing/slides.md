---
title: "**Testing**"
sub_title: "Go Programming — Topic 12"
author: Kevin Cunningham
---

## Opening scenario

You're setting up a new Go project. Next decision: which test
framework?

**Type in chat: think about the JUnit-vs-TestNG conversation, or the
pytest-vs-unittest conversation, you've had before. How long did that
decision take — and who ended up making it?**

We'll come back to this once you've seen what Go actually gives you on
day one.

<!--
speaker_note: |
  Let answers land in chat for 20-30 seconds. Expect stories about
  a whole meeting, a wiki page, a linter config PR that took a week to
  get approved. Bank a few of the specifics - which framework, how the
  decision got made, how long it took - you'll come back to them
  directly at the end.
-->

<!-- end_slide -->

## There is no decision

Go ships a `testing` package and a `go test` command as part of the
standard toolchain. Nothing to install. Nothing to compare. Nothing to
argue about in a PR review.

<!-- pause -->

A test file sits right next to the code it tests, named with a
`_test.go` suffix — no separate test directory, no config file
declaring where tests live.

<!-- pause -->

**Let that sit for a second: the framework decision that ate a meeting
in your last project doesn't exist here.**

<!--
speaker_note: |
  Deliberately slow down here rather than rushing to the code example.
  For a room with JUnit/pytest/Jest backgrounds, "there is no decision"
  is a genuine moment of relief, not a throwaway line - give it room to
  land before moving on.
-->

<!-- end_slide -->

## This isn't your first test

You wrote your first `_test.go` file in **Topic 2**, testing `divide`.
Every lab since then ended the same way.

<!-- pause -->

<!-- incremental_lists: true -->

- **Topic 7**: a worker pool test that asserted the aggregate result, not execution order — because concurrent code can't be tested by asserting timing
- **Topic 8**: a service tested against a fake repository, built in ten lines, with zero framework involved — dependency injection's actual payoff
- **Topic 10**: an HTTP handler tested with `httptest`, no real server, no network call
- **Topic 11**: a gRPC method called directly, because it's just a Go method with `context.Context` first

<!-- incremental_lists: false -->

**Today doesn't introduce testing. It formalizes a habit you already have** — table-driven cases, panic testing, coverage, benchmarking.

<!--
speaker_note: |
  This is the reframe for the whole topic: resist any temptation to
  present today's material as "now let's learn testing." They've been
  doing it for ten topics. Today gives the habit sharper tools -
  t.Run, coverage percentages, benchmarks - not a new skill from zero.
  If the room has been diligent about the "prove it with a test"
  exercises, name specific things you remember them building.
-->

<!-- end_slide -->

## The basic shape

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// discount.go
package discount

func TenPercentOff(price float64) float64 {
    return price * 0.9
}
```

<!-- column: 1 -->

```go
// discount_test.go
package discount

import "testing"

func TestTenPercentOff(t *testing.T) {
    got := TenPercentOff(100)
    want := 90.0
    if got != want {
        t.Errorf(
            "TenPercentOff(100) = %v, want %v",
            got, want,
        )
    }
}
```

<!-- reset_layout -->

**Demo:** `go test ./...` from the module root. No config, no test
runner to point at a directory — it finds every `_test.go` file on its
own.

<!--
speaker_note: |
  Run this live. Point out what ISN'T here compared to JUnit/pytest:
  no @Test annotation, no test class, no discovery config, no fixture
  decorator. Just a function named Test-something taking *testing.T.
-->

<!-- end_slide -->

## No built-in assertion library

Vanilla `testing.T` gives you exactly **four** reporting methods.

<!-- incremental_lists: true -->

- `Error` / `Errorf` — mark the test failed, keep running the rest of the function
- `Fatal` / `Fatalf` — mark it failed, stop this test function immediately

<!-- incremental_lists: false -->

That's it. No `assertEquals`. No `expect(x).toBe(y)`.

<!-- pause -->

**You write the comparison yourself, every time:**
`if got != want { t.Errorf(...) }`

<!--
speaker_note: |
  Say this directly and early - some people in the room WILL go
  looking for the built-in assertEquals later in the course and need
  to hear right now that it genuinely doesn't exist, not that they
  missed an import.
-->

<!-- end_slide -->

## testify exists — and it's not stdlib

In real Go codebases, a huge number pull in
`github.com/stretchr/testify` for its `assert`/`require` packages. It's
close to a de facto standard.

<!-- pause -->

**Say this plainly: testify is a third-party addition. It ships from
nowhere near the standard library.**

```go
// with testify - a dependency you add on purpose
assert.Equal(t, 90.0, got)

// vanilla testing.T - what you get for free
if got != want {
    t.Errorf("got %v, want %v", got, want)
}
```

<!--
speaker_note: |
  Name testify explicitly and early, don't wait for someone to ask.
  The framing that matters: it's popular, it's common in production
  code, AND it's an opt-in dependency, not something the language
  ships. All three of those are true at once.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Where Go disagrees with what you know
===

<!-- end_slide -->

## Testing for panics: no assertPanics

There's no `assertPanics(t, func)` anywhere in the standard library.
You reach for `defer`/`recover` — directly inside the test.

```go
func TestDivideByZeroPanics(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Error("expected a panic, but didn't get one")
        }
    }()
    Divide(10, 0)
}
```

<!-- pause -->

**These are the same primitives from Topics 2 and 5** — deferred
execution and recover, doing real work here rather than sitting in a
slide about error handling.

<!--
speaker_note: |
  Worth naming explicitly: this isn't a new mechanism invented for
  testing. It's defer and recover, which they've already used, now
  applied to a new job. The safety-net-for-a-stunt-rehearsal framing
  works well here - you're deliberately trying to cause the fall,
  purely to confirm it happens when it's supposed to.
-->

<!-- end_slide -->

## Prove the panic test actually checks something

A panic test that can't fail isn't testing anything — it's decoration.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// Divide panics on b == 0
func Divide(a, b float64) float64 {
    if b == 0 {
        panic("cannot divide by zero")
    }
    return a / b
}
```

<!-- column: 1 -->

```
Remove the panic.
Re-run the test.

If it still passes,
your test wasn't
checking the panic -
it was passing by
accident.
```

<!-- reset_layout -->

**Demo:** comment out the `panic(...)` line, run `go test`, and watch
`TestDivideByZeroPanics` correctly fail.

<!--
speaker_note: |
  This is the step people skip when writing panic tests in the real
  world - they write the recover scaffolding, see it pass once, and
  never confirm it can fail. Make the failure happen live.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Table-driven tests
===

<!-- end_slide -->

## Not a framework feature — just Go

JUnit has `@ParameterizedTest`. pytest has `@pytest.mark.parametrize`.
Dedicated, special test-framework syntax in both cases.

<!-- pause -->

Go's answer: **a slice of struct literals and a `for` loop.** Nothing
new to learn — you already know how to write both.

```go
cases := []struct {
    name  string
    price float64
    want  float64
}{
    {"ten percent off", 100, 90},
    {"zero price stays zero", 0, 0},
}
```

<!--
speaker_note: |
  Slow down here specifically - this is one of the highlight moments
  of the whole course, not a routine technique to move past quickly.
  The point isn't "here's a testing trick," it's "you already had
  every tool needed for this before today's session started."
-->

<!-- end_slide -->

## `t.Run` wraps each case as its own subtest

```go
func TestDiscounts(t *testing.T) {
    cases := []struct{ name string; price, want float64 }{
        {"ten percent off", 100, 90},
        {"zero price stays zero", 0, 0},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := TenPercentOff(tc.price)
            if got != tc.want {
                t.Errorf("got %v, want %v", got, tc.want)
            }
        })
    }
}
```

<!-- pause -->

Each case reports pass/fail **independently** — one bad case doesn't
hide the rest.

<!-- end_slide -->

## Target one case by name

```
go test -run TestDiscounts/zero_price_stays_zero
```

<!-- pause -->

**Type in chat: in pytest or JUnit, what did isolating a single
parameterized case to debug it actually involve?**

<!-- pause -->

Here it's a command-line flag against the same table you already
wrote — no separate debug configuration, no commenting out other
cases.

<!--
speaker_note: |
  Real-world angle if it helps land the point: a table-driven test is
  a quality-control checklist with pre-filled sample tickets run
  through the exact same inspection station, one at a time. The
  inspection process (the test body) never changes - only the data on
  each ticket does. That's the entire reason to pull cases into a
  table instead of writing one near-identical test per case.
-->

<!-- end_slide -->

## Coverage: built into the toolchain

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```
go test -cover ./...
```

Prints a coverage percentage. No JaCoCo (Java), no `coverage.py`
(Python) to install separately.

<!-- column: 1 -->

```
go test \
  -coverprofile=cover.out ./...

go tool cover \
  -html=cover.out
```

Opens a line-by-line HTML report — covered lines green, uncovered
lines red.

<!-- reset_layout -->

**Demo:** run the HTML report against a package with an untested
branch and find the red line together.

<!--
speaker_note: |
  Worth naming that coverage tells you what RAN, not whether the
  assertions were meaningful - 100% coverage with a weak assertion
  still shows green. Don't let the room treat the percentage as a
  quality score on its own.
-->

<!-- end_slide -->

## Benchmarking: same package, same command

```go
func BenchmarkTenPercentOff(b *testing.B) {
    for i := 0; i < b.N; i++ {
        TenPercentOff(100)
    }
}
```

```
go test -bench=. -benchmem
```

<!-- pause -->

No JMH-style separate setup — `go test` runs it.

<!-- end_slide -->

## `b.N` isn't a number you choose

`go test` runs your benchmark function repeatedly, adjusting `b.N`
across runs, until it has a **stable timing measurement**.

<!-- pause -->

Think of a stopwatch operator deciding to "time the line for at least
a second's worth of widgets" — however many that turns out to be —
rather than insisting on exactly 100 regardless of how fast the line
runs.

<!-- pause -->

`-benchmem` adds allocation counts (`B/op`, `allocs/op`) alongside the
`ns/op` timing — change the implementation, re-run, see if either
number moves.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **`testing` + `go test` ship with the toolchain** — no framework to install, choose, or configure
2. **No built-in assertion library** — four methods (`Error`/`Errorf`/`Fatal`/`Fatalf`), you write the comparison
3. **`testify` is common but is explicitly a third-party addition**, not stdlib
4. **Table-driven tests are just a slice of structs and a `for` loop** — `t.Run` gives per-case reporting and `-run` targeting
5. **Panic testing uses `defer`/`recover` by hand** — the same primitives from Topics 2 and 5
6. **Coverage and benchmarking are first-class, built into `go test`** — no separate JaCoCo/coverage.py/JMH equivalent
7. **None of this was new today** — you'd already written tests in ten straight topics; this session gave the habit table-driven cases, coverage, and benchmarks on top of what you had

<!-- end_slide -->

## Back to the opening scenario

You told us how long the JUnit-vs-TestNG or pytest-vs-unittest
conversation actually took — a meeting, a wiki page, a config PR
someone had to review.

<!-- pause -->

**Type in chat: now that you've seen `testing` and `go test`, how long
does that same decision take in Go?**

<!-- pause -->

There is no decision. It ships with the toolchain, day one, for every
Go project you'll ever start.

<!--
speaker_note: |
  This is the payoff for the opening chat poll - read a couple of
  the original answers back if you can. Let "zero" land as a genuine
  relief rather than rushing to the bridge slide. Don't turn it into
  a Go-vs-everything-else superiority debate - the point is the
  decision doesn't exist, not that other ecosystems are wrong to have
  one.
-->

<!-- end_slide -->

## Bridge to Topic 13

**We've established:**

<!-- incremental_lists: true -->

- Go's `testing` package and `go test` command need no separate setup — ever
- Table-driven tests turn a JUnit/pytest framework feature into ordinary Go you already knew
- Coverage and benchmarking are first-class parts of the same tool, not add-ons
- The habit itself started back in Topic 2 — today just gave it sharper tools

<!-- incremental_lists: false -->

You've now built, tested, and can run every piece of Go code from this
course — with a `_test.go` file behind most of it, not just this
topic's. **Topic 13: Docker & Deployment** — packaging that single
compiled binary from Topic 1 into a container, and getting it running
somewhere that isn't your laptop.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
