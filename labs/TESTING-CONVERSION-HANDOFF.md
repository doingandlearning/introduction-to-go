# Handoff: converting labs to pre-written failing tests

Working doc for continuing the test-pattern conversion started in Lab 2.
Point a fresh session at this file plus `labs/02-core-language-features/`
and `labs/03-flow-control-and-data-structures/` (both finished, second
template) and `labs/04-object-oriented-programming/` (next up).

## The decision (confirmed with Kevin)

Every implementation exercise, in every lab, should ship a **pre-written,
already-failing test** in `starter/` instead of the old "implement it,
then write your own test to prove it" pattern. Students never write a
test themselves before Topic 12 (Testing) — that's explicitly deferred.
Topic 2's slides ("Proving it: your first failing test") already teach
this framing: run `go test ./...` before touching anything, watch it
fail, implement until it's green.

## Scope

- **Converted:** Topic 2 slides (`slides/02-core-language-features/slides.md`), Lab 2 (`labs/02-core-language-features/`) — the template. Lab 3 (`labs/03-flow-control-and-data-structures/`) — done this session, see notes below.
- **Still to do:** Labs 4, 5, 6, 7, 8, 9, 10, 11.
- **Do not touch:** Lab 12 (Testing) — students learn to *write* tests here for the first time, on purpose, via a different pattern (production code finished, `t.Skip` test stubs to fill in). Lab 13 (Docker) — no unit-testable code, exercises are about Dockerfiles/compose/WSL.

## The pattern, exactly (from Lab 2)

For each exercise that has the student implement a function:

1. Write the test **fully** in `starter/.../X_test.go` (not a `t.Skip`
   stub — a complete, real test using the exact assertions from the old
   solution test, if one already existed). It must fail against the
   current `starter/` placeholder and pass against `solution/`.
2. If the function under test prints straight to stdout instead of
   returning/writing a checkable value (common in `cmd/*/main.go` demo
   functions), refactor it to accept an `io.Writer` and use
   `fmt.Fprintf`/`Fprintln` instead of `fmt.Printf`/`Println`. Update
   `main()` to pass `os.Stdout`. Tests then use a `bytes.Buffer` and
   check for substrings — favor loose "is this label/value present"
   checks over exact-format assertions, since a reasonably-correct
   student implementation may not match one exact byte-for-byte layout.
   (This is also a deliberate callback to Topic 2's "`fmt` doesn't just
   print to the terminal" slide — name that connection in the exercise
   text if it fits.)
3. Copy the same test file into `solution/` (it should already pass
   there, since solution code is finished).
4. Rewrite the exercise's **Tasks** list: first task is "run `go test
   ./... -v`, read the failure," then the implementation steps, then "re-run
   `go test ./... -v`, confirm it passes." Keep any `go run` verification
   step too where it adds value — the test isn't a replacement for seeing
   the real program work, just the primary check.
5. If the lab has a trailing "Prove it with a test" / "write your own
   test" exercise, **delete it** — its content (what to assert) gets
   folded into the exercises that already existed, as their pre-written
   test.
6. Exercises with nothing to implement (pure prediction/observation, or
   deliberately-broken standalone demos like Lab 7's deadlock/race
   programs) stay manual — don't force a test onto them. Lab 1's
   "break the build on purpose" and Lab 2's "predict the defer order"
   are the precedent.
7. Update the lab's intro paragraph (one sentence: tests ship pre-written
   and failing, run `go test ./... -v` first) and the **Summary** section's
   last bullet (swap "write a test" language for "read a failing test and
   make it pass").

## Known gotcha — check this on every lab before calling it done

**Unused imports in starter placeholders.** If a starter file imports a
package (e.g. `"fmt"`, `"strings"`) that only the *finished* version of a
function actually calls, and the placeholder body doesn't use it yet,
`go test`/`go build`/`go vet` will fail with "imported and not used" —
before the student even gets to see their test fail properly. This bit
Lab 2's `mathutils.go` (`fmt` imported but unused until `SafeDivide` was
implemented).

Fix: don't pre-import it. Have the student add the import themselves as
part of the exercise (Lab 1's `strings.ToUpper` and Lab 2's `fmt.Errorf`
now both do this — TODO comment says exactly which import line to add).

**How to check, per lab, before finishing:** for every `starter/**/*.go`
file, confirm every imported package is actually referenced somewhere
outside of comments in that file's current (unfinished) state. A quick
pattern (adjust path):

```bash
for f in $(find starter -name '*.go'); do
  echo "== $f =="
  # list imports, then grep for "<pkg>." usage outside comment lines
