# Handoff: converting labs to pre-written failing tests

Working doc for continuing the test-pattern conversion started in Lab 2.
Point a fresh session at this file plus `labs/02-core-language-features/`
and `labs/03-flow-control-and-data-structures/` (both finished, second
template), `labs/05-inheritance-and-interfaces/` (done a prior session —
**read its notes below before doing any lab with an interface-
satisfaction exercise**, it found a structural gotcha that generalizes),
`labs/06-functional-programming/` (done a prior session, notes below —
the cleanest conversion so far, no structural surprises, but it did turn
up a wrong expected value baked into the original lab), and
`labs/07-concurrency/` (done this session, notes below — the first lab
with a genuinely new gotcha: a stuck producer/consumer can't fail a test
normally, only hang it, so every new test wraps its wait in a
`select`+`time.After`). Lab 4 is still open and skipped three times in a
row now — Kevin asked to prioritize 5, then 6, then 7, since that's what
the course was covering next each time. It's first in line whenever
picked back up.

## The decision (confirmed with Kevin)

Every implementation exercise, in every lab, should ship a **pre-written,
already-failing test** in `starter/` instead of the old "implement it,
then write your own test to prove it" pattern. Students never write a
test themselves before Topic 12 (Testing) — that's explicitly deferred.
Topic 2's slides ("Proving it: your first failing test") already teach
this framing: run `go test ./...` before touching anything, watch it
fail, implement until it's green.

## Scope

- **Converted:** Topic 2 slides (`slides/02-core-language-features/slides.md`), Lab 2 (`labs/02-core-language-features/`) — the template. Lab 3 (`labs/03-flow-control-and-data-structures/`) — done a prior session, static checks only. Labs 5, 6, and 7 (`labs/05-inheritance-and-interfaces/`, `labs/06-functional-programming/`, `labs/07-concurrency/`) — all three done with a real Go toolchain available (see "Other constraints" below — "no Go toolchain" is no longer a safe assumption, check first).
- **Still to do:** Labs 4, 8, 9, 10, 11.
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

- **Go toolchain availability depends on the session/environment —
  check with `which go && go version` before assuming either way.**
  Labs 2 and 3 were converted with no toolchain available: every check
  was static (brace/paren balance, import-usage scanning, manual tracing
  of what a test asserts vs. what a placeholder currently returns), and
  Kevin was asked to run `go build ./... && go test ./...` for real
  before trusting them. Lab 5's session *did* have a real Go toolchain
  (in the cloud workspace, not on Kevin's machine via the device bridge —
  that bridge still has no `go`) and every claim in this doc about Lab 5
  is toolchain-verified, not inferred. If a future session has a
  toolchain: stage the lab into the workspace, edit there, run
  `gofmt -l . && go vet ./... && go test ./...` against both `starter/`
  and `solution/` for real, *then* push the verified files back via the
  device bridge — this caught a real bug in Lab 5 (see below) that
  static checking never would have.
- **Favor loose test assertions** (substring/label checks) over exact
  output matching wherever the exercise instructions leave any room for
  reasonable variation in student implementation — reduces false failures
  you can't verify by actually running the test.
- **Lab-specific heads-up:**
  - **Lab 7 (Concurrency):** done — see notes below. The old "prove it
    with a test" exercise's aggregate-outcome assertion (sorted results,
    a count) was exactly right and got reused verbatim in the folded-in
    Exercise 8 test. The `standalone/deadlock` and `standalone/race`
    programs are deliberately-broken demos, not implementation exercises
    — left manual, unchanged.
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

Originally 3 → 4 → 5 → 6 → 8 → 9 → 10 → 7 → 11 (save Concurrency and gRPC
for last since they're the least mechanical). Actual order so far has
been 3 → 5 → 6 → 7, with 5, 6, and 7 all jumping the queue ahead of 4 at
Kevin's request, since each was what the course was starting next at the
time. **Lab 4 is still open** and is next in line whenever picked back
up, then 8 → 9 → 10 → 11 as before. Do one lab at a time, verify
(toolchain if available, static checks if not), then stop and report
status rather than batching several labs silently — this conversion has
already turned up four real bugs this way (Lab 2's unused import, Lab
5's pre-existing `WelcomeAll` compile failure, Lab 6's wrong `// want`
comment for the `Filter` string case, Lab 7's unused-import/unused-var
gotcha hitting three separate `cmd/*` packages at once), all invisible
to a quick read-through, so keep the loop tight.

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
balance, import-usage scan) all passed — as always at the time, no Go
toolchain in that sandbox, so Kevin was asked to run
`go build ./... && go test ./...` in both `starter/` and `solution/`
before treating Lab 3 as live.

