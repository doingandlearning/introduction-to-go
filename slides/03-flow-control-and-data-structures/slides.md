---
title: "**Flow Control & Data Structures**"
sub_title: Go Programming — Topic 3
author: Kevin Cunningham
---

## Opening scenario

You've got a slice of data. Before you modify it, you "back it up" the
obvious way:

```go
backup := original
backup[0] = 999
```

You check `original` afterward. **It changed too** — even though you never
touched `original` directly.

**Type in chat: if you didn't write to `original`, why did it change?**

We'll come back to this once we've taken slices apart properly.

<!--
speaker_note: |
  Let guesses land for 20-30 seconds. Expect "shouldn't that be a copy,
  like in Python?" from the Python crowd and confused silence from the
  Java crowd, who are used to arrays being references anyway but not to
  this specific slicing behavior. Don't resolve it yet - bank it.
-->

<!-- end_slide -->

## One keyword, every loop shape

Go has exactly one looping keyword: `for`. No `while`, no `do-while`, no
`foreach` keyword — `for` covers all of it.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// classic C-style
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

<!-- pause -->

```go
// while-style
for condition {
    // ...
}
```

<!-- pause -->

<!-- column: 1 -->

```go
// infinite, exit via break
for {
    if done {
        break
    }
}
```

<!-- pause -->

```go
// range-based
for i, v := range someSlice {
    // ...
}
```

<!-- reset_layout -->

If you've spent years reaching for `while` in Python or Java, that reflex
just becomes `for` with no init and no post statement.

<!-- end_slide -->

## if/else with a scoped extra

`if`/`else` works the way you'd expect, with one addition: an **initializer
clause**, scoped to the whole if/else chain.

```go
if err := doSomething(); err != nil {
    return err
}
// err does not exist out here
```

<!-- pause -->

This is idiomatic Go's favorite error-handling shape — you'll see it
everywhere, constantly. The variable it declares can't leak into the rest
of the function and get reused by accident later.

<!-- end_slide -->

## switch without a condition

Go's `switch` is more flexible than the C/Java version you know.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// C/Java: an if/else chain
if score >= 90 {
    grade = "A"
} else if score >= 80 {
    grade = "B"
} else {
    grade = "C"
}
```

<!-- pause -->

<!-- column: 1 -->

```go
// Go: switch with no condition
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
default:
    grade = "C"
}
```

<!-- reset_layout -->

Each `case` can also list multiple values on one line: `case 1, 2, 3:`.

<!-- end_slide -->

## The switch gotcha: no fallthrough by default

**This is the opposite default from C, Java, and JavaScript.**

<!-- pause -->

In those languages, execution falls into the next case unless you write
`break;`. In Go, every case breaks automatically — you opt **into** the old
behavior explicitly with the `fallthrough` keyword.

```go
switch day {
case 1, 2, 3, 4, 5:
    open = true
case 6:
    open = true // Saturday: open, but only mornings
    fallthrough
case 7:
    // falls into here from Saturday deliberately
}
```

**Demo:** write a switch with no `fallthrough` and one with it side by side, and watch how much execution differs from a single keyword.

<!--
speaker_note: |
  Anyone with break-habit muscle memory will write extra break statements
  out of pure reflex - harmless in Go, since break exits a switch here
  too, but pointing out they're now unnecessary lands well. The bigger
  risk is the opposite: assuming Go behaves like C and being surprised
  when a case doesn't fall through into the next one silently.
-->

<!-- end_slide -->

## Pointers: no arithmetic, same idea

`&x` takes the address of `x`. `*p` dereferences pointer `p` — reads or
writes the value it points to.

<!-- pause -->

**What's missing compared to C:** no pointer arithmetic. You cannot do
`p + 1` to walk forward through memory. A Go pointer points at exactly one
thing, always.

```go
x := 5
p := &x
*p = 10
fmt.Println(x) // 10
```

<!-- end_slide -->

## Structs: grouping named fields

```go
type Item struct {
    Name  string
    Count int
}

i := Item{Name: "Widget", Count: 3}
fmt.Println(i.Name, i.Count)
```

<!-- pause -->

No classes, no constructors required — a struct literal is enough. You'll
add methods to structs in the next topic; for now, think of a struct as a
plain, typed bag of fields.

<!-- end_slide -->

## Structs are values too

Passing a struct to a function passes a **copy**, exactly like passing an
`int` — unless you pass a pointer.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
func reset(i Item) {
    i.Count = 0
    // only the copy changes
}

reset(item)
// item.Count unchanged
```

