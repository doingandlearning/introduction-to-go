---
title: "**Inheritance and Interfaces**"
sub_title: Go Programming — Topic 5
author: Kevin Cunningham
---

## Opening scenario

```go
type MyError struct{}

func (e *MyError) Error() string { return "something broke" }

func doWork() error {
    var e *MyError = nil
    // ... some logic that never sets e ...
    return e
}

func main() {
    err := doWork()
    if err == nil {
        fmt.Println("all good")
    } else {
        fmt.Println("failed:", err)
    }
}
```

**Type in chat: `e` is `nil`. What does this print?**

<!--
speaker_note: |
  Take guesses in chat - most people confidently say "all good" because
  e is visibly nil. Don't resolve it yet. We come back to this once
  we've covered how an interface value is actually stored, later in
  the deck.
-->

<!-- end_slide -->

## No `extends`, no class tree

Go has no class keyword, no inheritance keyword, and no hierarchy to
walk. Two mechanisms do the job instead:

<!-- incremental_lists: true -->

- **Struct embedding** — one type placed inside another, for code reuse
- **Interfaces** — method-signature sets, for polymorphism

<!-- incremental_lists: false -->

**Anchor fact, not a vibe:** embedding is composition, not inheritance.
There is no `super`, and there is no virtual dispatch. We'll come back to
exactly what that means in a moment — hold onto it.

<!-- end_slide -->

## Struct embedding

Put a type inside another struct with **no field name**, and its fields
and methods are *promoted* onto the outer type.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
type Person struct {
    Name string
}

func (p Person) Greet() string {
    return "Hi, I'm " + p.Name
}

type Employee struct {
    Person // embedded, no name
    Role   string
}
```

<!-- column: 1 -->

```go
e := Employee{
    Person: Person{Name: "Sam"},
    Role:   "Engineer",
}

e.Name    // "Sam"
          // promoted field

e.Greet() // "Hi, I'm Sam"
          // promoted method
```

<!-- reset_layout -->

**Demo:** build this live, then delete the `Role` field and show that nothing about the promotion changes — it has nothing to do with what else is on `Employee`.

<!-- end_slide -->

## Composition, not inheritance

`Employee` **contains** a `Person`. It did not extend one.

<!-- pause -->

That distinction has two concrete, checkable consequences:

<!-- incremental_lists: true -->

- **No virtual dispatch.** If `Employee` defines its own `Greet()`, it *shadows* the promoted one — it doesn't override it polymorphically, because there's no shared base-class pointer for a virtual call to resolve through
- **No `super.Greet()`.** If shadowed code wants the embedded version too, it has to say so explicitly

<!-- incremental_lists: false -->

<!--
speaker_note: |
  These two bullets are the anchor from the opening slide - point back
  to them explicitly. This is the moment to say plainly: this is not
  inheritance wearing a disguise, it behaves differently in ways that
  matter the first time you shadow a method.
-->

<!-- end_slide -->

## The manual "super call"

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
func (e Employee) Greet() string {
    return e.Person.Greet() +
        ", I work here"
}
```

<!-- column: 1 -->

`e.Person.Greet()` is doing the job `super.Greet()` would do in Java —
but it isn't special syntax. It's a normal field access (`e.Person`)
followed by a normal method call.

<!-- reset_layout -->

**Demo:** call `e.Greet()` before adding the override (promoted version runs), then after (shadowed version runs), then show `e.Person.Greet()` still reaches the original directly.

<!-- end_slide -->

## The kitchen island

<!-- pause -->

A kitchen island with a built-in dishwasher, not a parent-child
biological relationship.

<!-- incremental_lists: true -->

- The island didn't *inherit* "wash dishes" — it has a dishwasher bolted inside it
- Its front panel happens to expose the same buttons, so from outside it looks like the island can wash dishes
- Want a pre-rinse step? You build it into the island's own button, which calls the dishwasher's original cycle **explicitly**
- There's no automatic "parent dishwasher" your button calls into for free

<!-- incremental_lists: false -->

<!-- end_slide -->

<!-- jump_to_middle -->

Interfaces: behavior, not pedigree
===

<!-- end_slide -->

## Interfaces are method sets

An interface names a set of method signatures. Any type with a matching
method set satisfies it — automatically, with no declared link.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
type Greeter interface {
    Greet() string
}
```

<!-- column: 1 -->

```go
func WelcomeAll(gs []Greeter) {
    for _, g := range gs {
        fmt.Println(g.Greet())
    }
}
```

<!-- reset_layout -->

`Person` and `Employee` from earlier both already satisfy `Greeter` —
just by having a matching `Greet() string` method. Nothing else needed.

<!-- end_slide -->

## No `implements` keyword, anywhere

This is the single biggest interface-related difference from Java.

<!-- pause -->

```go
type Robot struct{}

func (r Robot) Greet() string { return "BEEP BOOP HELLO" }

