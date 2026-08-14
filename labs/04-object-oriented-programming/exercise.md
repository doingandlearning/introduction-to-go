# Lab 4: Object-Oriented Programming — a small library catalog

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

Both directories contain the same shape: a `cmd/library` package (the
program entry point) and an `internal/catalog` package holding a `Book`
struct and its methods. By the end, you'll have written both kinds of
receiver, built a validating constructor, triggered a real compiler
error on purpose, and fixed a subtle mixed-receiver bug.

---

## Exercise 1: A value-receiver method

**Objective:** Compute a derived value from a struct without mutating
it.

**Context:** `internal/catalog/catalog.go` has a `Book` struct and a
`TODO` where `EstimatedReadHours` should live.

`main.go` already builds a `Book` with a struct literal and prints its
`EstimatedReadHours()` — you don't need to write that part.

**Tasks:**

1. Implement `func (b Book) EstimatedReadHours() float64` so it returns
   `float64(b.PageCount) / 40` (a rough "40 pages per hour" estimate).
2. Run `go run ./cmd/library` and confirm the number under "Exercise 1"
   looks right for the page count in `main.go` (`380` pages ≈ `9.5`
   hours).

**Key Learning:** A value-receiver method operates on a copy of the
struct. Nothing about calling it can change the original — it's a
read-only computation with method syntax.

---

## Exercise 2: A pointer-receiver method, and the map gotcha

**Objective:** Mutate a struct through a method, then hit — and fix —
the map-of-structs gotcha.

**Context:** `Checkout` should decrement `CopiesAvailable` and increment
`TimesBorrowed`. It needs a pointer receiver, because it mutates the
`Book`.

`main.go` already calls `Checkout` on a plain `Book` variable and builds
a `map[string]catalog.Book` called `library` — you don't need to write
that part.

**Tasks:**

1. Implement `func (b *Book) Checkout() error`. It should return an
   error if `CopiesAvailable` is already `0`, otherwise decrement
   `CopiesAvailable`, increment `TimesBorrowed`, and return `nil`.
2. Run `go run ./cmd/library` and confirm `CopiesAvailable` changes
   between the "before" and "after" prints on the plain `Book`
   variable — that's the mutation sticking.
3. In `main.go`, find the commented-out line
   `// library["go-in-action"].Checkout()`. Uncomment it and run `go
   build ./cmd/library`. It won't compile — read the error carefully
   before doing anything else. Then re-comment the line.
4. Just below it, write the read-mutate-write-back fix yourself: read
   the `Book` out of the map into a local variable, call `Checkout` on
   the local variable, then write it back into the map under the same
   key. Print `library["go-in-action"].CopiesAvailable` afterward to
   confirm it changed.

