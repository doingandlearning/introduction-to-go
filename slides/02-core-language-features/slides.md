---
title: "**Core Language Features**"
sub_title: Go Programming — Topic 2
author: Kevin Cunningham
---

## Opening scenario

You're porting a small utility from Java. It has two overloaded methods:

```java
int add(int a, int b) { return a + b; }
double add(double a, double b) { return a + b; }
```

You write the Go version the way muscle memory tells you to — two
functions, both called `Add`, different parameter types.

<!-- pause -->

It doesn't compile. Not "compiles with a warning." **Doesn't compile.**

**Type in chat: if you can't overload a function name, what do you think Go expects you to do instead?**

<!--
speaker_note: |
  Let a few guesses land - "rename it," "use an interface," "generics"
  are all reasonable and all technically usable answers. Don't confirm
  or deny yet, just bank the guesses. We resolve this once we've covered
  functions properly.
-->

<!-- end_slide -->

## Three ways to declare a variable

```go
var x int = 5   // explicit type
var y = 5       // inferred type
z := 5          // short declaration — function-scope only
const Pi = 3.14159
```

<!-- pause -->

`:=` is the one you'll type constantly. It only works inside a function
body — package-level declarations need `var` or `const`.

<!--
speaker_note: |
  Point out that this isn't three interchangeable styles you pick by
  taste. Idiomatic Go leans hard on := inside functions and var at
  package level, largely because := isn't legal at package level at all.
-->

<!-- end_slide -->

## The type menagerie

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```
int8   int16   int32   int64
uint8  uint16  uint32  uint64
float32  float64
int   uint   (platform width)
string   bool
```

<!-- column: 1 -->

`byte` is just an alias for `uint8`.

`rune` is just an alias for `int32` — a single Unicode code point, not a
byte.

If you've used `char` in Java, a Go `rune` is the closer equivalent, not
`byte`.

<!-- reset_layout -->

<!-- end_slide -->

## No implicit numeric conversion. Ever.

```go
var i int = 5
var f float64 = 3.2

result := i + f
// compile error: mismatched types int and float64
```

<!-- pause -->

Java and JavaScript would silently widen `i` to a `float`/`float64` and
let this run. Go refuses.

```go
result := float64(i) + f   // this is the only way
```

**Every conversion is explicit and visible at the call site — no silent precision loss hiding in an arithmetic expression.**

<!--
speaker_note: |
  This is a strong candidate for the opening-scenario bug if you want a
  live demo instead of/alongside the overloading one - a Python or JS
  developer will absolutely have shipped an int-plus-float expression
  without a second thought, and Go stops them cold at compile time.
-->

<!-- end_slide -->

## Zero values, not null

Every declared-but-unassigned variable gets a deterministic **zero
value** — never `null`, never `undefined`, never garbage memory.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
var i int       // 0
var s string    // ""
var b bool      // false
var p *int      // nil
```

<!-- column: 1 -->

```go
type Point struct{ X, Y int}

var pt Point
// {0 0} — usable immediately,
// no constructor required
```

<!-- reset_layout -->

**Demo:** declare one variable of each basic type, unset, and print each — compare the zero values against what Python's `None` or Java's `null` would give you.

<!--
speaker_note: |
  The struct example is the one that lands hardest for Java devs -
  there's no "did you forget to `new` this" NullPointerException class
  of bug for a zero-value struct, because it was never null to begin
  with. It's a real, usable value from the moment it's declared.
-->

<!-- end_slide -->

## No `enum` keyword — build one with `iota`

Java, TypeScript, and Python (3.4+) all have a real `enum` keyword. Go
doesn't. The idiomatic replacement is a named type over `int`, plus
`iota` to number the constants for you.

```go
type Status int

const (
	StatusPending Status = iota // 0
	StatusActive                // 1
	StatusDone                  // 2
)
```

<!-- pause -->

`iota` resets to `0` at the start of each `const` block and increments by
one per line inside it — you're not assigning `0`, `1`, `2` by hand, and
inserting a new value in the middle renumbers everything after it
automatically.

<!--
speaker_note: |
  If someone asks "why not just use string constants instead," that's a
  reasonable alternative for logging/display, but StatusPending Status =
  iota gives you a distinct TYPE (Status), not just a string - a function
  that takes a Status can't accidentally be called with an unrelated int
  or a typo'd string. That type safety is the actual point.
-->

<!-- end_slide -->

## The gotcha: it's still just an `int`

```go
var s Status = 99 // compiles fine — nothing restricts this to 0, 1, or 2