<!-- column: 1 -->

```go
func resetPtr(i *Item) {
    i.Count = 0
    // the real thing changes
}

resetPtr(&item)
// item.Count is now 0
```

<!-- reset_layout -->

**Demo:** call both versions against the same struct and print `item.Count` after each call.

<!--
speaker_note: |
  This is the same value-vs-reference question delegates already have
  half-formed opinions about from other languages - Java passes object
  references by value, so mutating a field through them "just works,"
  which trains the wrong instinct for Go structs. Let that collision
  surface before explaining.
-->

<!-- end_slide -->

## Arrays: fixed size, baked into the type

```go
var a [5]int   // an array of exactly 5 ints
var b [6]int   // a completely different type from a
```

`[5]int` and `[6]int` are not interchangeable — the size is part of the
type, checked at compile time.

<!-- pause -->

**Arrays are value types**, copied on assignment and on function calls —
the reverse of Java, where arrays are always references:

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0]) // still 1 — independent copies
```

In practice this matters less than it sounds — idiomatic Go rarely uses
raw arrays directly.

<!-- end_slide -->

## Slices: what you actually use

A slice — `[]int`, no size in the brackets — is a flexible view over an
underlying array: a pointer, a length, and a capacity, bundled together.

```go
s := []int{1, 2, 3}       // slice literal
s = append(s, 4)           // grows as needed
fmt.Println(len(s), cap(s))
```

<!-- pause -->

95% of the time, this is the type you reach for. Arrays exist mostly as
the thing slices are built on top of.

<!-- end_slide -->

<!-- jump_to_middle -->

The single most common new-Go-dev bug
===

<!-- end_slide -->

## Slicing shares memory — it doesn't copy

```go
original := []int{1, 2, 3, 4, 5}
view := original[1:3]
```

`view` is **not** a new copy of elements 1 and 2. It's a new little struct
(pointer + length + capacity) pointing at the *same underlying array* as
`original`.

<!-- pause -->

```go
view[0] = 99
fmt.Println(original) // [1 99 3 4 5]
```

Nobody wrote `original[1] = 99`. It happened anyway, through `view`.

```go
type Slice {
    pointer memory_reference
    length int
    capacity int
}
```

<!--
speaker_note: |
  Say this plainly: re-slicing copies the header (pointer/len/cap), never
  the data. That single sentence is the whole mental model - everything
  that follows is a consequence of it.
-->

<!-- end_slide -->

## The real-world picture

A slice is a **viewfinder held over a strip of film**. The array is the
film; the slice is a rectangle of viewfinder showing you some section of
it — and possibly several other viewfinders are pointed at overlapping
sections of the same strip, simultaneously.

<!-- pause -->

Scratch the film through your viewfinder, and anyone with an overlapping
view sees the scratch. There's only one strip of film underneath.

<!-- end_slide -->

## Then append gets involved

Splicing more frames onto the end of your section of film (`append`ing
past capacity) goes one of two ways:

<!-- incremental_lists: true -->

- Room left on the existing strip: your new frames land right there, possibly overwriting someone else's neighboring section
- No room left: Go quietly photocopies the whole strip onto fresh film, and your viewfinder now points at the copy instead

<!-- incremental_lists: false -->

**Both look identical from the code's point of view.** Only checking the
data — or the capacity — tells you which one happened.

<!-- end_slide -->

## Watch it happen: cap() before and after

```go
original := []int{1, 2, 3, 4, 5}
view := original[1:3]
fmt.Println(len(view), cap(view)) // 2 4

view = append(view, 100) // still room — overwrites original[3]
fmt.Println(original)     // [1 99 3 100 5]

view = append(view, 200, 300, 400) // exceeds capacity now
fmt.Println(len(view), cap(view))  // cap jumped — reallocated

