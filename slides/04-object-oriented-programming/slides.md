---
title: "**Object-Oriented Programming**"
sub_title: Go Programming — Topic 4
author: Kevin Cunningham
---

## Opening scenario

Here's a `Vehicle` in Java: three fields, a constructor, one method — all
sealed inside `class Vehicle { ... }`.

```java
public class Vehicle {
    private String plate;
    private int mileage;

    public Vehicle(String plate, int mileage) {
        this.plate = plate;
        this.mileage = mileage;
    }

    public double estimatedServiceCost() {
        return mileage * 0.05;
    }
}
```

**Type in chat: Go has no `class` keyword at all. What do you think replaces this?**

We'll come back to this once you've seen the actual Go version.

<!--
speaker_note: |
  Let a few guesses land - "struct," "interface," and "I have no idea"
  are all common. Don't confirm or deny yet, just bank the guesses and
  move into the struct-plus-methods reveal.
-->

<!-- end_slide -->

## Data and behavior, declared separately

Go splits what Java bundles into one unit. A **struct** holds the data.
**Methods** are functions with a receiver, declared separately, and Go
links them together.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
type Vehicle struct {
    Plate   string
    Mileage int
}
```

<!-- pause -->  

<!-- column: 1 -->

```go
func (v Vehicle) EstServiceCost() float64 {
    return float64(v.Mileage) * 0.05
}
```

<!-- reset_layout -->

<!-- pause -->

**Two declarations, not one.** The struct doesn't know the method exists
until you read the whole file. There's no `this.field` scope holding it
all together — just a receiver name you choose yourself.

<!--
speaker_note: |
  This is the payoff for the opening slide - the mental shift is "data
  and behavior are two separate declarations Go links via the
  receiver," not "one indivisible class unit." Don't mention
  inheritance or embedding here - that's Topic 5, and bringing it up
  now makes people reach for a Go feature that doesn't exist yet in
  this course.
-->

<!-- end_slide -->

## A method is just a function with a receiver

`func (v Vehicle) EstServiceCost() float64` is barely different from a
free function `EstServiceCost(v Vehicle) float64`.

The receiver — `(v Vehicle)` — is sugar for writing `v.EstServiceCost()`
instead of `EstServiceCost(v)`. Nothing more magical happens.

<!-- pause -->

**This is why Go has no method overloading.** There's no class-scoped
namespace for method names to live in — `EstServiceCost` is just a
package-level identifier with a receiver attached, and Go doesn't allow
two functions with the same name in one package, full stop.

<!-- pause -->

**You've already seen this shape once.** Topic 2's `Stringer` example,
`func (p Point) String() string`, is exactly this pattern — a receiver
attached to a function, just named `String` because that's specifically
what `fmt` goes looking for.

<!-- end_slide -->

## Value receiver: works on a copy

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
type Item struct {
    Name     string
    Quantity int
    UnitCost float64
}

func (i Item) TotalValue() float64 {
    return float64(i.Quantity) *
        i.UnitCost
}
```

<!-- column: 1 -->

```go
it := Item{
    Name: "Widget", Quantity: 10,
    UnitCost: 5,
}
it.TotalValue() // 50

// TotalValue can't have
// mutated it - i is a copy.
```

<!-- reset_layout -->

**Photocopy-and-total desk:** hands you a number, the original form on
the counter is untouched.

<!-- end_slide -->

## Pointer receiver: works on the original

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
func (i *Item) ApplyDiscount(pct float64) {
    i.UnitCost *= (1 - pct/100)
}
```

<!-- column: 1 -->

```go
it.ApplyDiscount(10)
it.TotalValue() // 45

// Go auto-took the address of
// it here: it.ApplyDiscount(10)
// means (&it).ApplyDiscount(10)
```

<!-- reset_layout -->

<!-- pause -->

**"Approved" stamp:** stamps and hands back the *same* form, changed in
content. Different desk, same paper.

**Demo:** run both methods against one `Item` and print `UnitCost`
between each call — watch which one changes it and which one doesn't.

<!--
speaker_note: |
  Use code/cmd/receiverdemo for this. Print UnitCost before
  TotalValue, after TotalValue, and after ApplyDiscount - the first two
  should be identical, only the third changes. That's the whole point
  landing visually rather than just verbally.
-->

<!-- end_slide -->

## Pick one convention per type

Value vs pointer receiver is a real per-type decision — Go won't choose
for you, and it won't stop you from mixing them.

<!-- incremental_lists: true -->

- If **any** method on a type needs to mutate it, give **all** its
  methods pointer receivers
- Same if the struct is large — copying it on every value-receiver call
  gets expensive
- Mixing the two on one type is a real source of confusion: callers
  can't always tell which behavior they're getting just by reading the
  call site

<!-- incremental_lists: false -->

**Type in chat: if you saw `it.ApplyDiscount(10)` in review with no other
context, could you tell whether `it` just got mutated?**

<!--
speaker_note: |
  The honest answer is "not reliably" without checking the method
  declaration - that's exactly the problem consistency solves. Let
  that sit for a beat before moving on.
-->

<!-- end_slide -->

## No enforced constructor

Go's zero value **is** a valid instance. `var it Item` is immediately
usable — no null, no uninitialized-object exception.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
var it Item
// Name: "", Quantity: 0,
// UnitCost: 0
it.TotalValue() // 0 — runs fine
```