fmt.Println(s) // prints "99", not a name
```

<!-- pause -->

Unlike a real enum in Java or TypeScript, Go's compiler does **not**
restrict a `Status` variable to the values you declared. Nothing stops an
out-of-range value from compiling, and printing one with `%v` just shows
the underlying number — not a readable name.

**Demo:** declare `Status(99)`, print it, and compare against `StatusDone`
printed the same way — both come out as plain integers with no hint
they're "supposed" to be a small, closed set.

<!-- pause -->

The fix for the printing half of this is the `Stringer` interface, coming
up shortly in this same topic — implement one `String()` method and
`fmt` starts printing names instead of numbers automatically.

<!--
speaker_note: |
  Flag this as a real production gotcha, not a theoretical one: a
  function signature like func Apply(s Status) reads like it's closed to
  three values, but nothing enforces that at the type level the way a
  Java enum or a TypeScript union type would. Validate at the boundary
  (JSON decode, CLI parsing) if an out-of-range value would actually be a
  problem.
-->

<!-- end_slide -->

## Functions: multiple returns, named returns, variadics

```go
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
```

<!-- pause -->

Multiple return values aren't a tuple hack — they're how Go expresses
"this can fail" without exceptions or a wrapper type. More on that next.

<!-- end_slide -->

## No function overloading

You cannot declare `Add(a, b int) int` and `Add(a, b float64) float64` in
the same package. One name, one signature, full stop.

<!-- incremental_lists: true -->

- Give the functions different names: `AddInts`, `AddFloats`
- Accept an interface and type-switch inside
- Use generics — covered properly in Topic 6

<!-- incremental_lists: false -->

**This is usually the single biggest adjustment for a Java developer in this whole course.**

<!--
speaker_note: |
  Resolve the opening scenario's guesses here if you haven't already -
  most rooms guess at least one of the three options above. Generics
  exist and do solve this, but don't get pulled into a Topic 6 detour;
  name it and move on.
-->

<!-- end_slide -->

## No exceptions for ordinary failure

There is no `try`/`catch` in Go. A function that can fail returns
`error` as its last value, and the caller checks it immediately.

```go
result, err := divide(10, 0)
if err != nil {
	log.Fatal(err)
}
```

<!-- pause -->

**`if err != nil` is not boilerplate you'll eventually abstract away — it's the defining idiom of the language, written deliberately, every time.**

`panic`/`recover` exist, but they're reserved for truly exceptional
situations — a programming bug, not "the user typed bad input."

<!-- pause -->

One more piece of function-execution behavior before we look at the
failing test already waiting for this code: what actually runs when a
function returns — including an early return, or a panic.

<!--
speaker_note: |
  Someone will ask "isn't this just try/catch with extra steps?" The
  honest answer: it trades brevity for the compiler forcing you to look
  at every failure point instead of letting one catch block at the top
  of a call stack silently swallow five different failure modes.
-->

<!-- end_slide -->

## `defer`: run this when the function returns

```go
func readConfig() error {
	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}
```

<!-- pause -->

`defer` schedules the call to run when the surrounding function returns —
however it returns: normal return, early return, or even a panic.

**Real-world version:** clipping your car key to your house key hook the
moment you walk in the door, so future-you physically cannot leave the
house without both. You set up the guarantee right next to the point
where you created the thing that needs cleaning up — not in a `finally`
block somewhere far below.

<!-- end_slide -->

## The gotcha: arguments evaluate now, the call runs later

```go
for i := 0; i < 3; i++ {
	defer fmt.Println(i)
}
// prints: 2, 1, 0
```

<!-- pause -->

Two rules collide here:

<!-- incremental_lists: true -->

- `defer`'s **arguments** are evaluated immediately, when the `defer` statement runs
- the **call itself** executes later, in LIFO order (last deferred, first run)

<!-- incremental_lists: false -->

**Demo:** run this loop, then predict — before running it — what a `defer func(n int) { fmt.Println(n) }(i)` version would print instead.

That's the last new language rule before we look at the failing test
already sitting in this topic's lab.

<!--
speaker_note: |
  Walk through why it's 2, 1, 0 rather than 0, 1, 2: each defer captures
  the value of i at the moment defer ran (2, then 1, then 0 as the loop
  counted down through its final iterations), and LIFO unwinds them in
  reverse of the order they were deferred. The func(n int){...}(i) fix
  works because passing i as an argument copies it at defer-time, same
  as the plain fmt.Println(i) case - the real gotcha shows up when
  people close over the loop variable in a closure instead.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Proving it: your first failing test
===

<!-- end_slide -->

## From here on, every exercise starts red

Every lab for the rest of this course ships with its tests **already
written** — sitting in `starter/` right next to the code you're about to
implement. Run `go test ./...` before you write a single line, and they
fail. That's not a bug in the starter code. That's the starting line.

<!-- pause -->

Starting now, not in Topic 12. Your job is never "write a test to prove
your code works" — it's "make the test that's already there stop
failing." The test is the spec, stated precisely enough for a machine to
check, not something bolted on after you've already convinced yourself.

**How to *write* one — `t.Errorf`, subtests, table-driven cases — is
deliberately Topic 12's job**, saved for when you have a whole course's
worth of code to test in hindsight. Until then, you'll read plenty of
test code. You won't write any.

<!--
speaker_note: |
  Name John Arundel-style thinking here if it's useful to the room: a
  test isn't just verification after the fact, it's the specification
  you're implementing against. Landing this distinction early - "make it
  pass" now, "write one" in Topic 12 - avoids students feeling like
  they're expected to already know test syntax nobody's taught them yet.
-->

<!-- end_slide -->

## Reading a test before you can pass it

```go
// divide_test.go — already sitting in starter/, before divide exists
package main

