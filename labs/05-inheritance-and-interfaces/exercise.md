# Lab 5: Inheritance and Interfaces — a library front-desk system

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

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

**Tasks:**

1. Open `starter/internal/library/library.go`. Embed `Person` in
   `Volunteer` — no field name, just the type name on its own line —
   alongside the existing `ShiftHours int` field.
2. In `starter/cmd/frontdesk/main.go`, uncomment the `Volunteer` block in
   `main()`.
3. Run `go run ./cmd/frontdesk`. Confirm `v.Name` prints `"Sam"` and
   `v.Greet()` produces a full greeting — neither one has any code on
   `Volunteer` itself yet.

**Key Learning:** Embedding a type with no field name promotes its
exported fields and methods onto the outer type automatically. This is
composition with syntactic sugar, not inheritance — `Volunteer` *contains*
a `Person`.

---

## Exercise 2: Shadow Greet(), then call through explicitly

**Objective:** See the difference between shadowing and polymorphic
override, and write Go's equivalent of a `super` call.

**Tasks:**

1. In `starter/internal/library/library.go`, give `Volunteer` its own
   `Greet()` method. It should call `v.Person.Greet()` explicitly and
   append `" I'm volunteering today."` to the result.
2. Run `go run ./cmd/frontdesk` again. Confirm `v.Greet()` now includes
   the volunteering line, while `v.Person.Greet()` — called directly —
   still returns the original, un-appended greeting.

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

**Tasks:**

1. Define the `Greeter` interface with one method: `Greet() string`.
2. Implement `Greet()` on `Kiosk` (e.g. return
   `"SCAN YOUR CARD AT STATION <N> TO BEGIN."`).
3. Implement `Greet()` on `Mascot` (e.g. return something like
   `"*<Name> waves enthusiastically, says nothing*"`).
4. In `starter/cmd/frontdesk/main.go`, uncomment the `greeters` slice and
   the call to `library.WelcomeAll(greeters)`.
5. Run it. Confirm all three greetings print — then search every file
   you touched for the word `implements`. It shouldn't appear anywhere.

**Key Learning:** A type satisfies an interface purely by having the
right method set. `Person`, `Kiosk`, and `Mascot` share no base type and
declare no relationship to `Greeter` or to each other — this is
structurally identical to duck typing, checked at compile time.

---

## Exercise 4: Handle anything with a type switch

**Objective:** Use `any` and a type switch for the case where even a
one-method interface is too narrow.

**Context:** `starter/cmd/frontdesk/main.go` has a `logCheckIn(x any)`
function with a placeholder body.

**Tasks:**

1. Implement `logCheckIn` with a type switch (`switch v := x.(type)`)
   that handles at least `int` (treat it as a visitor count), `string`
   (treat it as a patron name), and a `default` case for anything else.
2. Run `go run ./cmd/frontdesk`. Confirm the three `logCheckIn` calls
   already in `main()` (an `int`, a `string`, and a `float64`) each hit
   the case you'd expect.

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

**Tasks:**

1. Run `go run ./cmd/frontdesk` as-is and look at the `CheckOut()` line
   in the output. Even on a clear success path, `err == nil` prints
   `false`. Before changing anything, explain to yourself (or out loud)
   why — in terms of the `(type, value)` pair an interface value really
   is.
2. Fix `CheckOut` in `starter/internal/library/checkout.go` so the
   success path returns a literal `nil` directly, instead of a variable
   declared as a concrete `*CheckoutError` pointer type.
3. Run it again. Confirm `err == nil` now prints `true` on the success
   path, and still returns a usable error when the patron is over their
   book limit.

**Key Learning:** An interface value is a `(type, value)` pair. A `nil`
pointer of a concrete type, returned through an interface-typed return
value, produces a pair like `(*CheckoutError, nil)` — the type half is
non-nil, so the whole pair isn't `nil`. The fix is structural: never let
a concrete pointer type leak into an interface return slot on the
success path — return a literal `nil` instead.

---

## Exercise 6: Prove it with a test

**Objective:** Turn the eyeball checks from Exercises 1–5 into automated
tests — the same `go test` habit from Topic 2, now aimed at embedding,
interface satisfaction, and the nil-interface gotcha.

**Context:** `starter/internal/library/` has three test files already
scaffolded: `library_test.go`, `greeter_test.go`, and `checkout_test.go`.
Each has a real test function signature in place; the body is a
`t.Skip(...)` placeholder for you to replace.

**Tasks:**

1. In `library_test.go`, implement `TestVolunteerPromotion`. Build a
   `Volunteer` with an embedded `Person` and assert:
   - `v.Name` reads through to the embedded `Person.Name` (promotion).
   - `v.Greet()` differs from `v.Person.Greet()` (shadowing) and equals
     the embedded greeting plus `" I'm volunteering today."` (the manual
     call-through from Exercise 2).
2. In `greeter_test.go`, implement `TestGreeterSatisfiedByMultipleTypes`.
   Put a `Person`, a `Kiosk`, and a `Mascot` into a single `[]Greeter`
   slice, loop over it, and assert each `Greet()` call returns the string
   you'd expect — with no shared base type and no `implements` anywhere.
3. In `checkout_test.go`, implement `TestCheckOutNilInterfaceGotcha`.
   Call `CheckOut` once with room under the patron's limit and confirm
   `err == nil` is genuinely `true`; call it again at the limit and
   confirm you get a non-nil error whose `Error()` message names the
   book.
4. Run `go test ./...` from the lab root. All three tests should pass.
5. Deliberately break it: in `checkout.go`, temporarily put the
   Exercise 5 bug back — declare `var problem *CheckoutError = nil` and
   `return problem` on the success path, instead of a literal `return
   nil`. Run `go test ./...` again and confirm `TestCheckOutNilInterfaceGotcha`
   fails. Then restore the fix and confirm it passes again.

**Key Learning:** A test doesn't just check a return value, it can check
the *shape* of one — whether an interface holds a genuinely nil pair or a
`(type, value)` pair that only looks nil at a glance. `if err == nil`
inside a test is exactly the same trap, and the same fix, as `if err ==
nil` in production code. A test that covers the nil-interface gotcha will
catch a regression the instant someone reintroduces it, long after
everyone's forgotten Exercise 5 happened.

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
- Write tests that verify embedding promotion, interface satisfaction
  across unrelated types, and the nil-pointer-in-an-interface gotcha —
  and watch them fail on purpose before trusting them to pass