done
```

(See the exact script used for Lab 2 in this session's history if you
need the working version — it's a simple grep loop, nothing fancy.)

## Other constraints to remember

- **No Go toolchain in this sandbox.** Nothing here can be compiled or
  run — every check is static: brace/paren balance, import-usage
  scanning, manual tracing of what a test asserts vs. what a placeholder
  currently returns. Say this explicitly when handing work back to Kevin,
  and ask him to run `go build ./... && go test ./...` for real on each
  converted lab before using it live.
- **Favor loose test assertions** (substring/label checks) over exact
  output matching wherever the exercise instructions leave any room for
  reasonable variation in student implementation — reduces false failures
  you can't verify by actually running the test.
- **Lab-specific heads-up:**
  - **Lab 7 (Concurrency):** the existing "prove it with a test" exercise
    already asserts an aggregate outcome (sorted results, a count) rather
    than exact output — that's the right shape, just needs to become
    pre-written instead of student-authored. The `standalone/deadlock`
    and `standalone/race` programs are deliberately-broken demos, not
    implementation exercises — leave them manual.
  - **Lab 11 (gRPC):** `bookingpb` is hand-written, explicitly labeled
    "ILLUSTRATIVE — NOT GENERATOR OUTPUT" (protoc never ran in this
    sandbox). Be careful that any new test still makes sense against that
    stand-in package, and don't let it imply the illustrative code is
    real generated output.
  - **Lab 10 (REST):** has `httptest`-based patterns already in mind per
    its slides' summary ("the handler is testable with
    `net/http/httptest`") — lean into that rather than reinventing an
    approach.

## Suggested order

Labs 3 → 4 → 5 → 6 → 8 → 9 → 10 → 7 → 11, i.e. save Concurrency and gRPC
for last since they're the least mechanical (aggregate assertions,
illustrative stand-in code). Do one lab at a time, verify with the static
checks above, then stop and report status rather than batching several
labs silently — this conversion has already turned up one real bug
(Lab 2's unused import) that only showed up once a human actually ran
`go test`, so keep the loop tight.

## Lab 3 — done this session, notes for the pattern going forward

Lab 3 had two wrinkles Lab 2 didn't, worth knowing before Lab 4:

- **Exercise 1 (loops) had zero tests before this session** — the old
  lab's trailing "prove it with a test" exercise only ever covered
  Exercises 2 and 3 (switches, struct semantics). Per the "every
  implementation exercise" rule, Exercise 1 got a brand-new
  `loops_test.go` (`TestPrintCatalog`, `TestCountAvailable`,
  `TestFindFirstAvailable`) even though nothing in the old lab tested it.
  Don't assume the old "prove it with a test" exercise's scope is the
  full set of exercises that need tests — check every implement-a-function
  exercise independently.
- **`PrintCatalog` printed straight to stdout via `fmt.Printf`** — same
  gotcha as Lab 2's `cmd/*/main.go` functions. Refactored to
  `PrintCatalog(w io.Writer, books []Book)` using `fmt.Fprintf`, updated
  `main()` to pass `os.Stdout` (which required adding an `"os"` import to
  `main.go` — that file isn't a TODO placeholder itself, so pre-importing
  it there is fine, unlike the mathutils.go-style gotcha).
- **Exercises 4, 5, 6 stayed manual, deliberately** — all three are
  inline `TODO` blocks directly in `cmd/librarian/main.go`'s `main()`,
  not separate functions with signatures. Ex4 (slice-aliasing bug) and
  Ex6 (map iteration order) are pure observation, same as Lab 2's
  `auditReport`/`closeZonesBuggy` precedent. Ex5 (word-frequency +
  comma-ok) does have real logic (building a map) that could in
  principle be extracted into a testable `catalog` function, but doing
  so wasn't requested and would mean inventing a new function signature
  not in the original lab — left manual to avoid scope creep. Revisit
  this call if Kevin wants Ex5 test-covered too; extracting a
  `catalog.CountCheckouts(log string) map[string]int` would be the
  natural approach, following the same io.Writer-refactor spirit.
- Old Exercise 7 ("Prove it with a test") deleted; its assertions were
  already word-for-word what `solution/internal/catalog/decisions_test.go`
  and `structs_test.go` had, so they were reused directly as the
  pre-written starter tests rather than rewritten from scratch.

All Lab 3 changes were pushed to
`/Users/kevincunningham/code/neueda/go-programming-reup/course/labs/03-flow-control-and-data-structures/`
on Kevin's machine via the device bridge. Static checks (brace/paren
balance, import-usage scan) all passed — as always, no Go toolchain in
this sandbox, so ask Kevin to run `go build ./... && go test ./...` in
both `starter/` and `solution/` before treating Lab 3 as live.
