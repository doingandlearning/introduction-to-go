---
title: "**Functional Programming**"
sub_title: Go Programming — Topic 6
author: Kevin Cunningham
---

## Opening scenario

You need a constructor for a `Server` type. It has 6 optional settings —
timeout, retries, TLS config, logger, max connections, idle interval.
Most callers only ever set 2 of them.

In Python you'd give the function default keyword arguments. In Java
you'd overload the constructor five different ways, or reach for a
Builder class.

**Type in chat: Go has no default parameter values and no function
overloading at all. What does that constructor's signature look like?**

We'll come back to this once we've seen what Go does with functions as
values.

<!--
speaker_note: |
  Let a few guesses land - people usually propose a giant struct-of-config
  argument, or "just make everything a pointer and check for nil," or
  simply "I don't know, that seems like a real gap." All are reasonable.
  Bank the guesses, don't resolve yet.
-->

<!-- end_slide -->

## Functions are values, full stop

In Go, a function is a value with a type, exactly like an `int` or a
`string`. You can assign it to a variable, pass it as an argument, return
it from another function, or store it in a struct field.

```go
var op func(int, int) int = func(a, b int) int { return a + b }
op(2, 3) // 5
```

<!-- pause -->

Python and JavaScript developers already expect this. If you're coming
from Java, this is closer to a lambda assigned to a functional
interface — except there's no interface to declare first.

<!-- end_slide -->

## Closures capture their surroundings

A function literal defined inside another function keeps a live
reference to the variables around it — even after the outer function has
returned.

```go
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

c := makeCounter()
c() // 1
c() // 2
c() // 3
```

<!-- pause -->

**Each call to `makeCounter` creates a brand-new `count`.** Two counters
from two calls never see each other's state.

<!--
speaker_note: |
  Demo this live: create two counters, interleave calls to each, show the
  sequences stay independent. This is the moment to head off "wait, is
  count a global?" - it isn't, and proving it with a second counter is
  more convincing than saying so.
-->

<!-- end_slide -->

## Closures are delegation, not magic

Think of a manager in a small office who doesn't do the labor personally.
The manager holds onto "whichever worker function got handed to them,"
and calls it when the moment is right.

<!-- incremental_lists: true -->

- The manager doesn't care who wrote the worker function
- The manager doesn't care what the worker function does internally
- The manager only cares that it has the right shape — the right
  signature

<!-- incremental_lists: false -->

That's every higher-order function in Go: something that takes a
function as a parameter, or hands one back, without needing to know
what's inside it.

<!-- end_slide -->

## Higher-order functions, made reusable with generics

Before Go 1.18 (2022), writing a reusable `Filter` meant `interface{}`
and type assertions, or a separate function per type. Generics remove
that tax.

```go
func Filter[T any](items []T, predicate func(T) bool) []T {
    var result []T
    for _, item := range items {
        if predicate(item) {
            result = append(result, item)
        }
    }
    return result
}
```

<!-- pause -->

```go
evens := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
// [2 4 6]

long := Filter([]string{"go", "rust", "c", "python"}, func(s string) bool { return len(s) > 3 })
// [rust python]
```

**One function definition, two completely different element types.**

<!-- end_slide -->

## If you're expecting generics everywhere old blog posts aren't

Generics landed in **Go 1.18, mid-2022** — genuinely recent as language
features go.

<!-- pause -->

Anything written before that date, and a lot written since out of habit,
solves this problem a different way:

<!-- incremental_lists: true -->

- `interface{}` (now spelled `any`) plus a type switch or type assertion
- Code generation — a tool that writes a typed version of your function
  per type, before compilation
- Reflection, via the `reflect` package — slow, and a last resort

<!-- incremental_lists: false -->

If you land on a Stack Overflow answer or a blog post doing any of this
for a "generic-looking" problem, check the date before you copy it.

<!-- end_slide -->

<!-- jump_to_middle -->

Go is still a loop-first language
===

<!-- end_slide -->

## No Map, no Filter, no Reduce in the standard library

Python has comprehensions built into the syntax. JavaScript has
`Array.prototype.map/filter/reduce` on every array. Go's `slices`
package (1.21+) is mostly sort, search, and compare helpers — it is
**not** a Map/Filter/Reduce toolkit.

<!-- pause -->

This isn't a gap Go forgot to close. It's a stance.

**This is deliberately how idiomatic Go stays.** For "transform every
element," the idiomatic answer is usually just a `for` loop.

<!-- end_slide -->

## Why a loop beats a chain, to Go's culture

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// A combinator chain hides
// what happens on error
result := Map(Filter(items, ok), transform)
```

<!-- column: 1 -->

```go
// The loop keeps every step,
// and every error, visible
var result []Item
for _, it := range items {
    if !ok(it) {
        continue
    }
    v, err := transform(it)
    if err != nil {
        return nil, err
    }
    result = append(result, v)
}
```

<!-- reset_layout -->

**Type in chat: which version would you rather debug at 2am from a stack
trace with no line numbers for the lambda?** Experienced Go developers
generally pick the loop, and mean it.

<!--
speaker_note: |
  This is a genuine cultural belief in the Go community, not a strawman -
  worth saying plainly so delegates don't think Go "just doesn't have"
  map/filter/reduce out of an oversight. It has Filter and Map available
  via generics (previous slides), the community mostly still reaches for
  loops, especially once error handling enters the picture.
-->

<!-- end_slide -->

## No currying, no operator overloading

Go gives you no syntactic shortcut for partial application and no way to
overload an operator to compose functions.

```go
// Partial application - by hand, with a closure, every time
func add(a, b int) int { return a + b }

