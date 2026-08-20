# Lab 5: Inheritance and Interfaces — a library front-desk system

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

**Every exercise below ships with its test already written**, sitting in
`starter/internal/library/library_test.go`, `greeter_test.go`,
`checkout_test.go`, and `starter/cmd/frontdesk/main_test.go`. Run `go
test ./...` from `starter/` right now, before you change anything. Your
job in each exercise is to make the named test pass, not to write a new
one; writing tests yourself is Topic 12's job.

**One heads-up specific to this lab:** Exercises 1, 2, 3, and 5 all live
in the same `internal/library` package, and Go compiles a package's
test files together as a single unit. Until Exercises 1, 2, and 3 are
*all* done, that package won't build for testing at all — running `go
test ./...` right now shows a batch of compiler errors from more than
one file at once, not one clean failure. That's Go's package-is-the-
compilation-unit rule doing exactly what it's supposed to, not a bug in
this test setup — each exercise below tells you which errors are yours
to look at right now and which belong to a later exercise. Exercise 4
lives in its own package (`cmd/frontdesk`) and is unaffected by any of
this — its test fails cleanly and independently from the start.

Both directories contain the same shape: a `cmd/frontdesk` package (the
program entry point) and an `internal/library` package (the types you'll
build out). The scenario: Sunnyvale Public Library is rolling out a
front-desk greeting system that needs to work for human staff,
volunteers, self-service kiosks, and a costumed reading-time mascot —
without a class hierarchy in sight.

---

## Exercise 1: Embed Person in Volunteer

**Objective:** Confirm struct embedding promotes a field and a method
with zero extra code.

**Context:** `starter/internal/library/library.go` has a `Person` type
with a `Name` field and a `Greet()` method already implemented. `Volunteer`
is stubbed out with a `TODO` where the embedding should go.
`TestVolunteerPromotion` in `library_test.go` is already written and
covers both this exercise and the next — right now it won't even
compile, since it references `v.Person`, `v.Name`, and `v.Greet()`, none
of which exist on `Volunteer` yet.

**Tasks:**

1. Run `go test ./...` from `starter/`. Read the compiler errors coming
   from `library_test.go` — `unknown field Person`, `v.Name undefined`,
   and so on. (You'll also see one from `greeter_test.go` — ignore that
   one for now, it's Exercise 3's, not yours.)
2. Open `starter/internal/library/library.go`. Embed `Person` in
   `Volunteer` — no field name, just the type name on its own line —
   alongside the existing `ShiftHours int` field.
3. In `starter/cmd/frontdesk/main.go`, uncomment the `Volunteer` block in
   `main()`.
4. Run `go run ./cmd/frontdesk`. Confirm `v.Name` prints `"Sam"` and
   `v.Greet()` produces a full greeting — neither one has any code on
   `Volunteer` itself yet.
5. Re-run `go test ./...`. `library_test.go`'s compile errors are gone,
   but `TestVolunteerPromotion` now fails one assertion instead — the
   shadowing check. That's expected; Exercise 2 finishes it.
   (`greeter_test.go` is still blocking the package build — still
   Exercise 3's.)

**Key Learning:** Embedding a type with no field name promotes its
exported fields and methods onto the outer type automatically. This is
composition with syntactic sugar, not inheritance — `Volunteer` *contains*
a `Person`.

---

## Exercise 2: Shadow Greet(), then call through explicitly

**Objective:** See the difference between shadowing and polymorphic
override, and write Go's equivalent of a `super` call.

**Context:** `TestVolunteerPromotion` compiles now but is still failing
after Exercise 1, on its shadowing assertion.

**Tasks:**

1. In `starter/internal/library/library.go`, give `Volunteer` its own
   `Greet()` method. It should call `v.Person.Greet()` explicitly and
   append `" I'm volunteering today."` to the result.
2. Run `go run ./cmd/frontdesk` again. Confirm `v.Greet()` now includes
   the volunteering line, while `v.Person.Greet()` — called directly —
   still returns the original, un-appended greeting.
3. Run `go test ./...`. `TestVolunteerPromotion` now passes in full.
   (`internal/library` as a whole still won't build yet — Exercise 3 is
   next.)

**Key Learning:** Shadowing replaces a promoted method for callers of the
outer type; it doesn't override it polymorphically, because there's no
shared base-class pointer for a virtual call to resolve through. There's
no `super.Greet()` keyword — `v.Person.Greet()` is just a normal field
access plus a normal method call, and it's the only way to reach the
embedded version once it's shadowed.

---

## Exercise 3: Three unrelated types, one interface, no "implements"

**Objective:** Prove interfaces are satisfied structurally, with no
declared relationship required.

**Context:** `starter/internal/library/greeter.go` has a `Greeter`
interface stub and two unimplemented types, `Kiosk` and `Mascot`.
`Person` (from Exercise 1) already has a `Greet() string` method.
`TestGreeterSatisfiedByMultipleTypes` in `greeter_test.go` is already
written — it's the thing that's been blocking `go test ./...` from
building at all since Exercise 1, because `Greeter` currently has no
methods for it to call.

**Tasks:**

1. Define the `Greeter` interface with one method: `Greet() string`.
2. Implement `Greet()` on `Kiosk` (e.g. return
   `"SCAN YOUR CARD AT STATION <N> TO BEGIN."`).
3. Implement `Greet()` on `Mascot` (e.g. return something like
   `"*<Name> waves enthusiastically, says nothing*"`).
4. `WelcomeAll`, at the bottom of `greeter.go`, is commented out for the
   same reason the test wouldn't build — uncomment it now, and add
   `"fmt"` back to this file's imports.
5. In `starter/cmd/frontdesk/main.go`, uncomment the `greeters` slice and
   the call to `library.WelcomeAll(greeters)`.
6. Run it. Confirm all three greetings print — then search every file
   you touched for the word `implements`. It shouldn't appear anywhere.
7. Run `go test ./...`. `internal/library` builds now: `TestVolunteerPromotion`
   and `TestGreeterSatisfiedByMultipleTypes` both pass.
   (`TestCheckOutNilInterfaceGotcha` is visible for the first time too,
   and fails — that's Exercise 5, next.)

**Key Learning:** A type satisfies an interface purely by having the
right method set. `Person`, `Kiosk`, and `Mascot` share no base type and
declare no relationship to `Greeter` or to each other — this is
structurally identical to duck typing, checked at compile time. That
same compile-time check is exactly why this package's tests couldn't
build until this exercise was done: an interface with no methods is a
real, valid type, and nothing calling `.Greet()` through it can compile
until the method exists.

---

## Exercise 4: Handle anything with a type switch

**Objective:** Use `any` and a type switch for the case where even a
one-method interface is too narrow.

**Context:** `starter/cmd/frontdesk/main.go` has a `logCheckIn` function
with a placeholder body. It takes an `io.Writer` as its first parameter
instead of printing straight to stdout — same reason as Topic 2's "`fmt`
doesn't just print to the terminal" slide — which is what lets
`TestLogCheckIn`, in `starter/cmd/frontdesk/main_test.go`, check it with
a `bytes.Buffer` instead of spawning the real program. This exercise
lives in its own package, so its test has been failing cleanly and
independently the whole time, unaffected by Exercises 1-3.

**Tasks:**

1. Run `go test ./...`. `TestLogCheckIn` fails — the placeholder writes
   the raw event with no type switch, so none of the expected words show
   up in the output.
2. Implement `logCheckIn` with a type switch (`switch v := x.(type)`)
   that handles at least `int` (treat it as a visitor count), `string`
   (treat it as a patron name), and a `default` case for anything else —
   writing through `fmt.Fprintf(w, ...)` to the `w` it's given, not
   `fmt.Printf`.
3. Run `go run ./cmd/frontdesk`. Confirm the three `logCheckIn` calls
   already in `main()` (an `int`, a `string`, and a `float64`) each hit
   the case you'd expect.
4. Re-run `go test ./...` and confirm `TestLogCheckIn` passes.

**Key Learning:** Inside each `case`, the switch variable is re-typed to
that case's concrete type — inside `case int`, `v` is an `int`, not an
`any` you have to cast yourself.

---

## Exercise 5: Fix the nil-pointer-in-an-interface gotcha

**Objective:** Reproduce a bug that catches senior engineers, understand
why it happens, then fix the underlying pattern.

**Context:** `starter/internal/library/checkout.go` has a `CheckOut`
function that returns `error`. As written, it declares a `*CheckoutError`
variable, leaves it `nil` on the success path, and returns it directly.
`TestCheckOutNilInterfaceGotcha` in `checkout_test.go` is already
written, and became visible for the first time once Exercise 3 got the
package building — it fails on its very first assertion.

**Tasks:**

1. Run `go test ./...`. `TestCheckOutNilInterfaceGotcha` fails
   immediately: even on a comfortably-under-the-limit call, `err != nil`
   is true. Run `go run ./cmd/frontdesk` too and look at the `CheckOut()`
   line in that output — same symptom, `err == nil` prints `false` on a
   clear success path. Before changing anything, explain to yourself (or
   out loud) why, in terms of the `(type, value)` pair an interface value
   really is.
2. Fix `CheckOut` in `starter/internal/library/checkout.go` so the
   success path returns a literal `nil` directly, instead of a variable
   declared as a concrete `*CheckoutError` pointer type.
3. Run it again. Confirm `err == nil` now prints `true` on the success
   path, and still returns a usable error when the patron is over their
   book limit.
4. Re-run `go test ./...` and confirm `TestCheckOutNilInterfaceGotcha`
   passes — every test in `starter/` should be green now.
5. Optional, worth doing once: temporarily put the bug back — declare
   `var problem *CheckoutError = nil` and `return problem` on the success
   path, instead of a literal `return nil`. Run `go test ./...` again and
   confirm `TestCheckOutNilInterfaceGotcha` fails. Then restore the fix
   and confirm it passes again — that's the test catching a regression
   the instant someone reintroduces the exact bug you just fixed.

**Key Learning:** An interface value is a `(type, value)` pair. A `nil`
pointer of a concrete type, returned through an interface-typed return
value, produces a pair like `(*CheckoutError, nil)` — the type half is
non-nil, so the whole pair isn't `nil`. The fix is structural: never let
a concrete pointer type leak into an interface return slot on the
success path — return a literal `nil` instead.

---

## Summary

By the end of this lab you should be able to:

- Embed a type with no field name and predict exactly which fields and
  methods get promoted
- Explain why shadowing a promoted method isn't polymorphic override, and
  write the manual call-through Go uses instead of `super`
- Make three unrelated types satisfy the same interface with no declared
  relationship, and explain why that's not a coincidence
- Use `any` with a type switch to handle values whose type isn't known
  ahead of time
- Diagnose and fix the nil-pointer-in-an-interface gotcha by explaining
  the `(type, value)` pair behind every interface value
- Explain why Go compiles a package's tests as a single unit, and read a
  batch of compiler errors spanning several files to work out which one
  belongs to the exercise you're on
- Read a failing test's output to work out what's still missing, and make
  it pass — without needing to write a test yourself yet