view[0] = -1
fmt.Println(original) // unchanged — sharing broke
```

**Demo:** run `code/cmd/slicealias` live and read the `len`/`cap` at every
step out loud before revealing whether `original` changed.

<!--
speaker_note: |
  This is the slide to slow down on. Walk line by line, predict the
  output with the room before running it, then run it. The moment cap()
  jumps to a value bigger than the previous cap is the exact moment
  aliasing stops - point at that number specifically.
-->

<!-- end_slide -->

## Range: one keyword, five collection types

`range` iterates arrays, slices, strings (by rune, not byte), maps, and
channels — always yielding index/key plus value.

```go
for i, v := range []string{"a", "b", "c"} { /* i, v */ }
for i, r := range "héllo"                  { /* i, r rune */ }
for k, v := range map[string]int{"x": 1}   { /* k, v */ }
```

<!-- pause -->

If you only need the value, discard the index with `_`: `for _, v := range
s`. Discarding is explicit — Go won't let an unused variable slide.

<!-- end_slide -->

## Maps: comma-ok, not exceptions

```go
m := map[string]int{"a": 1}
v, ok := m["b"]
// v == 0 (int's zero value), ok == false
```

<!-- pause -->

No `KeyError`, no `NullPointerException`. A missing key returns the zero
value silently — `ok` is how you tell "genuinely zero" apart from "never
there."

<!-- end_slide -->

## Maps have no guaranteed order — on purpose

Unlike Python 3.7+ dicts (insertion-ordered) or Java's `LinkedHashMap`,
ranging over a Go map gives a **different order on different runs**,
deliberately randomized.

<!-- pause -->

Think of it as a **coat-check counter, not a queue**: you hand over your
coat and get a ticket back; the attendant hands coats back whenever they
feel like reaching for them, not in arrival order. If you need an order,
you track it yourself — a sign-in sheet, not just hooks on a wall.

**Demo:** run the same range-over-map loop twice in one program, then run the whole program again — compare all three orders.

<!-- end_slide -->

## nil slice vs. nil map: not the same nil

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
var s []int
s = append(s, 1)
// fine — allocates
// on first use
```

<!-- column: 1 -->

```go
var m map[string]int
m["x"] = 1
// panic: assignment
// to entry in nil map
```

<!-- reset_layout -->

<!-- pause -->

Both are safe to **read** (`len(s)` is 0, `m["x"]` returns the zero value).
Only the map panics on write. Initialize with `make(map[K]V)` or a literal
before you write to one.

<!--
speaker_note: |
  Great "gotcha of the day" - short, memorable, and something people
  genuinely hit in week one of writing real Go, usually as a nil-pointer-
  style panic that's confusing the first time because the surrounding
  code looks fine.
-->

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **One loop keyword**: `for` covers classic, while-style, infinite, and range-based looping
2. **`switch` doesn't fall through by default**: opposite of C/Java/JS, opt in with `fallthrough`
3. **Structs are values**: passed by copy unless you pass a pointer explicitly
4. **Arrays are fixed-size value types; slices are the flexible view you actually use**
5. **Slicing shares the underlying array**: mutating through one slice mutates through all of them, until `append` forces a reallocation
6. **Maps have no guaranteed order, and nil maps panic on write** — nil slices don't
7. **A test is still just a `_test.go` file and `TestX(t *testing.T)`** — apply that to this lab's `switch`/`fallthrough` and pointer-mutation code

<!-- end_slide -->

## Back to the opening scenario

```go
backup := original
backup[0] = 999
```

`original` changed because `backup` was never a copy of the data — it was
a copy of a small header (pointer, length, capacity) pointing at the exact
same underlying array as `original`.

<!-- pause -->

**Type in chat: what would you write instead, if you actually wanted an independent copy?**

(`copy(dst, src)` into a freshly made slice, or `append([]T{}, original...)`
— either forces new backing memory.)

<!--
speaker_note: |
  Let a few answers land before confirming copy() - some people will
  guess append-to-empty-slice, which also works and is worth validating.
  The point is that "backup := original" was never doing what it looked
  like it was doing.
-->

<!-- end_slide -->

## Bridge to Topic 4

**We've established:**

<!-- incremental_lists: true -->

- `for`, `if`, and `switch` cover all of Go's control flow, with a couple of reversed defaults to unlearn
- Structs are plain value types; pointers are how you opt into mutation
- Slices alias their underlying array — the source of the most common real-world Go bug
- Maps trade guaranteed order for safety, and nil maps and nil slices behave asymmetrically
- The lab's last step is a test again, same as Topic 2's

<!-- incremental_lists: false -->

**Topic 4: Object-Oriented Programming** — methods on structs, receivers,
and how Go builds behavior around data without a `class` keyword.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