## Lab 5 — done this session, notes for the pattern going forward

This session had a real Go toolchain, so everything below is verified by
actually running `gofmt`, `go vet`, `go build`, and `go test` against
both `starter/` and `solution/` — not inferred. Lab 5 was picked (out of
suggested order) because Kevin's course was starting it next.

- **Found and fixed a real, pre-existing bug, unrelated to the test
  conversion itself.** `internal/library/greeter.go`'s `WelcomeAll`
  calls `g.Greet()` on a `[]Greeter` — but in `starter/`, the `Greeter`
  interface starts with zero methods (that's Exercise 3's TODO). That
  one line meant `go build ./...` failed from a completely untouched
  starter checkout, before any exercise, for any student — not a test
  problem, a "the lab doesn't run at all" problem. Nobody caught this in
  Labs 2/3's static-only conversion because there was no toolchain to
  build against. **Fix:** comment out `WelcomeAll`'s body (with a TODO
  telling students to uncomment it once `Greeter` has `Greet()`, and
  re-add the now-unused `"fmt"` import) as part of Exercise 3, exactly
  the same "don't pre-import/pre-reference what isn't implemented yet"
  principle as the known unused-import gotcha, just applied to a method
  call instead of an import. **Check every lab still on the list for
  this same shape** — any pre-built helper function that calls a method
  on an interface/type the student hasn't finished yet will break the
  same way, and only a real `go build` catches it.
