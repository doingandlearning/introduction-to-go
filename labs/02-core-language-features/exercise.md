# Lab 2: Core Language Features — a warehouse inventory CLI

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

**Every exercise below ships with its test already written**, sitting in
`starter/internal/mathutils/mathutils_test.go` and
`starter/cmd/inventory/main_test.go`. Run `go test ./...` from `starter/`
right now, before you change anything — every test fails. That's the
starting line, not a bug. Your job in each exercise is to make the named
test pass, not to write a new one; writing tests yourself is Topic 12's
job.

Both directories contain the same shape: a `cmd/inventory` package (the
program entry point) and an `internal/mathutils` package (a helper
package you'll import). The scenario: a small tool that reports on a
warehouse's stock — total units, average unit price, and a per-item
formatted breakdown. By the end, you'll have exercised multiple return
values, zero values, `defer`, a two-package import, and the core `fmt`
verbs.

---

## Exercise 1: Divide safely

**Objective:** Write a function that reports failure as a value instead
of crashing — and watch a pre-written test go from red to green as you do.

**Context:** `starter/internal/mathutils/mathutils.go` has a `TODO` where
`SafeDivide` should live. `TestSafeDivide` in `mathutils_test.go` is
already written and already failing. The inventory report needs
`SafeDivide` to compute average unit price (`totalValue / totalQuantity`)
without panicking when a warehouse has zero items in stock.

**Tasks:**

1. Run `go test ./...` from `starter/`. Read `TestSafeDivide`'s failure —
   it's checking a normal division and a divide-by-zero case, and the
   placeholder implementation satisfies neither.
2. Implement `SafeDivide(a, b float64) (float64, error)` in
   `starter/internal/mathutils/mathutils.go`. Return an error (via
   `fmt.Errorf`, not `errors.New`, so you get formatting practice) when
   `b == 0`, otherwise return `a / b` and a `nil` error. The file doesn't
   import `"fmt"` yet — add that line yourself; it's not there until
   something in the file actually needs it.
3. Re-run `go test ./...`. `TestSafeDivide` should now pass.
4. Run `go run ./cmd/inventory` with the provided sample data and confirm
   you get a sensible average — the same behavior the test just checked,
   now visible in the real program.

**Key Learning:** A function that can fail returns `error` as its last
value. There's no exception to catch — the caller checks `err != nil`
immediately, right where the call happens. The test told you exactly
what "correct" meant before you wrote a line of implementation.

---

## Exercise 2: Zero values, not null

**Objective:** See Go's deterministic zero values firsthand, and make a
pre-written test pass by printing them correctly.

**Context:** `starter/cmd/inventory/main.go` has a `TODO` inside
`reportZeroValues` — a function meant to declare one variable of each
basic type without assigning it, plus an unassigned `Item` struct, and
print all of them to the `io.Writer` it's given. `TestReportZeroValues`
in `main_test.go` checks that a label for each one shows up in the
output.

**Tasks:**

1. Run `go test ./...`. `TestReportZeroValues` fails — the placeholder
   only prints the header line.
2. Declare unset variables of type `int`, `float64`, `string`, `bool`,
   and `*Item` (a pointer). Print each with `%v` to the function's `w
   io.Writer` parameter (`fmt.Fprintf(w, "int:      %v\n", qty)`, and so
   on), labeled the way the test expects: `"int:"`, `"float64:"`,
   `"string:"`, `"bool:"`, `"*Item:"`.
3. Declare an unset `Item` (the struct already defined in this file —
   `Name string`, `Quantity int`, `UnitPrice float64`) and print it with
   `%+v`, labeled `"Item:"`.
4. Re-run `go test ./...` and confirm `TestReportZeroValues` passes, then
   `go run ./cmd/inventory` and compare the printed zero values to what
   you'd expect from `None` in Python or `null` in Java.

**Key Learning:** Every Go type has a deterministic zero value —
`0`, `""`, `false`, `nil` for pointers — and a zero-value struct is
immediately usable, not a null reference waiting to blow up. Passing an
`io.Writer` in instead of calling `fmt.Println` directly is what let this
be tested at all without spawning the real program.

---

## Exercise 3: Predict the defer order

**Objective:** Reason about `defer`'s LIFO ordering before you run it.

**Context:** `starter/cmd/inventory/main.go` has a function,
`auditReport`, with three `defer` statements already written — nothing to
implement here, and no test either. This one's purely for observation.

**Tasks:**

1. Read `auditReport` without running it. Write down, on paper or in a
   comment, the order you expect its three deferred lines to print in.
2. Run `go run ./cmd/inventory`. Compare the actual output to your
   prediction.
3. Add a fourth `defer` statement to the function and predict again
   before re-running.

**Key Learning:** Deferred calls run in LIFO order — last deferred, first
executed — when the surrounding function returns, regardless of how many
statements sit between the `defer` and the return.

---

## Exercise 4: Fix the loop-variable defer gotcha

**Objective:** Reproduce a real `defer` bug, then fix it — verified by a
pre-written test checking the exact LIFO order.

**Context:** `starter/cmd/inventory/main.go` has a function,
`closeZonesBuggy`, that loops over a slice of warehouse zone names and
defers a `fmt.Println` for each one — reproducing the classic gotcha.
`closeZonesFixed` is your job, and `TestCloseZonesFixed` in
`main_test.go` is already written, checking that `"closing zone: C"`,
`"closing zone: B"`, and `"closing zone: A"` come out in exactly that
order.

**Tasks:**

1. Run `closeZonesBuggy` (called from `main`) and observe the printed
   order. It won't be the order you'd naively expect from reading the
   loop top to bottom.
2. Run `go test ./...` — `TestCloseZonesFixed` fails, since the
   placeholder only prints the header line.
3. Implement `closeZonesFixed(w io.Writer)`: loop over the same zones
   slice, deferring a closure that takes the zone name as an explicit
   parameter and prints it to `w` —
   `defer func(zone string) { fmt.Fprintln(w, "closing zone:", zone) }(z)`.
4. Re-run `go test ./...` and confirm `TestCloseZonesFixed` passes, then
   `go run ./cmd/inventory` and confirm the printed order matches.

**Key Learning:** A deferred call's arguments are evaluated the moment
`defer` runs, not when the call executes. Passing the loop variable as a
parameter — instead of closing over it — is what locks in the value you
actually wanted, and it's exactly what the pre-written test is checking
for, line by line.

---

## Exercise 5: A second package

**Objective:** Import your own package across a module, the same way
you'd import `mathutils` from anywhere else in a larger codebase — and
make `TestAdd` pass in the process.

**Context:** `starter/internal/mathutils/mathutils.go` has a `TODO` for
`Add`, an exported function summing a slice of ints — used by
`cmd/inventory` to total up item quantities across the warehouse.
`TestAdd` in `mathutils_test.go` is already written and already failing.

**Tasks:**

1. Run `go test ./...`. `TestAdd` fails on every case, including
   `Add()` with no arguments at all — read what the placeholder returns
   versus what each case expects.
2. Implement `Add(nums ...int) int` in `mathutils.go`, using a variadic
   parameter and a `for range` loop to sum them.
3. Re-run `go test ./...` and confirm `TestAdd` passes for every case in
   the table, including the zero-argument one.
4. Run `go run ./cmd/inventory` and confirm the total matches what you'd
   get by summing the sample data by hand.

**Key Learning:** Import paths are rooted at the module name in
`go.mod`. `example.com/core-lab/internal/mathutils` is one import,
regardless of how many files live inside that package directory. Notice
that the zero-argument case was already in the test table before you
wrote a line of `Add` — a good implementation handles it without any
special-casing.

---

## Exercise 6: Compare the struct-formatting verbs

**Objective:** See exactly what each of `%v`, `%+v`, `%#v`, and `%T` is
for, side by side, against the same value — checked by a pre-written test.

**Context:** `starter/cmd/inventory/main.go` has a `TODO` inside
`printItemFormats` where all four verbs should run against one `Item`.
`TestPrintItemFormats` in `main_test.go` checks that all four verb labels
and the `%T` output (`main.Item`) show up.

**Tasks:**

1. Run `go test ./...` — `TestPrintItemFormats` fails, since the
   placeholder only prints the header line.
2. Print the `sample Item` parameter with `%v`, `%+v`, `%#v`, and `%T` to
   the function's `w io.Writer`, one per line, each labeled with the verb
   that produced it (`fmt.Fprintf(w, "%%v   -> %v\n", sample)`, and so on).
3. Re-run `go test ./...` and confirm `TestPrintItemFormats` passes, then
   `go run ./cmd/inventory` and read the four real lines. Note which one
   you'd reach for while skimming logs, which for debugging a specific
   field, which for pasting a literal back into source, and which for
   confirming a type at runtime.

**Key Learning:** These four verbs aren't redundant — each answers a
different question (`%v` "what's roughly in here," `%+v` "which field is
which," `%#v` "give me a compilable literal," `%T` "what type even is
this").

---

## Summary

By the end of this lab you should be able to:

- Write a function that reports failure via a returned `error`, and call
  it with the `if err != nil` pattern
- Explain why a zero-value struct needs no constructor
- Predict `defer`'s LIFO execution order before running the code
- Diagnose and fix the loop-variable `defer` gotcha
- Import a package you wrote yourself using a module-rooted import path
- Choose the right `fmt` verb (`%v`, `%+v`, `%#v`, `%T`) for a given task
- Read a failing test's output to work out what's still missing, and make
  it pass — without needing to write a test yourself yet