**Key Learning:** Go can auto-take the address of a local, addressable
variable (that's what made step 2 work), but a map value has no stable
address — the map can move it at any time. Go refuses to compile a
pointer-method call on a map entry rather than hand you a pointer that
might go stale.

---

## Exercise 3: A validating constructor

**Objective:** Build the `NewX` convention: a plain function that
validates and constructs.

`main.go` already calls `NewBook` once with valid input and once with a
negative copy count, handling both return values — you don't need to
write that part.

**Tasks:**

1. Implement `func NewBook(title, author string, copies int) (*Book,
   error)`. It should reject a negative `copies` with an error made via
   `fmt.Errorf`, including the rejected value in the message. Otherwise
   it returns a `*Book` with the given fields and `TimesBorrowed`
   starting at `0`.
2. Run `go run ./cmd/library` and confirm the "Exercise 3" output now
   shows a real created book on the valid call and a real rejection
   message on the invalid one, instead of the placeholder zero-valued
   book both calls currently produce.

**Key Learning:** `NewX` is a convention for "validate, then build," not
a language feature. Nothing stops any caller from writing
`catalog.Book{...}` directly and skipping validation entirely — the
compiler will not force anyone through `NewBook`.

---

## Exercise 4: The zero value is a real value

**Objective:** Confirm a Go struct is usable with no constructor call at
all.

`main.go` already declares `var zero catalog.Book` with no struct
literal and no `NewBook` call, and immediately calls
`zero.EstimatedReadHours()` on it.

**Tasks:**

1. Run `go run ./cmd/library` (after finishing Exercise 1) and look at
   the "Exercise 4" output.
2. Note what you get, and why: every field is at its type's zero value
   (`PageCount` is `0`), so the calculation runs against a fully valid,
   if uninteresting, `Book` and returns a sane `0` — nothing crashes.

**Key Learning:** Coming from Java, an uninitialized object reference is
`null`, and calling a method on it throws `NullPointerException`. In Go,
`var b catalog.Book` is already a complete, zero-valued instance — there
is no "not yet constructed" state a method call can crash into.

---

## Exercise 5: Mixed receivers, on purpose

**Objective:** See why mixing value and pointer receivers on one type is
a real bug, not just a style nit.

**Context:** `catalog.go` has a `Return` method already written —
*deliberately* with a value receiver, even though `Checkout` on the same
type uses a pointer receiver.

`main.go` already calls `b5.Return()` on a `Book` with a known
`CopiesAvailable` and prints the value before and after.

**Tasks:**

1. Run `go run ./cmd/library` and look at the "Exercise 5" output.
   `CopiesAvailable` won't have changed — figure out why before reading
   further.
2. Fix `Return` in `catalog.go` by changing its receiver from `(b Book)`
   to `(b *Book)`, matching `Checkout`.
3. Re-run `go run ./cmd/library`. Confirm `CopiesAvailable` now
   increments correctly.

**Key Learning:** A value receiver silently compiles and silently does
nothing useful when the method's whole job is to mutate — there's no
error, just a change that vanishes when the copy goes out of scope. The
fix is consistency: once any method on a type needs to mutate it, every
method on that type should use a pointer receiver, so callers can trust
one rule instead of checking each method's signature.

---

## Exercise 6: Prove it with a test

**Objective:** Apply the `_test.go` habit from Topic 2 to the receiver
behavior you just built — a test is where "did this actually mutate?"
stops being a guess.

**Context:** You already know the mechanics from Topic 2: a `_test.go`
file, a `TestX(t *testing.T)` function, `t.Errorf`/`t.Fatalf`, no
framework needed. Nothing new there — this exercise just points that
habit at `NewBook`, `EstimatedReadHours`, and `Checkout`.

**Tasks:**

1. Create `internal/catalog/catalog_test.go` (package `catalog`, only
   import `testing`).
2. Write `TestNewBook_Success`: call `NewBook` with valid input, check
   the returned fields match what you passed in and `TimesBorrowed` is
   `0`.
3. Write `TestNewBook_RejectsNegativeCopies`: call `NewBook` with a
   negative `copies`, confirm you get a non-nil error and a nil
   `*Book`.
4. Write `TestEstimatedReadHours_ValueReceiverDoesNotMutate`: build a
   `Book`, call `EstimatedReadHours()`, check the returned value, then
   check the `Book`'s fields are still what you set them to — make the
   "a value receiver can't mutate" guarantee visible in test form,
   instead of just taking it on faith.
5. Write `TestCheckout_PointerReceiverMutates`: build a `Book` with some
   `CopiesAvailable`, call `Checkout()`, confirm `CopiesAvailable`
   decremented and `TimesBorrowed` incremented — the pointer-receiver
   counterpart to step 4.
6. Run `go test ./...` and confirm every test passes.
7. Now deliberately break it: temporarily change `Checkout`'s receiver
   from `(b *Book)` to `(b Book)`. Run `go test ./...` again and read
   the failure — `TestCheckout_PointerReceiverMutates` should fail,
   because the mutation no longer sticks to anything outside the copy.
   Change the receiver back to `(b *Book)` and confirm the test passes
   again.

**Key Learning:** "Did this method actually mutate what I expected?" is
exactly the kind of thing that's easy to eyeball wrong — especially with
a value/pointer receiver mix-up like Exercise 5's, where the code
compiles fine either way and only the *behavior* is wrong. A test that
asserts on the field after the call catches that silently-vanishing
mutation for free, every single run — not just the one time you
happened to print and check it by eye.

---

## Summary

By the end of this lab you should be able to:

- Write a value-receiver method that computes a derived result without
  mutating its receiver
- Write a pointer-receiver method that mutates its receiver, and explain
  why the receiver kind has to match the intent
- Read and fix the "cannot take address of a map value" compiler error
  using read-mutate-write-back
- Write a validating `NewX` constructor following Go's convention, and
  explain why it isn't enforced by the compiler
- Explain why a Go zero-value struct is usable immediately, unlike an
  uninitialized object reference in Java
- Recognize a value receiver that should be a pointer receiver from its
  silently-vanishing mutation, and fix it
- Write a `_test.go` file that pins down the value-vs-pointer receiver
  distinction — asserting a value receiver leaves its original untouched
  and a pointer receiver mutates it — reusing the testing habit from
  Topic 2