- **Structural finding that applies beyond Lab 5: a Go package compiles
  its test files as one unit.** Exercises 1, 2, 3, and 5 all live in
  `internal/library`. Once `WelcomeAll` was fixed, `go build ./...`
  passed from the start — but `go test ./...` still couldn't build
  *that package's test binary* until Exercises 1, 2, and 3 were **all**
  done, because `greeter_test.go`'s `TestGreeterSatisfiedByMultipleTypes`
  assigns concrete types into a `[]Greeter` — a compile-time interface
  check, not a runtime assertion — and Go compiles every `_test.go` file
  in a package together. Verified by hand, stage by stage: after
  Exercise 1 alone, `library_test.go` compiles but the package still
  won't build (blocked by `greeter_test.go`); same after Exercise 1+2;
  only after Exercise 3 does the package build, at which point
  `TestVolunteerPromotion` and `TestGreeterSatisfiedByMultipleTypes` are
  both already green (silently, the whole time) and
  `TestCheckOutNilInterfaceGotcha` (Exercise 5) becomes visible for the
  first time, failing cleanly on its own. **This means any lab with more
  than one interface-satisfaction exercise sharing a package will have
  this same "batch of compiler errors until several exercises are done
  together" shape** — it can't be avoided without splitting types into
  separate packages, which is out of scope (changes the lab's file
  layout, wasn't asked for). The fix that *was* applied: name this
  explicitly in the lab's intro and in each affected exercise's Context,
  so it reads as an intentional Go lesson ("a package is the unit of
  compilation") rather than a broken test setup. `cmd/frontdesk`
  (Exercise 4, `logCheckIn`) is a separate package and was unaffected —
  its test fails and passes in clean isolation the whole time.
- **`logCheckIn` (Exercise 4) had no test before this session**, same
  situation as Lab 3's Exercise 1 — the old lab never covered it. Same
  fix as Lab 3's `PrintCatalog`: refactored to take `io.Writer` as its
  first parameter (`logCheckIn(w io.Writer, x any)`), call sites updated
  to pass `os.Stdout`, new `main_test.go` added to both `starter/` and
  `solution/` with loose substring checks (`"42"` + `"visitors"`,
  `"Priya"` + `"patron"`, `"3.14"`) rather than exact string matches.
- Old "Exercise 6: Prove it with a test" deleted. Its three test bodies
  were already fully written in `solution/` (not `t.Skip` stubs there,
  just missing from `starter/`) — reused verbatim as the pre-written
  `starter/` tests rather than rewritten, same as Lab 3's Exercise 7.
  Its final "deliberately break it again, confirm the test catches the
  regression" step was folded into Exercise 5's tasks as an optional
  last step, since it's specific to the nil-interface bug Exercise 5 is
  already about.

All Lab 5 changes were pushed to
`/Users/kevincunningham/code/neueda/go-programming-reup/course/labs/05-inheritance-and-interfaces/`
on Kevin's machine via the device bridge, after `gofmt -l .`,
`go vet ./...`, `go build ./...`, and `go test ./...` all ran clean
against `solution/` (fully green) and matched the documented behavior
exactly against `starter/` (build/test failures only where expected, at
every intermediate exercise stage, verified by hand-completing exercises
one at a time in a scratch copy). Safe to treat as live without a
re-verification pass — but worth a quick real run in the classroom
anyway, same as always.

## Lab 6 — done this session, notes for the pattern going forward

The cleanest conversion so far. All five exercises' TODOs live in one
file, `internal/orders/orders.go`, and every placeholder body already
returns a syntactically valid stub (`return 0`, `return items`
unmodified, `return &CoffeeOrder{}`, etc.) — no interface-satisfaction
compile hazard like Lab 5's, verified the same way: hand-completed
Exercise 1 alone in a scratch copy and confirmed only its test went
green, everything else failed on its own, no shared-package build
failures blocking anything. `cmd/pipeline/main.go`'s five exercise
functions are pure `go run`-facing demo/print code that call into
`internal/orders` — none of them had a TODO themselves (all the real
logic lives in the package), so no `io.Writer` refactor was needed
anywhere in this lab, unlike Labs 3 and 5.

- **Found a second real bug, this time in a test oracle, not the
  production code.** The original lab's `cmd/pipeline/main.go` had
  `// want [latte espresso chai]` next to the `Filter(drinkNames, len >
  3)` call — but `"cola"` has length 4, so it passes `len(s) > 3` too.
  The correct filtered result is four items
  (`[latte cola espresso chai]`), not three. This comment was wrong in
  the *original* lab, before any conversion touched it — I copied it
  faithfully into `TestFilter`'s `wantLong` on the first pass, ran the
  test against a known-correct `solution/`, and it failed. That's what
  caught it: a real toolchain checks the test author's assumptions, not
  just the student's implementation. Fixed in three places: the test's
  `wantLong` slice, and the `// want` comment in both `starter/` and
  `solution/`'s `main.go`. **Lesson for future labs:** don't trust an
  existing `// want` comment as ground truth just because it predates
  this conversion — recompute it by hand (or run the real solution) 
  before writing it into a pre-written test, the same way you'd verify
  any other assertion.
- All six tests (`TestNewOrderCounter`, `TestFilter`, `TestDollars`,
  `TestMap`, `TestNewCoffeeOrder`, `TestInvoiceTotal`) live in the one
  `orders_test.go` file, matching Lab 2's `mathutils_test.go` shape
  (several small test functions in one file) rather than Lab 5's
  one-file-per-exercise split — whichever shape matches how the
  existing `_test.go` files in a given lab are already organized is the
  one to follow; don't invent a new split.
- **`TestNewCoffeeOrder` already existed almost verbatim** as the old
  "Exercise 6: Prove it with a test" — same fields, same three cases
  (zero/one/stacked options). Reused directly as Exercise 4's
  pre-written test rather than rewritten, same as Lab 3 and Lab 5's
  precedent for an old trailing test exercise whose assertions already
  match an earlier implementation exercise.
- **`Invoice.Total()` had a `TODO` that no exercise's task list actually
  assigned** in the original lab — Exercise 5's Context paragraph just
  asserted it would be "already implemented by the time you reach this
  exercise," without saying by whom or when. Folded implementing it into
  Exercise 5's own task list (as task 2, before the method-value
  assignment in `main.go`) rather than leaving that gap for a student to
  stumble on. `TestInvoiceTotal` checks both `Total()`'s arithmetic and
  the method-value assignment pattern itself, so it also stands in for
  what Exercise 5 can't get a compiler to check (that the student wrote
  `inv.Total` and not `inv.Total()` in `main.go` — that's a coding-style
  point the test can't see, only `go run`'s printed output confirms it).