import "testing"

func TestDivide(t *testing.T) {
	got, err := divide(10, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 5 {
		t.Errorf("divide(10, 2) = %v, want 5", got)
	}
}
```

```
$ go test ./...
./errdemo_test.go:8:9: undefined: divide
FAIL
```

<!-- pause -->

<!-- incremental_lists: true -->

- `TestX(t *testing.T)` — the function the test runner calls; you'll recognize the shape without writing it yet
- `t.Fatalf`/`t.Errorf` — what the test prints when it doesn't get what it expects
- The failure above is the honest starting state — nothing to fix in the test, only in the code it's waiting on

<!-- incremental_lists: false -->

**Demo:** run `go test ./...` against this topic's starter code, read the failure, then implement just enough of `divide` to turn it green.

<!-- pause -->

One more thing before this topic closes: you've been calling
`fmt.Println` since Topic 1 without a second thought. Time to actually
learn what's in that package.

<!--
speaker_note: |
  This is deliberately the bare minimum reading of a test, not the full
  testing package - table-driven tests, testify, coverage, and
  benchmarking are all Topic 12's job. The point here is just: a failing
  test is a to-do list, not a verdict, and reading one is a skill on its
  own before writing one ever comes up.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

fmt: a proper primer
===

<!-- end_slide -->

## The three basic printers

```go
fmt.Print("no newline, ")
fmt.Print("just concatenation")

fmt.Println("adds a newline")
fmt.Println("and a space between args", 42)

fmt.Printf("formatted: %s is %d years old\n", "Go", 16)
```

<!-- pause -->

`Print` and `Println` insert spaces between operands only when neither
side is already a string. `Printf` gives you no automatic spacing or
newline at all — the format string controls everything.

**This is the one that trips people up: forget the `\n` in `Printf` and your output just runs together.**

<!-- end_slide -->

## The core Printf verbs

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```
%v    default format
%+v   adds field names
%#v   Go-syntax literal
%T    the type itself
```

<!-- column: 1 -->

```
%d    integer
%s    string
%q    quoted string
%f    float
%.2f  float, 2 decimal places
%t    bool
%p    pointer
```

<!-- reset_layout -->

You will use `%v`, `%d`, `%s`, and `%.2f` constantly. The rest earn their
keep in debugging sessions.

<!-- end_slide -->

## Compare all four, side by side

```go
type Point struct{ X, Y int }
p := Point{X: 3, Y: 4}

fmt.Printf("%v\n", p)   // {3 4}
fmt.Printf("%+v\n", p)  // {X:3 Y:4}
fmt.Printf("%#v\n", p)  // main.Point{X:3, Y:4}
fmt.Printf("%T\n", p)   // main.Point
```

<!-- pause -->

<!-- incremental_lists: true -->

- `%v` — quick glance, what's roughly in here
- `%+v` — debugging, which field is which
- `%#v` — copy-pasteable back into Go source, a literal you could compile
- `%T` — "what type even is this," useful with `interface{}`/`any` values

<!-- incremental_lists: false -->

**Demo:** run `cmd/fmtdemo` and match each line of output to the verb that produced it.

<!--
speaker_note: |
  Have delegates guess the output before running it, especially %#v -
  the "main.Point{X:3, Y:4}" form surprises people the first time
  because it includes the package-qualified type name, not just the
  struct name.
-->

<!-- end_slide -->

## Formatting into a string, an error, or anywhere else

```go
msg := fmt.Sprintf("user %d not found", id)

err := fmt.Errorf("loading config: %w", innerErr)
```

<!-- pause -->

`Sprintf` builds a string instead of printing one — use it wherever you'd
otherwise be gluing strings together by hand. `Errorf` does the same for
`error` values, and `%w` wraps an inner error so it can be unwrapped
later (`errors.Is` / `errors.As` — a Topic 3+ concern, filed away for
now).

<!-- end_slide -->

## `fmt` doesn't just print to the terminal

```go
fmt.Fprintln(os.Stderr, "warning: cache miss")

var buf bytes.Buffer
fmt.Fprintf(&buf, "%d items processed", count)
```

<!-- pause -->

`Fprintf`/`Fprintln` write to anything satisfying `io.Writer` — standard
error, an open file, a network connection, an in-memory buffer used in a
test. `Println`/`Printf` are really just `Fprintln`/`Fprintf` pinned to
`os.Stdout`.

**This is the generalization that makes `fmt` work everywhere, not a
stdout-only convenience function.**

<!-- end_slide -->

## Controlling your own `%v`: the `Stringer` interface

```go
type Point struct{ X, Y int }

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

fmt.Println(Point{3, 4})   // (3, 4) — not {3 4}
```

<!-- pause -->

Implement one method, `String() string`, and every `%v` and `Println`
call involving that type uses it automatically. `fmt` checks for this
interface before falling back to its default struct formatting.

**One small method, and your debug output stops looking like raw memory
dumps and starts looking like something you'd choose to read.**

<!-- pause -->

**This is the fix for the `Status(99)` problem a few slides back.** Give
`Status` its own `String() string` and it prints `"Pending"` instead of
`0` — `fmt` finds it the same way it just found `Point`'s.

The `(p Point)` between `func` and the name is called a **receiver** —
it's what turns a plain function into something you can call as
`p.String()`. Full treatment is Topic 4's job; for now, just recognize
the shape when you see it.

<!--
speaker_note: |
  Two callbacks in one slide - the enum gotcha from earlier in this
  topic, and a forward flag for the receiver syntax nobody's formally
  taught yet. Don't get pulled into a full methods digression here; name
  it, connect it, and let Topic 4 do the real work.
-->

<!-- end_slide -->

## Organizing packages

A package is a **directory**. Every `.go` file inside it declares the
same `package name` at the top.

<!-- pause -->

<!-- incremental_lists: true -->

- Files in the same package share a namespace — no imports needed between them (same rule as Topic 1's `farewell.go`)
- Other packages import this one by its **import path**, usually rooted at the module name in `go.mod`
- Import happens at package granularity, not per-file — you can't import "just one file"

<!-- incremental_lists: false -->

**Demo:** open `mathutils` in this topic's lab and trace the import path from `go.mod` through to the `import` line in `main.go`.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Three declaration styles**: `var x int = 5`, `var y = 5`, `z := 5` (function-scope only) — plus `const`
2. **No implicit numeric conversion**: every widening or narrowing is an explicit `float64(i)`-style call
3. **Zero values replace null**: every type has a deterministic, immediately usable default
4. **No `enum` keyword**: `iota` inside a `const` block numbers a named type for you, but the compiler never restricts a variable to just those values
5. **No function overloading**: rename, use an interface, or reach for generics later
6. **No exceptions for ordinary errors**: multiple returns plus `if err != nil`, everywhere, on purpose
8. **`fmt` generalizes past stdout**: `Sprintf`, `Errorf`, `Fprintf` against any `io.Writer`, and `Stringer` for custom formatting
9. **A package is a directory**: same-package files share a namespace for free; imports work at package granularity
10. **Every lab from here on starts red**: the test is already written in `starter/` — your job is to make it pass, not to write it; writing one yourself is Topic 12's job

<!-- end_slide -->

## Back to the opening scenario

You tried to write two `Add` functions and Go refused to compile them.

**Now you know why, and what to do instead:** rename them, accept an
interface, or wait for generics in Topic 6. None of those is a
workaround — they're the idiomatic answer.

<!-- pause -->

**Type in chat: which of the three fixes would you reach for first in the utility you were porting — and why?**

<!--
speaker_note: |
  There's no single right answer here - renaming is often genuinely
  clearer than an overload would have been, since AddInts vs AddFloats
  documents intent that Add/Add never could. Let that reframe land: the
  restriction isn't purely a limitation, it sometimes produces more
  readable call sites.
-->

<!-- end_slide -->

## Bridge to Topic 3

**We've established:**

<!-- incremental_lists: true -->

- Types convert explicitly, never silently
- Errors are values, checked immediately, not caught after the fact
- `defer` guarantees cleanup runs, in LIFO order, right next to the resource that needs it
- Every lab from here on starts with a failing test — implement until `go test ./...` passes; writing tests yourself is Topic 12's job
- `fmt` is a small toolkit that works across strings, errors, and any `io.Writer`

<!-- incremental_lists: false -->

**Topic 3: Flow Control and Data Structures** — `if`/`for`/`switch`,
arrays, slices, maps, and the range loop that ties them together. The
lab's last step will be a test again, same as this one.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
