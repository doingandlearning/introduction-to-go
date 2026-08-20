# Lab 1: Go Essentials — a two-package greeter CLI

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

Both directories contain the same shape: a `cmd/greeter` package (the
program entry point) and an `internal/greeting` package (a helper package
you'll import). By the end, you'll have run every command in the Go
toolchain you saw in the lecture, and broken the build on purpose at
least twice.

---

## Exercise 1: Get it running

**Objective:** Confirm your toolchain works before writing anything.

**Tasks:**

1. From `starter/`, run `go run ./cmd/greeter`. It should compile — but the
   output won't be right yet, because `greeting.go` isn't finished.
2. Run `go build ./cmd/greeter` — **name the package explicitly, not
   `./...`**. `starter/` has two packages (`cmd/greeter` and
   `internal/greeting`), and `go build ./...` matches both — when `build`
   compiles more than one package it discards the result instead of
   writing a binary, silently and without error. Naming a single main
   package is what actually produces a file.
3. Run the binary it just built: `./greeter` on macOS/Linux,
   `.\greeter.exe` on Windows (PowerShell's native path separator —
   `./greeter.exe` usually works too, but `.\` is the safer, idiomatic
   form for a native `.exe`, not just a PowerShell command). Delete it
   once you've confirmed it runs.
4. Run `go install ./cmd/greeter`. Find the binary it produced — check
   `$HOME/go/bin` (or `%USERPROFILE%\go\bin` on Windows). Run it from
   there, outside the project directory entirely.

**Key Learning:** `go run`, `go build`, and `go install` all compile your
code — the difference is what they do with the result afterward, not
whether they "interpret" it. There's no such thing as Go interpreting a
file. And `go build ./...` is for *checking* a whole module compiles —
not for producing a binary, unless it happens to match exactly one main
package. Name the package explicitly when you actually want the file.

---

## Exercise 2: Finish the greeting package

**Objective:** Implement the exported function the CLI depends on.

**Context:** `starter/internal/greeting/greeting.go` has a `TODO` where
`Greet` should live. `starter/cmd/greeter/main.go` already calls
`greeting.Greet(name)` and expects a string back.

**A note before you start:** the function signatures are already written
for you — you're only filling in the bodies. Doing that needs a `return`
statement and string concatenation with `+`, neither of which the lecture
has covered yet on purpose. Both work exactly like they do in Python,
Java, or C, so lean on what you already know here; Go's function syntax
proper (multiple returns, named returns, variadics) is Topic 2's job.

**Tasks:**

1. Open `starter/internal/greeting/greeting.go`. Implement `Greet(name
   string) string` so it returns `"Hello, <name>!"`.
2. Add a second, **unexported** helper function `shout(s string) string`
   that returns `s` uppercased, and call it from `Greet` so the greeting
   comes back shouted: `"HELLO, <NAME>!"`.
3. Run `go run ./cmd/greeter` again from `starter/`. Confirm the greeting
   is correct.
4. Now, from `cmd/greeter/main.go`, try calling `greeting.shout("test")`
   directly instead of going through `Greet`. What error do you get, and
   why does it name the specific thing you can't do?

*Hint if stuck:* `Greet` needs one line — `return shout("Hello, " + name +
"!")` — string concatenation with `+` works the same as it does in
Python, Java, or C. `shout` needs `strings.ToUpper(s)`, which means
adding a second import line above `package greeting`, the same shape as
the `import "fmt"` you've already seen in Topic 1 — just point it at
`"strings"` instead.

**Key Learning:** Capitalization is the only visibility rule Go has. An
unexported name isn't just "discouraged" from outside its package — it's
genuinely invisible, and the compiler enforces it, not a linter.

---

## Exercise 3: Same package, no import needed

**Objective:** See when Go does — and doesn't — require an import.

**Tasks:**

1. In `starter/cmd/greeter/`, create a second file, `farewell.go`, in the
   same `package main` as `main.go`.
2. In `farewell.go`, write an unexported function `farewell(name string)
   string` that returns `"Goodbye, <name>."` — no import needed, since
   it's the same package as `main.go`.
3. Call `farewell(name)` from `main()` in `main.go` and print the result
   after the greeting.
4. Run it. Confirm both lines print.

**Key Learning:** Package membership, not file boundaries, is what
requires an import. Two files in the same directory with the same
`package` declaration share everything unexported between them for free.

---

## Exercise 4: Break the build, on purpose

**Objective:** Learn to read Go's compiler errors by causing three of
them deliberately, one at a time. Fix each before moving to the next.

**Tasks:**

1. Remove the `package main` line from `main.go`. Run `go build ./...`.
   Read the error. Put the line back.
2. Move the opening `{` of `func main()` onto its own line (Allman
   style). Run `go build ./...` again. The error will not obviously say
   "brace in the wrong place" — that's automatic semicolon insertion
   biting you. Figure out what it's actually telling you, then fix it.
3. Add an import you never use anywhere in the file (e.g. `"strings"`
   with no calls to it). Run the build. Note that this is a hard failure,
   not a warning — then remove the unused import.

**Key Learning:** Go's compiler errors are precise about what it expected,
even when the underlying cause (like semicolon insertion) isn't obvious
from the error text alone. Reading them line-by-line beats guessing.

---

## Exercise 5: Format and vet

**Objective:** Use the two tools that check things the compiler doesn't.

**Tasks:**

1. In any file you've edited, deliberately mis-indent a block (mix tabs
   and spaces, or misalign a brace). Run `gofmt -l .` from the project
   root — it should list the file. Run `gofmt -w .` and look at the diff.
2. Add a `fmt.Printf` call somewhere with a mismatched verb — e.g.
   `fmt.Printf("Age: %d\n", "not a number")`. Run `go build ./...` (it
   still compiles). Then run `go vet ./...` and read what it flags.
   Remove the bad Printf afterward.

**Key Learning:** `gofmt` and `go vet` catch different classes of problem
than the compiler does — formatting consistency and probable mistakes,
respectively, neither of which stop your program from building.

---

## Summary

By the end of this lab you should be able to:

- Explain the difference between `go run`, `go build`, and `go install`
- Export or hide an identifier deliberately, and predict the compile
  error when something unexported is reached from outside its package
- Explain why two files can share a package with no import between them
- Read a Go compiler error and locate the actual cause, even when the
  message doesn't name it directly
- Use `gofmt` and `go vet` for the two classes of problem the compiler
  doesn't check