addFive := func(b int) int { return add(5, b) }
addFive(3) // 8
```

<!-- pause -->

There's no `add.curry()`, no `functools.partial` equivalent baked into
the language. If you want it, you write the wrapping closure yourself.

<!-- end_slide -->

<!-- jump_to_middle -->

The idiom that actually matters: functional options
===

<!-- end_slide -->

## The real problem: no named or default parameters

Go has neither. Every argument is positional, every call must supply
all of them (or use variadic args). That collides badly with a
constructor that has several optional settings.

<!-- pause -->

**This is the gap from the opening scenario.** In Java you'd write five
overloaded constructors. In Python you'd write default keyword args. In
Go, neither tool exists — so Go grew a different idiom entirely.

<!-- end_slide -->

## Functional options

A constructor takes a variadic slice of functions, each of which mutates
the thing being built. Every stdlib-adjacent library that needs
optional configuration uses this — gRPC, the AWS SDK, and more.

```go
type Server struct {
    timeout time.Duration
    retries int
}

type ServerOption func(*Server)

func WithTimeout(d time.Duration) ServerOption {
    return func(s *Server) { s.timeout = d }
}

func WithRetries(n int) ServerOption {
    return func(s *Server) { s.retries = n }
}
```

<!-- end_slide -->

## Functional options, built

```go
func NewServer(opts ...ServerOption) *Server {
    s := &Server{
        timeout: 30 * time.Second,
        retries: 3,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

<!-- pause -->

```go
NewServer()                                          // defaults: 30s, 3 retries
NewServer(WithTimeout(5 * time.Second))               // 5s, 3 retries
NewServer(WithTimeout(5*time.Second), WithRetries(10)) // 5s, 10 retries
```

**All three calls are valid. None of them needed a different function
name or an overload.**

<!--
speaker_note: |
  Walk through NewServer() first and have delegates predict the field
  values before you show them - the defaults surviving untouched when no
  option is passed is the part that needs to visibly land.
-->

<!-- end_slide -->

## Functional options are a coffee counter

Think of ordering a coffee. The base order already has sensible
defaults filled in — regular size, whole milk, one shot.

<!-- incremental_lists: true -->

- You hand over zero or more modifier slips: oat milk, extra shot,
  decaf
- Each slip independently knows how to adjust the order
- You never need a separately printed order form for every possible
  combination of milk, shots, and size

<!-- incremental_lists: false -->

That's `WithTimeout` and `WithRetries` — modifier slips for a struct
instead of a coffee.

<!-- end_slide -->

## What this replaces

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Java**, no functional options:

```java
new Server();
new Server(timeout);
new Server(timeout, retries);
// ...or a whole separate
// Builder class
```

<!-- column: 1 -->

**Go**, with functional options:

```go
NewServer()
NewServer(WithTimeout(t))
NewServer(WithTimeout(t), WithRetries(n))
```

<!-- reset_layout -->

One pattern resolves two gaps at once: no overloading, and no
named/default parameters. The Builder-class alternative is real, and
we'll meet it properly in Topic 9.

<!-- end_slide -->

## Method values: a function that remembers its receiver

You can assign a method — without calling it — to a variable. Go bundles
the receiver in with it automatically.

```go
type Invoice struct {
    Subtotal float64
    TaxRate  float64
}

func (i Invoice) Total() float64 {
    return i.Subtotal * (1 + i.TaxRate)
}

inv := Invoice{Subtotal: 100, TaxRate: 0.2}

getTotal := inv.Total // no call, no parens - a method value
getTotal()             // 120, still knows about inv
```

<!-- pause -->

**`getTotal` needed no reference to `inv` passed in again.** It captured
the receiver the moment it was assigned, the same way a closure captures
a variable.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Functions are values**: assign, pass, return, store them like any
   other type
2. **Closures capture live variables from the enclosing scope**: each
   call to the outer function gets its own independent state
3. **Generics (1.18+) make higher-order functions reusable**: one
   `Filter` or `Map` works across element types — but this is recent,
   pre-2022 code fakes it with `interface{}` or codegen
4. **No built-in Map/Filter/Reduce, and the culture prefers loops**: not
   a missing feature, a readability stance
5. **Functional options solve the real gap**: no named or default
   parameters, no overloading — a variadic list of `func(*T)` fixes both
6. **A test proves the defaults, not just a demo run**: the same
   `go test` mechanics from Topic 2 catch a functional-options default
   silently breaking, before a customer notices their "plain" order
   arrived as a large

<!-- end_slide -->

## Back to the opening scenario

You asked what a 6-optional-setting constructor looks like in Go.

**It's `NewServer(opts ...ServerOption)`**, with sensible defaults set
first and each option mutating the struct afterward.

<!-- pause -->

**Type in chat: compared to the Java overloads or Builder class you
were picturing earlier, does the functional options version feel like
more code or less, once you've written the options once?**

<!--
speaker_note: |
  The honest answer is "about the same amount of code, but it composes
  better" - the option functions are written once and combine freely,
  where overloaded constructors multiply combinatorially as settings are
  added. Let that land rather than oversell it.
-->

<!-- end_slide -->

## Bridge to Topic 7

**We've established:**

<!-- incremental_lists: true -->

- Functions are ordinary values, and closures capture live state
- Generics make higher-order functions reusable without giving up types
- Go still prefers an explicit loop over a combinator chain
- Functional options are how Go handles optional configuration without
  overloading or default parameters
- A short `go test` against `NewCoffeeOrder` — same mechanics from
  Topic 2 — is enough to catch a broken default before it ships

<!-- incremental_lists: false -->

**Topic 7: Concurrency** — goroutines, channels, and how Go's answer to
"do several things at once" is a language feature, not a library.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
