---
title: "**Go Essentials**"
sub_title: Go Programming — Topic 1
author: Kevin Cunningham
---

## Accessing Resources

- Navigate to https://learn.neueda.com
- Enter the course code `181003`
- Create a username or login
- The registration code is `EVBYWI` 

<!-- end_slide -->

## Opening scenario

You've written a small script in Python. It works perfectly on your laptop.
You zip it up and send it to a teammate. They open it and — nothing runs.
Wrong Python version, missing packages, no virtualenv activated.

**Type in chat: how many extra steps would your teammate need before your script actually runs on their machine?**

We'll come back to this once we've built and shipped a Go program.

<!--
speaker_note: |
  Let answers land in chat for 20-30 seconds - most people will say
  "depends," "a few," or list interpreter/venv/pip install steps.
  Don't resolve it yet, just bank the number they gave.
-->

<!-- end_slide -->

## What Go actually is

Go is a compiled, statically typed language that produces a single,
self-contained binary.

If you've spent your career in Python (interpreted, dynamically typed) or
Java/TypeScript (compiled-ish, but to bytecode or to JS that still needs a
runtime), that sentence is doing more work than it looks like.

<!-- pause -->

**No interpreter. No JVM. No runtime on the target machine at all.**

A Go binary is statically linked by default — it needs nothing but itself.

<!-- end_slide -->

## Why Go looks the way it does

Go was designed at Google in 2007 by engineers frustrated with C++ build
times measured in tens of minutes, and with codebases where a decade of
teams had each picked their own way to do the same thing.

<!-- pause -->

**The design brief was small, on purpose.** Go's own documentation states
the goal directly: compile fast, read easily at scale, and give large
teams of engineers — different backgrounds, different experience levels —
one obvious way to write anything, not several.

<!--
speaker_note: |
  Frame this as the lens for the entire course, not a one-off aside.
  Every "Go doesn't have X" moment from here on (no overloading, no
  exceptions, no inheritance) traces back to this slide - point back to
  it explicitly the first few times it comes up.
-->

<!-- end_slide -->

## Simplicity is a constraint, not an accident

Go has roughly 25 keywords. Java has 50+. C++ has 90+, and the number
keeps growing with every revision of the standard.

<!-- incremental_lists: true -->

- One looping construct: `for`. No `while`, no `do-while`, no dedicated `foreach`.
- One way to format code: `gofmt` — not a house style guide, a compiler-adjacent tool.
- No function overloading, no operator overloading, no implicit numeric conversions.
- No class hierarchies — composition and interfaces instead, covered properly in Topic 5.

<!-- incremental_lists: false -->

**Every "no" on this list is a feature Go's designers looked at and
rejected on purpose** — not a gap they didn't get to, a boundary held
deliberately.

<!--
speaker_note: |
  Worth calling back to explicitly when "No function overloading" and
  "No exceptions for ordinary failure" land in Topic 2 - those slides
  explain what you do instead, this slide explains why the restriction
  exists in the first place.
-->

<!-- end_slide -->

## Consistency: one way to read any Go codebase

Rob Pike, one of Go's original designers, put it bluntly: **"gofmt's
style is no one's favorite, yet gofmt is everyone's favorite."**

<!-- pause -->

The point was never that tabs beat spaces, or that one brace style is
objectively correct. The point is that **every Go codebase you'll ever
open looks like every other one**, because none of that is left as a
per-team decision.

**Real-world payoff:** move between Go teams, or open an unfamiliar
open-source Go repo, and you skip the "whose style guide is this" ramp-up
entirely — something no Python or JavaScript team gets for free.

<!-- end_slide -->

## The trade-off, stated honestly

Simplicity isn't free. Go's smallness shows up as things that can feel
like missing features on your first week with the language:

<!-- incremental_lists: true -->

- More repetition — `if err != nil` at nearly every call site (Topic 2)
- No generics-driven abstraction culture the way Java or TypeScript have
  — real generics only arrived in 2022, and idiomatic Go still reaches
  for them sparingly (Topic 6)
- Verbosity in places a more expressive language would compress into one
  line

<!-- incremental_lists: false -->

**This course names every one of these trade-offs as they come up** — not
to relitigate them, but so you can explain *why* Go made that call the
first time someone on your team asks.

<!--
speaker_note: |
  Keep this evenhanded rather than defensive. Some delegates will
  genuinely prefer a more expressive language for their own projects,
  and that's a legitimate position - the goal is that they can name the
  trade-off precisely, not that they leave the room converted.
-->

<!-- end_slide -->

## The basic loop

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go")
}
```

<!-- column: 1 -->

```
go run main.go   # compile + run,
                 # throws away the binary

go build         # produces an executable
                 # named after the module

./your-binary    # run it directly —
                 # no runtime needed

go install       # builds, places binary
                 # in $HOME/go/bin
```

<!-- reset_layout -->

**Demo:** run all four commands against the same file and watch what each one leaves behind on disk.

<!--
speaker_note: |
  Run this live rather than showing pre-baked output. Have delegates
  watch the working directory before and after `go build` - the
  appearance of a binary named after the folder surprises people who
  expect an explicit output flag, like `gcc -o`.
-->

<!-- end_slide -->

## One package, one entry point

Every Go file belongs to a **package**, declared at the top with
`package something`.

An executable program needs exactly one package called `main`, and within
it, exactly one `func main()` — that's the entry point.

<!-- pause -->

No classes to instantiate. No `if __name__ == "__main__"`. No
`public static void main(String[] args)`. Just a function.

<!-- end_slide -->

<!-- jump_to_middle -->

Where Go disagrees with what you know
===

<!-- end_slide -->

## Capitalization is access control

In Java and TypeScript you write `public` or `private` explicitly. Go has
no such keyword.

Whether an identifier — function, type, variable, struct field — is
visible outside its package is decided entirely by whether its name
starts with an uppercase letter.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
func DoSomething() {}
// exported — visible to
// other packages
```

