# Topic 3 sample code — Flow Control & Data Structures

Four small programs, one module. Run everything from this directory
(`code/`).

## `cmd/switchdemo`

Grades a set of scores with a conditionless `switch` (`switch { case
score >= 90: ... }`), then walks through a day-of-week `switch` that uses
`fallthrough` deliberately to group Saturday's hours into Sunday's case.

```
go run ./cmd/switchdemo
```

Watch day 6 (Saturday) in the output — it hits its own case, then falls
into Sunday's case body because of the explicit `fallthrough`. Comment
that line out and re-run to see the difference a single keyword makes.

## `cmd/slicealias`

The big one. Reproduces the slice-aliasing bug step by step: takes a
sub-slice, mutates through it, then appends twice — once within capacity
(silently overwrites the original), once past capacity (silently
reallocates and breaks the sharing).

```
go run ./cmd/slicealias
```

Watch the `len`/`cap` printed at every step. The moment `cap` jumps to a
bigger number than the line before it, that's the reallocation — mutating
through `view` after that point no longer touches `original`. This is the
single most common real bug pattern new Go developers write; read every
line of output before moving on.

## `cmd/nilmap`

Shows the nil-slice-vs-nil-map asymmetry in one program: a nil slice is
safe to read and safe to `append` to, a nil map is safe to read but
panics the instant you write to it.

```
go run ./cmd/nilmap
```

The program is written to panic on purpose — that's the point. Read the
panic message (`assignment to entry in nil map`), then open the source
and find the two commented-out fixes (`make(map[string]int)` or a map
literal).

## `cmd/wordcount`

Builds a `map[string]int` word-frequency counter over a sentence using
`strings.Fields`, uses the comma-ok idiom to tell "appeared zero times"
apart from "never seen," then ranges over the same map twice in a row.

```
go run ./cmd/wordcount
```

Compare the two `range` passes printed in the same run, then run the
whole program again as a separate execution and compare again. The order
is deliberately randomized by Go — don't expect it to match Python's
insertion-ordered dicts or Java's `LinkedHashMap`.
