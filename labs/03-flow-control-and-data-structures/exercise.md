# Lab 3: Flow Control & Data Structures — a library circulation desk

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

Both directories share the same shape: a `cmd/librarian` package (the
program entry point) and an `internal/catalog` package (the helper code
you'll implement). The scenario throughout is a small library circulation
desk — a slice of books, a waitlist, and a log of checkouts — which gives
every language feature in this topic something concrete to operate on.

---

## Exercise 1: Three shapes of `for`

**Objective:** Use all three non-range shapes of Go's single loop keyword.

**Context:** `starter/internal/catalog/loops.go` has three `TODO`
functions operating on a `[]Book` (a `Book` is `{Title string; Copies
int}`, defined in `catalog.go`).

**Tasks:**

1. Implement `PrintCatalog(books []Book)` with a **classic C-style** `for`
   loop (`for i := 0; i < len(books); i++`) that prints each book's title
   and copy count.
2. Implement `CountAvailable(books []Book) int` with a **while-style**
   `for` loop (a bare condition, no init or post statement) that counts
   how many books have `Copies > 0`.
3. Implement `FindFirstAvailable(books []Book) (Book, bool)` with an
   **infinite** `for {}` loop that scans until it finds a book with
   `Copies > 0`, then `break`s. Return the zero `Book` and `false` if the
   loop finishes without finding one.
4. Run `go run ./cmd/librarian` from `starter/` and confirm the catalog
   prints, the available count is right, and the first available book is
   correct.

**Key Learning:** Go has one looping keyword. The classic, while-style,
and infinite forms aren't different keywords with different rules —
they're the same `for` with different pieces of the header omitted.

---

## Exercise 2: Switches — grading fees and desk hours

**Objective:** Use a conditionless `switch` and deliberate `fallthrough`.

**Context:** `starter/internal/catalog/decisions.go` has two `TODO`
functions.

**Tasks:**

1. Implement `LateFeeTier(daysLate int) string` using a `switch` with
   **no condition** (`switch { case daysLate >= 30: ... }`). Return
   `"suspended"` for 30+ days late, `"warning"` for 7-29 days, `"none"`
   for under 7.
2. Implement `DeskSchedule(day int) string` using `switch day { ... }`
   where `day` is 1 (Monday) through 7 (Sunday). Group Monday-Friday
   together as `"open 9am-6pm"`. Give Saturday its own case that ends
   with a deliberate `fallthrough` into Sunday's case, so both return
   `"open 10am-2pm only"`.
3. Run the program. Confirm day 6 and day 7 produce the same string via
   the fallthrough, then comment out the `fallthrough` line and confirm
   day 6 now returns something else (or nothing useful) instead.

**Key Learning:** Go's `switch` doesn't fall through by default — the
opposite of C, Java, and JavaScript. Every case is an implicit break;
`fallthrough` is how you opt back into the old behavior, one case at a
time, explicitly.

---

## Exercise 3: Struct value vs. pointer semantics

**Objective:** Prove to yourself when a function can and can't mutate the
struct you passed it.

**Context:** `starter/internal/catalog/structs.go` has two `TODO`
functions, both operating on `Book`.

**Tasks:**

1. Implement `ResetCopies(b Book) Book` — takes a `Book` **by value**,
   sets `Copies` to `0` on its local copy, and returns it.
2. Implement `ResetCopiesPtr(b *Book)` — takes a `*Book`, and sets
   `Copies` to `0` **through the pointer**, with no return value needed.
3. In `cmd/librarian/main.go`, call `ResetCopies(book)` on a book from
   your catalog and print `book.Copies` afterward — confirm it's
   unchanged.
4. Call `ResetCopiesPtr(&book)` on the same book and print `book.Copies`
   again — confirm it's now `0`.

**Key Learning:** Passing a struct to a function passes a copy, exactly
like an `int` — the function can only mutate the original if you pass a
pointer explicitly. This is the reverse of what "pass an object" trains
you to expect in Java, where object references make field mutation look
automatic.

---

## Exercise 4: Reproduce the slice-aliasing bug

**Objective:** Trigger the single most common real Go bug on purpose, then
find the exact line that stops it from happening.

**Context:** `cmd/librarian/main.go` has a `TODO` block under `// --
Exercise 4 --`. You'll build a waitlist scenario from scratch here, not
in the `catalog` package — this one's about observing behavior, not
returning a value.

**Tasks:**

1. Create `waitlist := []string{"Alice", "Bob", "Carla", "Dev", "Eve"}`.
2. Take `nextUp := waitlist[1:3]` (a sub-slice). Print `len`/`cap` of
   both `waitlist` and `nextUp`.
3. Set `nextUp[0] = "REPLACED"`. Print `waitlist` — confirm it changed,
   even though you never assigned to `waitlist` directly.
4. Append one name to `nextUp` (something that still fits within its
   capacity). Print `waitlist` again — confirm a name in `waitlist` got
   silently overwritten.
5. Append enough more names to `nextUp` to exceed its capacity. Print
   `cap(nextUp)` before and after this append and confirm it jumped.
6. Mutate `nextUp[0]` one more time and print `waitlist` — confirm it's
   now unaffected, because the reallocation in step 5 broke the sharing.

**Key Learning:** Re-slicing copies a small header (pointer, length,
capacity) — never the underlying data. Two slices with overlapping
ranges share memory until an `append` forces one of them past its
capacity, at which point Go silently allocates a new backing array and
the sharing ends. Whether a given `append` mutates shared data or not
depends entirely on capacity at that moment — which is exactly why this
bug is intermittent in real code.

---

## Exercise 5: Word-frequency counter with comma-ok

**Objective:** Build a `map[string]int` and use comma-ok to distinguish
"counted, zero" from "never seen."

**Context:** Another `TODO` block in `cmd/librarian/main.go`, under `//
-- Exercise 5 --`.

**Tasks:**

1. Given a checkout log string, e.g. `"pratchett tolkien pratchett le-guin
   pratchett tolkien"`, split it with `strings.Fields` and build a
   `map[string]int` counting occurrences of each author.
2. For three lookups — one author who appears several times, one who
   never appears in the log at all, and (if you want a sharper edge
   case) one you explicitly set to `0` in the map beforehand — use the
   comma-ok idiom (`n, ok := counts[name]`) to print either `"<name>:
   <n> checkout(s)"` or `"<name>: never checked out"`.
3. Confirm the explicitly-zeroed author and the never-seen author print
   different messages, even though both report `n == 0`.

**Key Learning:** A missing map key and a key whose value happens to be
the zero value look identical if you only check the value. `ok` is the
only reliable way to tell "never written" from "written, and it's zero."

---

## Exercise 6: Map iteration order

**Objective:** Confirm for yourself that Go's map order isn't just
"whatever order you don't understand yet" — it's genuinely randomized.

**Context:** A final `TODO` block in `cmd/librarian/main.go`, under `//
-- Exercise 6 --`.

**Tasks:**

1. Using the map from Exercise 5, `range` over it and print the keys in
   order, twice in a row, in the same run.
2. Run the whole program two or three separate times (`go run
   ./cmd/librarian` again from the shell) and compare the orders you
   see across runs.

**Key Learning:** Go deliberately randomizes map iteration order so that
no code can accidentally depend on an order that was never promised —
unlike Python 3.7+ dicts or Java's `LinkedHashMap`, which are
insertion-ordered. If your code needs a stable order, track it
separately (e.g. a slice of keys) — the map itself won't give you one.

---

## Exercise 7: Prove it with a test

**Objective:** Extend the testing habit from Topic 2 to this lab's real
code — a conditionless `switch` with `fallthrough`, and value vs. pointer
struct semantics.

**Context:** `starter/internal/catalog/` has no test files yet. You
implemented `LateFeeTier` and `DeskSchedule` in Exercise 2, and
`ResetCopies`/`ResetCopiesPtr` in Exercise 3 — now write the tests the
same way Topic 2's lecture and Exercise 7 did for `divide`: a `_test.go`
file, `TestX(t *testing.T)` functions, no framework.

**Tasks:**

1. Create `starter/internal/catalog/decisions_test.go`, `package
   catalog`, importing `testing`. Write `TestLateFeeTier` covering at
   least the boundary days — one just under 7, one at exactly 7, one just
   under 30, one at exactly 30 — asserting the exact tier string each
   time.
2. In the same file, write `TestDeskSchedule` covering a weekday, and
   both Saturday (`6`) and Sunday (`7`) — assert Saturday and Sunday
   return the *same* string, which is only true because of the
   `fallthrough` you wrote in Exercise 2.
3. Create `starter/internal/catalog/structs_test.go`, same package and
   import. Write `TestResetCopies` that calls `ResetCopies` on a `Book`
   with `Copies` set to a nonzero number, and asserts **both** that the
   returned copy has `Copies == 0` **and** that the original `Book` you
   passed in is unchanged.
4. In the same file, write `TestResetCopiesPtr` that calls
   `ResetCopiesPtr` on a `&Book` and asserts the original now has
   `Copies == 0`.
5. Run `go test ./...` from `starter/`. All four tests should pass once
   Exercises 2 and 3 are correctly implemented.
6. Deliberately break the fallthrough — comment out the `fallthrough`
   line in `DeskSchedule` — and re-run `go test ./...`. Confirm
   `TestDeskSchedule` fails and tells you *why* (Saturday and Sunday no
   longer match). Put the `fallthrough` back and confirm the test passes
   again.

**Key Learning:** In Exercise 2 you confirmed the fallthrough behavior by
eyeballing printed output once. A test turns that same behavior into an
explicit, checked contract — "Saturday and Sunday must always produce the
same string" is no longer something you have to remember to re-check by
hand every time this file changes; `go test` checks it for you, every
time, forever.

---

## Summary

By the end of this lab you should be able to:

- Write all three non-range shapes of `for`, and know when to reach for
  each one
- Predict whether a `switch` case falls through, and use `fallthrough`
  deliberately when you want the old C/Java/JS default back
- Explain why passing a struct by value doesn't let a function mutate the
  caller's copy, and fix that by passing a pointer
- Reproduce the slice-aliasing bug on demand, and locate the exact
  `append` call that breaks sharing by watching `cap()` change
- Use the comma-ok idiom to distinguish a missing map key from a
  present key holding a zero value
- Explain why Go map iteration order is randomized, and what to do if
  your code actually needs a stable order
- Write `TestX(t *testing.T)` functions against this lab's own `switch`
  and pointer-mutation code, and watch a deliberately broken
  `fallthrough` get caught by the test instead of by eyeballing output