<!-- column: 1 -->

```go
func doSomething() {}
// unexported — package-
// private, full stop
```

<!-- reset_layout -->

**There's no "protected," no "internal," no friend classes. The package boundary is the only boundary.**

<!--
speaker_note: |
  Expect "where's public/private?" within the first ten minutes from
  anyone with a Java background. Export something live, then rename it
  lowercase and show the compile error from another package - watching
  it fail is what makes the rule stick, not just stating it.
-->

<!-- end_slide -->

## gofmt is not a linter suggestion

Run `gofmt -w .` (or let your IDE do it on save) and your code's
formatting is no longer up for discussion.

Tabs for indentation, brace placement, import grouping — all decided
for you.

<!-- pause -->

**Type in chat: coming from Python (PEP 8, broken constantly) or JS (a whole ecosystem of Prettier configs), is one enforced format a relief or a loss of control?**

<!--
speaker_note: |
  This usually splits the room a little - some people love not having
  the argument anymore, others resent losing stylistic choice. Both are
  reasonable reactions, let both surface before moving on.
-->

<!-- end_slide -->

## Semicolons exist, you just don't write them

Go's grammar requires semicolons; the lexer inserts them automatically at
the end of a line under certain rules.

The practical consequence: **opening braces must stay on the same line**
as the statement they belong to.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// Compiles:
func main() {
    fmt.Println("ok")
}
```

<!-- column: 1 -->

```go
// Does NOT compile:
func main()
{
    fmt.Println("nope")
}
```

<!-- reset_layout -->

**Demo:** break a working function by moving the brace to its own line and read the resulting error carefully — it won't say "wrong brace placement."

<!--
speaker_note: |
  If you've spent years with Allman-style braces in a C#/Java shop,
  this catches your fingers more than once. The error Go actually
  gives is confusing (something about a missing function body) - that
  confusion is the point of running it live.
-->

<!-- end_slide -->

## Unused things are compile errors

An unused local variable or an unused import is a hard compile failure
in Go — not a warning you can ignore.

Python will let unused imports slide. TypeScript, depending on config,
might too. Go won't build at all.

<!-- pause -->

**This bites hardest mid-demo:** comment out a line, forget its now-orphaned import is still there, and the whole build stops.

<!--
speaker_note: |
  Build in a beat for "yes, this is intentional, here's why" rather
  than treating it as a side note when it inevitably happens live.
-->

<!-- end_slide -->

## Why ship one binary at all?

Handing someone a recipe card versus handing them a vacuum-sealed,
fully cooked meal.

<!-- incremental_lists: true -->

- The recipe card (interpreted code) needs a kitchen, the right pans, the right skill, every time
- The vacuum-sealed meal (a compiled Go binary) just needs to be opened

<!-- incremental_lists: false -->

**This is exactly why Docker, Kubernetes, and Terraform are written in Go** — tools that run on thousands of differently-configured machines, where "ship one binary" massively simplifies distribution.

<!-- end_slide -->

## Setting up your IDE

VS Code with the official **Go extension** (`golang.go`, published by the
Go team) is the path of least resistance. It wires up:

<!-- incremental_lists: true -->

- `gopls` — the language server (autocomplete, go-to-definition, refactors)
- `gofmt` on save
- Inline `go vet` warnings
- Debugging and test running, built in

<!-- incremental_lists: false -->

GoLand is the alternative if you want something heavier, with more configuration out of the box.

<!-- end_slide -->

## Modules, not GOPATH

Modern Go organizes code with **modules**, not the old `GOPATH` convention.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```
mkdir my-project
cd my-project
go mod init example.com/my-project
```

<!-- column: 1 -->

Creates a `go.mod` file that names your module and pins your Go version and dependencies — the same job `package.json` does for a Node project.

<!-- reset_layout -->

**Demo:** `go mod init`, then open `go.mod` and point out what's in it before any code exists.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Compiled, statically typed, single binary**: no runtime required on the target machine
2. **Simplicity and consistency are deliberate design goals**: a small keyword set and one obvious way to do most things, not limitations Go hasn't gotten around to fixing
3. **Capitalization is the entire visibility mechanism**: no `public`/`private` keywords
4. **`gofmt` ends the formatting debate**: one enforced style, no config
5. **Braces stay on the statement line**: automatic semicolon insertion enforces it
6. **Unused vars and imports are compile errors**: not warnings you can ignore

<!-- end_slide -->

## Back to the opening scenario

Your Python script needed a teammate to have the right interpreter, the
right packages, and the patience to set both up — every time.

**A Go binary just needs to be run.** `go build`, hand over the file, done.

<!-- pause -->

**Type in chat: now that you've seen `go build`, how many of the steps you listed earlier actually go away?**

<!--
speaker_note: |
  This is the payoff for the opening chat poll - read a few of the
  original answers back if you can, then contrast with "zero setup
  steps" for the compiled binary. Don't let this become a Go-vs-Python
  debate, the point is the deployment model, not language superiority.
-->

<!-- end_slide -->

## Bridge to Topic 2

**We've established:**

<!-- incremental_lists: true -->

- Go compiles to one dependency-free binary
- Simplicity and consistency are the design goals behind nearly every
  restriction you'll hit from here on — expect this course to name the
  trade-off each time one shows up
- Visibility is controlled by capitalization, not keywords
- Formatting and unused-code rules are enforced by the toolchain, not convention

<!-- incremental_lists: false -->

**Topic 2: Core Language Features** — variables, types, functions, deferred
execution, formatted I/O, and how packages are organized within a module.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