<!-- column: 1 -->

```java
Vehicle v; // declared, not built
v.estimatedServiceCost();
// NullPointerException
```

<!-- reset_layout -->

**In Java, an unconstructed object doesn't exist.** In Go, it exists —
just with every field at its type's zero value — and any method can run
against it immediately.

<!-- end_slide -->

## `NewX`: convention, not syntax

Go has no constructor keyword. The convention is a plain function named
`NewX` that builds and validates, often returning a pointer plus an
error.

```go
func NewItem(name string, qty int, cost float64) (*Item, error) {
    if qty < 0 {
        return nil, fmt.Errorf("quantity cannot be negative, got %d", qty)
    }
    if name == "" {
        name = "Unknown"
    }
    return &Item{Name: name, Quantity: qty, UnitCost: cost}, nil
}
```

<!-- pause -->

**Nothing stops you from skipping `NewItem` entirely** and writing
`Item{...}` directly, bypassing every check inside it. The convention
only holds if every caller chooses to use it — the compiler won't force
anyone through this function.

<!-- end_slide -->

<!-- jump_to_middle -->

The map gotcha
===

<!-- end_slide -->

## Map values aren't addressable

```go
catalog := map[string]Item{
    "widget": {Name: "Widget", Quantity: 10, UnitCost: 5},
}

catalog["widget"].ApplyDiscount(10)
// cannot call pointer method on catalog["widget"]
// cannot take address of catalog["widget"]
```

<!-- pause -->

Go usually auto-takes the address for you — that's what made
`it.ApplyDiscount(10)` work earlier on a local variable. A map entry has
no stable address: the map can rehash and move it at any time, so Go
refuses to hand you a pointer into it at all.

**Demo:** uncomment the line in `cmd/mapgotcha/main.go` and run `go
build` — read the actual compiler error before revealing the fix.

<!--
speaker_note: |
  Trigger this error live, don't just show it in slides. Delegates
  remember the error message far better after watching it happen than
  after reading it on a slide. Let them read it themselves before you
  explain the "no stable address" reasoning.
-->

<!-- end_slide -->

## Fix: read, mutate, write back

```go
widget := catalog["widget"]   // read: copies the value out
widget.ApplyDiscount(10)      // mutate: the local copy
catalog["widget"] = widget    // write back: replaces the map entry
```

<!-- pause -->

Three steps, no shortcut. This is also the moment mixed receivers on one
type get genuinely confusing — if `ApplyDiscount` had been a value
receiver instead, this whole dance would compile fine and silently do
nothing, because you'd be discounting a copy that's thrown away.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **No `class` keyword**: a struct (data) plus methods (functions with a
   receiver) plays that role, declared separately
2. **A method is a function with a receiver**: sugar for
   `Func(receiver, ...)`, which is also why Go has no overloading
3. **Value receiver copies, pointer receiver mutates the original**: pick
   one convention per type and stick to it
4. **The zero value is a valid instance**: no null, no
   uninitialized-object crash
5. **`NewX` is convention, not syntax**: validates and builds, but
   nothing forces callers to use it
6. **Map values aren't addressable**: pointer-receiver calls on them
   need read-mutate-write-back
7. **A test still catches what a value-receiver method silently doesn't
   mutate**: apply Topic 2's `_test.go` habit here too

<!-- end_slide -->

## Back to the opening scenario

The Java `Vehicle` sealed data, constructor, and method into one
`class` block — reachable only through that object, from construction
onward.

**The Go version is two separate declarations**: a `Vehicle` struct with
no behavior of its own, and a function that happens to declare `v
Vehicle` as its first argument, written as a receiver.

<!-- pause -->

**Type in chat: which fields would you make pointer receivers on now that
mileage needs to update after a service visit?**

<!--
speaker_note: |
  The expected answer is "all of them, once ApplyService needs to
  mutate Mileage" - reinforce the "pick one convention per type" rule
  from earlier rather than letting the room debate field-by-field
  receiver choices, which isn't how Go receivers work (they're per
  method, chosen for the whole type by convention, not per field).
-->

<!-- end_slide -->

## Bridge to Topic 5

**We've established:**

<!-- incremental_lists: true -->

- Structs hold data, methods (functions with a receiver) provide
  behavior — two separate declarations, not one class
- Value receivers copy, pointer receivers mutate — pick one per type
- `NewX` is a validating-constructor convention, not a language rule
- Map values aren't addressable, so pointer methods need
  read-mutate-write-back
- The lab's last step is a test again, same as every topic since Topic 2

<!-- incremental_lists: false -->

**Topic 5: Inheritance and Interfaces** — how Go shares behavior across
types without a `class` hierarchy at all, and what actually replaces
`extends`.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
