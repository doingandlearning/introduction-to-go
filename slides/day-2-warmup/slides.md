---
title: "**Day 2 Warm-up**"
sub_title: Six things from Day 1 that compiled clean and still bit someone
author: Kevin Cunningham
---

## Before we touch Topic 4

Six short snippets. Six questions: **what does this print, or why won't
it compile?**

Guess in chat before each reveal — nothing new here, just yesterday's
material with the volume turned up.

<!--
speaker_note: |
  Pure recap/energy-check opener before diving into OOP, not new
  teaching. Don't let any one round run past 2-3 minutes - the point is
  quick recall and a bit of fun, not reteaching Topics 1-3. Move on the
  moment the room has the answer, even if the explanation could go
  deeper.
-->

<!-- end_slide -->

## Round 1 — Topic 1: what's wrong here?

```go
package main

import "fmt"

func main()
{
    fmt.Println("nope")
}
```

**Type in chat: why won't this compile?**

<!-- pause -->

Go's lexer inserts a semicolon at the end of `func main()` because the
line ends there — the opening brace on its own line has nothing left to
attach to. The compiler error won't even say "brace placement"; it just
complains the function body is missing.

**The rule:** opening braces stay on the same line as the statement they
belong to. Always.

<!-- end_slide -->

## Round 2 — Topic 2: what prints?

```go
var i int = 5
var f float64 = 3.2

result := i + f
fmt.Println(result)
```

**Type in chat: what does this print?**

<!-- pause -->

Nothing — **it doesn't compile.** `mismatched types int and float64`.
Java and JavaScript would silently widen `i` to a float and let this
run. Go refuses; every conversion has to be spelled out:
`float64(i) + f`.

<!-- end_slide -->

## Round 3 — Topic 2: what prints?

```go
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusDone
)

s := Status(99)
fmt.Println(s)
```

**Type in chat: what does this print?**

<!-- pause -->

`99`. `Status` is still just an `int` under the hood — nothing stops an
out-of-range value from compiling, and without a `String()` method,
`fmt` has no name to print, only the number.

<!-- end_slide -->

## Round 4 — Topic 3: does this run?

```go
var m map[string]int
m["x"] = 1
```

**Type in chat: does this compile? If it compiles, what happens when it
runs?**

<!-- pause -->

Compiles fine. Then **panics** — `assignment to entry in nil map`.

A nil map is perfectly safe to *read*: `m["x"]` would just hand back the
zero value, `0`. Writing to one is what panics. Compare that to
`var s []int; s = append(s, 1)` — completely fine. Nil slice and nil map
are not the same nil.

<!-- end_slide -->

## Round 5 — Topic 3: what prints?

```go
switch day := 6; day {
case 1, 2, 3, 4, 5:
    fmt.Println("weekday")
case 6:
    fmt.Println("saturday")
case 7:
    fmt.Println("sunday")
}
```

**Type in chat: just "saturday," or "saturday" and "sunday" both?**

<!-- pause -->

Just `saturday`. Go's `switch` doesn't fall through by default — the
*opposite* default from C, Java, and JavaScript. Reaching `case 7` too
would need an explicit `fallthrough`.

<!-- end_slide -->

## Round 6 — Topic 3: the one that ties it together

```go
func addItem(items []string, item string) {
    items = append(items, item)
}

list := []string{"a", "b"}
addItem(list, "c")
fmt.Println(list)
```

**Type in chat: does `list` end up with "c" in it, or not?**

<!-- pause -->

No `"c"`. Still `[a b]`. A slice is a small header — pointer, length,
capacity — and Go passes that header **by value**, exactly like a
struct. `addItem` reassigns its own local copy of the header; nothing
about that reaches back to `list` in `main`.

<!-- pause -->

**This is the same value-vs-pointer question Topic 4 is about to spend
an entire topic on** — methods, receivers, and exactly when Go copies
something versus when it shares it.

<!--
speaker_note: |
  Land this one deliberately as the bridge, not just another gotcha. If
  someone answers "yes, it has c," compare it to yesterday's
  original[1:3] aliasing example: same header-sharing mechanism behind
  both, opposite outcome here, because addItem only ever has a COPY of
  the header - nothing aliases list's header at all. That contrast is
  the whole point of the round.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

On to Topic 4
===

<!-- end_slide -->

## Into Object-Oriented Programming

Every one of those six gotchas traces back to one thing: Go being
explicit about what's copied, what's shared, and what's allowed to fail
silently versus loudly.

**Topic 4 puts a name on the value-vs-pointer question Round 6 just
raised** — methods, receivers, and how Go builds behavior around data
without a `class` keyword.

<!-- end_slide -->