// Robot satisfies Greeter with zero declaration of intent.
// WelcomeAll([]Greeter{Person{Name: "Sam"}, Robot{}}) just works.
```

`Robot`'s author never saw `Greeter`. Doesn't matter. The method set
matches, so it qualifies.

**Type in chat: TypeScript folks — does this feel familiar? Java folks — does this feel unsettling?**

<!--
speaker_note: |
  TS devs will often say "yeah, that's just structural typing" almost
  immediately - their interfaces work the same way. Java devs tend to
  need a beat longer; let "wait, I don't have to declare anything?"
  actually land instead of rushing past it. Budget real time here.
-->

<!-- end_slide -->

## The loading dock

A loading dock doesn't check whether a package is *certified* to be
lifted, scanned, and signed for — it just tries to lift it, scan it, and
get a signature.

<!-- incremental_lists: true -->

- A department-issued crate gets through
- Someone's personal Amazon return gets through
- A box from a company that's never heard of this dock gets through

<!-- incremental_lists: false -->

**The dock cares about behavior, not pedigree.** That's exactly how a Go
interface treats any type with the right methods.

<!-- end_slide -->

## Small interfaces are idiomatic

*"The bigger the interface, the weaker the abstraction."*

<!-- pause -->

The standard library takes this seriously — not just as a proverb:

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

<!-- column: 1 -->

```go
type Stringer interface {
    String() string
}
```

`io.Reader`, `io.Writer`, `fmt.Stringer` — one method each. Compare to
Java's tendency toward larger, more comprehensive interfaces.

<!-- reset_layout -->

**Demo:** show how many unrelated stdlib types satisfy `io.Reader` — files, network connections, byte buffers, string readers — none of them related to each other by hierarchy.

<!-- end_slide -->

## Polymorphism without a hierarchy

No "is-a" tree to walk. Unrelated types — a `Person`, a `Robot`,
something from a completely different package — are interchangeable
anywhere a `Greeter` is expected, as long as each independently has the
right method.

<!-- pause -->

**This is structurally identical to duck typing** — Python and JS
developers already think this way. Go just checks it at compile time
instead of at call time.

<!-- end_slide -->

## When even a small interface is too narrow: `any`

Sometimes you genuinely don't know the type ahead of time. `any` (the
modern alias for `interface{}`) accepts literally anything.

```go
func describe(x any) {
    switch v := x.(type) {
    case int:
        fmt.Println("int:", v*2)
    case string:
        fmt.Println("string:", strings.ToUpper(v))
    default:
        fmt.Printf("something else: %v\n", v)
    }
}
```

**A type switch, not a type check.** `v` is re-typed inside each `case`
— inside `case int`, `v` is an `int`, not an `any` you have to cast.

<!-- end_slide -->

<!-- jump_to_middle -->

Back to the opening scenario
===

<!-- end_slide -->

## An interface value is a pair

`err == nil` printed **`false`** — even though `e` was visibly `nil`.

<!-- pause -->

An interface value is really a pair: **(type, value)**.

```go
var e *MyError = nil
var err error = e
```

<!-- pause -->

The pair stored in `err` is `(*MyError, nil)` — the *type* part is
non-nil even though the *value* part is. `err == nil` checks the whole
pair, and the pair isn't zero, because a concrete type got recorded.

<!--
speaker_note: |
  Draw the (type, value) box on the whiteboard if you have one - two
  cells, one labeled type one labeled value. Fill in *MyError and nil
  respectively, then point at the type cell: "that's why it's not nil."
-->

<!-- end_slide -->

## Nobody is immune to this one

This trips up **senior engineers**, not just beginners.

<!-- pause -->

**The fix:** don't let a concrete pointer type leak into an interface
return slot. Return a literal `nil` in the success path instead of a
variable declared as a concrete pointer type.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// Leaky pattern:
func doWork() error {
    var e *MyError = nil
    return e // (*MyError, nil)
}
```

<!-- column: 1 -->

```go
// Fixed:
func doWork() error {
    if somethingWentWrong {
        return &MyError{}
    }
    return nil // the real nil
}
```

<!-- reset_layout -->

**Demo:** run both versions and print `err == nil` for each — same-looking success path, different result.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Embedding promotes fields and methods**: no field name on the embedded type is what triggers it
2. **Embedding is composition, not inheritance**: no `super`, no virtual dispatch — shadowing replaces, it doesn't override
3. **Interfaces are satisfied implicitly**: no `implements` keyword exists in Go at all
4. **Small interfaces are idiomatic**: `io.Reader`, `io.Writer`, `fmt.Stringer` are one method each
5. **A nil pointer stored in an interface is not itself nil**: an interface value is a (type, value) pair, and only the value half was nil
6. **This is testable, not just explainable**: a `TestX(t *testing.T)` that checks `err == nil` on a success path catches the nil-interface gotcha the moment anyone reintroduces it

<!-- end_slide -->

## Bridge to Topic 6

**We've established:**

<!-- incremental_lists: true -->

- Struct embedding gives you code reuse through composition, with promotion sugar on top
- Interfaces give you polymorphism through method sets, satisfied structurally with zero declared intent
- The (type, value) pair behind every interface value explains the nil-interface gotcha, and it's worth respecting even as an expert
- The lab's closing exercise puts a `TestX` around exactly that gotcha, plus interface satisfaction across `Person`, `Kiosk`, and `Mascot` — the same `go test` habit from Topic 2, now aimed at this topic's content

<!-- incremental_lists: false -->

**Topic 6: Functional Programming** — functions as values, closures, and
how Go leans on small interfaces and function types instead of deep
class hierarchies to keep code composable.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