All Lab 6 changes were pushed to
`/Users/kevincunningham/code/neueda/go-programming-reup/course/labs/06-functional-programming/`
on Kevin's machine via the device bridge, after `gofmt -l .`,
`go vet ./...`, `go build ./...`, and `go test ./...` all ran clean
against `solution/` (fully green, including a `go run` sanity check that
its printed output matches the corrected `// want` comments) and matched
the documented behavior against `starter/` (every test fails at
baseline; Exercise 1 alone, hand-completed in a scratch copy, turns only
`TestNewOrderCounter` green with everything else still failing
independently — no cross-exercise blocking). Safe to treat as live
without a re-verification pass.

## Lab 7 — done this session, notes for the pattern going forward

Lab 7 is shaped differently from every prior lab: no shared `internal/`
package at all. Each exercise is its own standalone `cmd/*` binary
(`dispatch`, `shiftend`, `courier`, `courierbuffered`, `radio`,
`workerpool`, `throttledpool`), plus a `standalone/` directory holding
two already-complete, deliberately-broken demo programs (`deadlock`,
`race`) that were never in scope for tests. That package-per-exercise
layout turned out to be a genuine advantage — verified stage-by-stage in
a scratch copy that completing one exercise's package never affects any
other package's test result, unlike Lab 5's shared-package coupling.
No lab this size has been this clean structurally.

- **New gotcha, specific to concurrency: a stuck producer/consumer can't
  fail a test normally, it can only hang one.** If a student's
  `dispatchOrders` or `worker` never sends and never closes its channel,
  the old naive test shape (`for v := range ch { ... }` then assert)
  blocks forever — `go test` eventually kills it with its own multi-minute
  default timeout and a goroutine dump, which is a terrible classroom
  experience compared to every other lab's instant, readable failure.
  **Fix, applied to all three affected tests** (`cmd/courier`,
  `cmd/courierbuffered`, `cmd/workerpool`): wrap the wait in a
  `select`/`time.After(2 * time.Second)` (a labeled `for`/`select`/`break`
  loop over the channel for the two producer tests; a `wg.Wait()` moved
  into its own goroutine signaling a `done` channel for the worker-pool
  test) so an unimplemented exercise fails in ~2 seconds with an explicit
  "did you forget to send/close/call wg.Done()?" message instead of
  hanging. This is worth calling out to Kevin explicitly: it's also a
  nice pedagogical callback, since it's the exact pattern Exercise 5
  teaches students to write by hand a little further into the same lab —
  said so directly in the lab's intro and in Exercise 5's Context.
  **Check this on every future concurrency-adjacent lab** — any
  pre-written test that waits on a channel or a `WaitGroup` fed by
  student code needs this guard, not just a bare receive/`Wait()`.
- **Found the same "unused import/declared-and-not-used blocks the whole
  test build" gotcha as Lab 2, but hitting three separate packages at
  once** (`cmd/courier`, `cmd/courierbuffered`, `cmd/workerpool`) —
  worse than previous labs because each was independently a real,
  pre-existing "doesn't build from an untouched starter checkout" bug,
  invisible without a real toolchain. `cmd/courier/main.go` and
  `cmd/courierbuffered/main.go` each imported `"fmt"` for a print
  statement that lives entirely inside the TODO'd portion of `main()` —
  removed the import, added a one-line TODO telling the student to add
  `"fmt"` back in as part of that task, same fix as the known gotcha.
  `cmd/workerpool/main.go`'s `main()` had `_ = jobs` and `_ = results`
  already present (silencing those two on purpose) but was missing the
  equivalent line for `wg` — added `_ = &wg` (a bare `_ = wg` fails
  `go vet` with "assignment copies lock value": `sync.WaitGroup`
  contains a `noCopy` sentinel specifically to catch that mistake, so the
  fix has to take the address, not the value). **`cmd/radio` has the same
  shape of failure** (`courierA`/`courierB` declared and unused until the
  `select` TODO is filled in) but was deliberately left as-is: Exercise 5
  has no pre-written test, so nothing blocks on it — it's simply an
  unfinished exercise correctly failing to build until finished, the same
  as any other incomplete TODO, not a "the lab is broken" bug. Only fix
  an unused-import/declared-and-unused failure when it's blocking a test
  from running; don't chase down every occurrence in a lab as if it were
  uniformly a defect.
- **The old lab's exercise numbering and its `cmd/*` source-file header
  comments had already drifted out of sync before this session touched
  anything.** `starter/cmd/shiftend/main.go` (Exercise 2, `defer`) had a
  visibly newer mtime than every other file in the lab — it was added
  later, after the numbering below it was already set, shifting every
  subsequent exercise's *true* number up by one in `exercise.md` without
  anyone propagating that shift into each file's own `// Exercise N:`
  header comment. Confirmed by literally reading the numbers: `exercise.md`
  called `cmd/courier` "Exercise 3" while the file's own header said
  "Exercise 2," and so on down the list — every file from `courier`
  onward was exactly one behind. Fixed all of them
  (`courier`, `courierbuffered`, `radio`, `standalone/deadlock`,
  `standalone/race` + its two `//go:build ignore` reference files,
  `cmd/workerpool/main.go`) to match `exercise.md`'s real numbering.
  Interestingly, `cmd/throttledpool/main.go`'s header already said
  "Exercise 9 (optional/stretch)" — the number it *should* be after
  deleting the trailing "prove it with a test" exercise (see next bullet)
  — strong evidence Kevin (or an earlier pass) had already started this
  exact renumbering and got partway through. **Lesson: when a lab's
  exercise count doesn't match cleanly, check every source file's own
  header comment against `exercise.md`, not just the file you're
  currently editing** — a silent drift like this is invisible unless you
  read every header side by side.
- Old "Exercise 9: Prove it with a test" deleted — unlike every prior
  lab, this one wasn't the *last* exercise (Exercise 10, the optional
  `throttledpool` stretch, came after it), so deleting it required
  renumbering Exercise 10 down to 9, not just removing a trailing
  section. Its test body was already fully written in `solution/`
  (`TestWorkerPoolAggregateResults`, sorted-results aggregate assertion —
  exactly the right shape per this doc's existing Lab 7 heads-up) and was
  reused verbatim as the pre-written `starter/` test, upgraded with the
  timeout guard above. Its final "deliberately break the aggregate logic,
  confirm the test catches it, fix it back" step was folded into the new
  Exercise 8's own task list as task 9, same as Lab 5 and Lab 6's
  precedent for a trailing test exercise's regression-check step.
- `checkOutCourier` (Exercise 2) was the only function needing an
  `io.Writer` refactor in this lab — it's the only exercise where the
  *content* of what gets printed is what the test needs to verify (that
  both deferred cleanup prints still fire on the early-return path).
  `dispatchOrders` (Exercises 3 and 4) needed no refactor: its test only
  cares about the values that travel over the channel, not what gets
  printed along the way, so `courierbuffered`'s "about to send order n"
  line stays on real stdout, unasserted, exactly as before.
- Verified with the real Go toolchain end to end: `gofmt -l .`, `go vet
  ./...`, `go build ./...`, `go test ./...`, and `go test -race ./...`
  all clean against `solution/`; `starter/` fails exactly as documented
  above (four packages red via clean assertions or timeouts, `cmd/radio`
  red via an expected build failure, nothing else affected); hand-completed
  each of Exercises 2, 3, 4, and 8 one at a time in a scratch copy and
  confirmed each package flips green independently with zero effect on
  any other package — the cleanest isolation of any lab converted so far.
  Also copied `main_mutex.go` and `main_channel.go` out of their
  `//go:build ignore` guard into standalone scratch modules and confirmed
  both compile on their own, matching what Exercise 7's instructions ask
  students to do.

All Lab 7 changes were pushed to
`/Users/kevincunningham/code/neueda/go-programming-reup/course/labs/07-concurrency/`
on Kevin's machine via the device bridge. Safe to treat as live without a
re-verification pass.
