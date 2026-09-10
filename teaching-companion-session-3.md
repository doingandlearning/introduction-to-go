# Topic 3 Teaching Companion — Flow Control & Data Structures

For today's session: 13:00–16:00, per `timings.md` (Session 3: 13:00–14:30,
break 14:30–14:50, Session 4: 14:50–16:00). Built from `slides.md`,
`exercise.md`, and the `code/` samples already in the repo.

---

## 1. Session plan — teach, then lab, twice

Yes, split it. The lab's six exercises already fall into two natural
halves that mirror the two lecture halves, so you don't need to invent a
new structure — just point the existing 90/break/70 split at the existing
exercise groups.

**Block A (control flow + structs) → Lab Exercises 1–3.**
**Block B (slices + maps, the hard half) → Lab Exercises 4–6.**

Given the room is less advanced than the prereqs assumed, bias extra time
toward the two hardest ideas — pointers/struct semantics and slice
aliasing — and let arrays stay brief, exactly as the slides already do
("in practice this matters less than it sounds").

### Session 3 — 13:00–14:30 (90 min)

| Time | Segment |
|---|---|
| 13:00–13:08 | Opening scenario hook (`backup := original`) — take guesses, **don't resolve it**, bank it for later |
| 13:08–13:16 | `for` — the one loop keyword, all four shapes |
| 13:16–13:24 | `if`/`else` with the scoped initializer + conditionless `switch` |
| 13:24–13:34 | `switch` fallthrough gotcha — live side-by-side demo, this is a "type in chat" moment worth pausing on |
| 13:34–13:44 | Pointers — `&`/`*`, no arithmetic. Go slower here than the slides imply; see §2 for scaffolding |
| 13:44–13:56 | Structs + value semantics — the `reset` vs `resetPtr` demo, live |
| 13:56–14:01 | Arrays — brief. Fixed size baked into the type, value-copied, then move on |
| 14:01–14:03 | Frame the lab: "Exercises 1–3, tests are already written and failing — make them pass" |
| 14:03–14:28 | **Lab block 1** — Exercises 1–3 (loops, switches, struct value/pointer). Circulate; this is where the pointer confusion from 13:34 will resurface concretely |
| 14:28–14:30 | Pulse check — who's green on all three tests? |

### Break — 14:30–14:50 (20 min)

### Session 4 — 14:50–16:00 (70 min)

| Time | Segment |
|---|---|
| 14:50–14:56 | Slices — "what you actually use," quick literal + `append` + `len`/`cap` |
| 14:56–15:10 | **The core concept**: slicing shares memory. Viewfinder/film metaphor, then live `go run ./cmd/slicealias`, reading `len`/`cap` out loud at every step before revealing output. Don't rush this — it's the one idea the whole afternoon hinges on |
| 15:10–15:14 | Range — one keyword, five collection types |
| 15:14–15:22 | Maps — comma-ok, then the nil-map-panics-on-write gotcha, then map iteration order (`cmd/wordcount` demo if time allows) |
| 15:22–15:24 | Frame the lab: "Exercises 4–6 are observational — no pre-written tests, you're reading output, not chasing green checkmarks" |
| 15:24–15:52 | **Lab block 2** — Exercises 4–6 (reproduce slice aliasing, comma-ok word counter, map order). Exercise 4 is the one to hover near — it's a direct rerun of the concept from 14:56 |
| 15:52–15:57 | Close the loop: back to `backup := original` from 13:00 — now they can answer it themselves |
| 15:57–16:00 | Summary + bridge to Topic 4 |

**Buffer notes, given the room's level:**

- Both lab blocks are sized tight (25 and 28 min for six exercises total). If Block 1 overruns, it's fine to cut Arrays lecture time to zero — the slides already treat it as optional ("idiomatic Go rarely uses raw arrays directly") and nothing in the lab exercises depends on lecturing arrays separately from slices.
- Exercise 6 (map iteration order) is the lightest of the six — pure observation, no code to write beyond a `range` loop reusing Exercise 5's map. If Block 2 is short on time, it's the one to let people finish at their own pace after 16:00 or fold into tomorrow's warm-up, rather than the one to rush.
- If the room is visibly stuck on pointers at 13:44, it's worth trading 5 minutes from Arrays to extend there — pointer confusion compounds into struct-semantics confusion immediately afterward, so it's cheaper to fix at the source.

---

## 2. Supplementary practice — for bolstering weaker foundations

The slides move fast on the assumption of "familiar with pointers/OOP
from C or Java." Use these as extra worked examples or a handout for
people who need a slower run-up before or during the lab blocks. Each is
a **predict-the-output** question — get them to write down what they
think happens before running it. Answers are at the bottom of each
group; hints are inline.

### 2a. Pointers — before structs

**Q1.** What does this print?

```go
x := 5
p := &x
fmt.Println(*p)
```

*Hint: `&x` is "the address of x." `*p` means "go to that address and
read what's there." Read `*p` as "the thing p points at," not as a
multiplication.*

**Q2.** What does this print, and why is it different from Q1?

```go
x := 5
p := &x
*p = 10
fmt.Println(x)
```

*Hint: writing through `*p` writes to the same memory `x` lives in.
There's only one `x` — `p` doesn't have its own copy.*

**Q3.** Will this compile?

```go
x := 5
p := &x
p2 := p + 1
```

*Hint: no. This is the line to point at when someone asks "so it's just
like a C pointer?" — Go pointers can't be walked like C pointers can.
There is no `p + 1`.*

**Q4.** What prints?

```go
func double(n int) {
    n = n * 2
}

x := 5
double(x)
fmt.Println(x)
```

*Hint: `n` is a parameter — a fresh copy of `x`'s value, local to
`double`. Nothing about calling a function changes the caller's variable
unless you pass a pointer. This is the exact same rule Q5 tests with a
struct instead of an int.*

**Q5.** What prints?

```go
func doublePtr(n *int) {
    *n = *n * 2
}

x := 5
doublePtr(&x)
fmt.Println(x)
```

*Hint: contrast directly against Q4. Same function shape, but this one
was handed the address, so it can reach back and change the original.*

Answers: Q1 → `5`. Q2 → `10`. Q3 → does not compile (`invalid operation:
p + 1 (mismatched types *int and untyped int)`). Q4 → `5` (unchanged).
Q5 → `10` (changed).

### 2b. Structs — same rule, different-looking code

**Q6.** Predict the output:

```go
type Point struct{ X, Y int }

func moveRight(p Point) {
    p.X = p.X + 1
}

pt := Point{X: 0, Y: 0}
moveRight(pt)
fmt.Println(pt.X)
```

*Hint: this is Q4 again, wearing a struct instead of an int. `p` inside
`moveRight` is a whole separate copy of the struct — every field
included.*

**Q7.** Now this one:

```go
func moveRightPtr(p *Point) {
    p.X = p.X + 1
}

pt := Point{X: 0, Y: 0}
moveRightPtr(&pt)
fmt.Println(pt.X)
```

*Hint: Q5's pattern again. The common trip-up: people expect this
because "Java objects are references" — but Go structs aren't objects in
that sense, they're values, exactly like an `int`, until you explicitly
take a pointer.*

Answers: Q6 → `0`. Q7 → `1`.

### 2c. Arrays vs. slices — the type signature gives it away

**Q8.** Are `a` and `b` the same type? Will this compile?

```go
var a [4]int
var b [5]int
a = b
```

*Hint: no — the size is part of the type. `[4]int` and `[5]int` are as
different as `int` and `string` are. This is the one thing to remember
about arrays even if nothing else sticks: the length is compile-time,
baked in.*

**Q9.** And this?

```go
var a []int
var b []int
a = b
```

*Hint: compiles fine — no size in the brackets means it's a slice, and
slices of the same element type are always assignment-compatible
regardless of how many elements are currently in them.*

### 2d. Slice aliasing — the one to spend the most time on

Work through these **in order**, predicting `len`/`cap`/contents before
each run. This is a direct paper-and-pencil version of `cmd/slicealias`.

**Q10.** Start here:

```go
original := []int{10, 20, 30, 40, 50}
view := original[1:3]
fmt.Println(view)        // ?
fmt.Println(len(view), cap(view))  // ?
```

*Hint: `view` holds indices 1 and 2 of `original` — values `20, 30`.
`len` is the count of elements in the slice (2). `cap` is how far `view`
could grow before running out of room in the *underlying array* — from
index 1 to the end of `original`, so 4.*

**Q11.** Now:

```go
view[0] = 999
fmt.Println(original)   // ?
```

*Hint: `view[0]` is the same memory as `original[1]`. There is one
array; `view` and `original` are two different windows onto it.*

**Q12.** Now:

```go
view = append(view, 100)
fmt.Println(original)   // ?
fmt.Println(len(view), cap(view))  // ?
```

*Hint: `cap(view)` was 4 and `len(view)` was 2 before this — there was
room for one more element without reallocating. So `append` wrote
straight into the next slot of the shared array, which happens to be
`original[3]`. Nobody wrote `original[3] = 100` — it happened as a side
effect.*

**Q13.** Now push it past capacity:

```go
view = append(view, 200, 300, 400)
fmt.Println(len(view), cap(view))  // ? — compare cap to before
view[0] = -1
fmt.Println(original)  // ?
```

*Hint: this append needed more room than was left, so Go allocated a
brand-new, bigger array, copied `view`'s current contents into it, and
pointed `view` at the new array. `original` never knew any of this
happened — `view` now lives somewhere else entirely, so mutating `view`
no longer touches `original`.*

**The one sentence to make them repeat back to you:** *"Re-slicing
copies the header — pointer, length, capacity — never the data
underneath it. Sharing lasts exactly until an `append` runs out of
capacity."*

Answers: Q10 → `[20 30]`, `2 4`. Q11 → `[10 999 30 40 50]`. Q12 →
`[10 999 30 100 50]`, `3 4`. Q13 → cap jumps to something ≥7 (implementation-
dependent, e.g. 8); after the reallocation, `original` is unchanged by the
`view[0] = -1` line.

### 2e. Maps — comma-ok

**Q14.** What prints for each line?

```go
m := map[string]int{"a": 1, "b": 0}
va, oka := m["a"]
vb, okb := m["b"]
vc, okc := m["c"]
fmt.Println(va, oka)  // ?
fmt.Println(vb, okb)  // ?
fmt.Println(vc, okc)  // ?
```

*Hint: `"b"` is genuinely in the map with value `0`. `"c"` was never
inserted at all. Both `vb` and `vc` come back as `0` — the only way to
tell them apart is `okb` (`true`) vs `okc` (`false`). If you only ever
look at the value, these two cases are indistinguishable.*

Answers: `1 true`, `0 true`, `0 false`.

### 2f. Switch fallthrough

**Q15.** What does this print?

```go
n := 2
switch n {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two")
    fallthrough
case 3:
    fmt.Println("three")
case 4:
    fmt.Println("four")
}
```

*Hint: matches `case 2`, prints "two", then `fallthrough` forces
execution into `case 3`'s body regardless of whether `n == 3` — it does
not re-check the condition, it just runs the next case's code. Execution
stops there because `case 3` has no `fallthrough` of its own.*

Answer: prints `two` then `three`, and stops — `four` never prints.

---

## 3. Anticipated questions — experienced C / Python developers

### From the C side

**"So a Go pointer is just like a C pointer?"**
Same two operators (`&`, `*`), same idea of "an address," but Go removes
pointer arithmetic entirely — no `p + 1`, no pointer-to-array decay, no
`void*` you cast freely between types. A Go pointer always points at
exactly one thing of exactly one type. You get the "can mutate the
original" power without the "can walk off the end of an array" danger.

**"Is a Go slice basically a C fat pointer — pointer + length?"**
Close, but it's pointer + length + *capacity* (three fields, not two).
The capacity is what makes `append` able to grow in place sometimes and
reallocate other times — that's the mechanic behind the slide on slice
aliasing.

**"Where does a slice's backing array live — stack or heap?"**
Determined by escape analysis, same as any Go value — if the compiler
can prove nothing outlives the function, it can stay on the stack;
otherwise it escapes to the heap. Don't over-promise a rule here; the
honest answer is "the compiler decides, and you can check with `go build
-gcflags="-m"` if you actually need to know."

**"Are structs padded/aligned like in C? Can I control layout?"**
Yes, Go structs are padded for alignment the same way C structs are, and
field order affects size (classic "reorder fields to shrink the struct"
trick applies). There's no `#pragma pack` equivalent in the language
itself, though the `unsafe` package exposes low-level layout if you
truly need it — well outside what this course covers.

**"Passing a struct by value — isn't that expensive for big structs,
like in C?"**
Yes, same cost model as C: a value pass copies every field. For large
structs passed frequently, passing a pointer avoids the copy — but that's
a performance decision, separate from the mutation decision. You can
pass a pointer purely for cheapness even if you never intend to mutate
through it (methods with pointer receivers, covered in Topic 4, usually
settle this by convention).

**"No `union`, no `void*` — how do I do type punning or hold 'any type'
generically?"**
`interface{}` (aliased as `any`) is the "can hold anything" type, checked
at runtime via type assertions or type switches — you'll see this
properly in Topics 5–6. It's not memory-layout type punning like a C
union; it's a completely different mechanism (a Go interface value is
itself a type-tag-plus-pointer pair under the hood).

### From the Python side

**"`original[1:3]` in Python gives me a new list. Why doesn't Go's
slicing copy?"**
This is the single biggest transfer-negative habit in the room — flag it
explicitly. Python slicing on a list always allocates a new list. Go
slicing on a slice never allocates — it hands back a new *view* over the
same backing array. It's less like Python slicing and more like a
NumPy view (`arr[1:3]` on a NumPy array *also* aliases, for the same
reason) — if anyone in the room has used NumPy, that's the closer
analogy.

**"Are Go maps like Python dicts? I thought dicts became ordered in
3.7+."**
Same hash-table idea, but Go deliberately did the opposite of what
Python did — Python 3.7 made insertion order a language guarantee; Go
randomizes iteration order on purpose, specifically so nobody's code can
accidentally depend on an order the map never promised. If you need
order in Go, you keep a separate slice of keys yourself.

**"Why do I need `make(map[string]int)` sometimes but slices seem to
work fine as `var s []int`?"**
A nil slice is safe to read and safe to `append` to (append allocates on
first use). A nil map is safe to *read* but panics the instant you
*write* to it. It's an intentional asymmetry, not an oversight — flag it
as one of the most common first-week panics.

**"No list comprehensions? How do I filter/map a slice?"**
There isn't comprehension syntax — you write the loop out with `for`.
This feels like a step backward for a week and then stops mattering;
Topic 6 (Functional Programming) covers the idiomatic higher-order-
function patterns Go does have.

**"Is a Go struct basically a Python dataclass?"**
Similar spirit (a typed bag of named fields, cheap to construct via a
literal) but static, not dynamic — you can't add a field at runtime, and
there's no `__init__` unless you write a constructor function yourself.
Topic 4 covers methods and constructor conventions.

**"`KeyError` vs. comma-ok — why not just raise like Python does?"**
Go doesn't have exceptions at all (that's deliberate, and comes up again
properly around error handling). A missing map key returning a silent
zero value is consistent with Go's general "no exceptions, check
explicitly" philosophy — comma-ok is the map-specific version of the
`err != nil` pattern from Topic 2.

**"Does `append` always return a new slice? Do I have to reassign it?"**
Always reassign the result (`s = append(s, x)`), every time, even though
it sometimes mutates the existing backing array in place. The reason
you must reassign regardless: you can't tell from the call site alone
whether this particular `append` reallocated or not, and if it did, the
old slice variable is now stale. Reassigning is the one habit that's
correct in both cases.

### General / either background

**"Why does Go bake the array length into the type instead of letting
me query it at runtime, like a C array via `sizeof` or Python's
`len()`?"**
It buys compile-time safety — `[4]int` and `[5]int` genuinely can't be
mixed up by accident — at the cost of flexibility, which is exactly why
slices exist as the practical, resizable alternative built on top. The
slides make this trade explicit: arrays are the foundation, slices are
what you reach for.

---

## 4. "Why does this matter?" — the deep answer

This is worth landing early, ideally right after the opening scenario
hook, because it reframes the whole afternoon from "here's Go trivia"
into "here's the thing that actually bites people in production."

**The short version:** everything in Topic 3 is really one idea wearing
several outfits — *Go makes copying explicit and sharing implicit,
except for the handful of types (slices, maps, pointers, channels) where
it's the reverse* — and almost every real bug new Go developers write in
their first few months traces back to guessing wrong about which side of
that line a piece of code is on.

**Why Go is built this way at all.** Go was designed by people who'd
watched large C++ codebases at Google rot under two different failure
modes: unpredictable aliasing (nothing stops you handing out a pointer
and forgetting who else holds one) and unpredictable cost (implicit
copies, operator overloading, hidden allocations — code that *looks*
cheap but isn't). Go's answer is to make the copy-vs-share question
*visible in the type itself* rather than something you have to trace
through call chains to know. An `int`, a `struct`, an array — plain
value, always copied, always safe to hand around without a second
thought. A slice, a map, a pointer — carries a reference to shared state,
always, no exceptions, no way to accidentally get a deep copy you didn't
ask for. Once someone internalizes that one axis, most of Topic 3's
"surprises" stop being surprises.

**Why the slice-aliasing bug specifically matters so much.** It's in the
slides as "the single most common new-Go-dev bug" for a concrete reason:
it's *intermittent*. The exact same line of code — `view = append(view,
x)` — either quietly corrupts a caller's data or safely reallocates away
from it, and which one happens depends entirely on capacity at that
moment, which depends on slice history the reader often can't see from
the line in front of them. That's the profile of the worst kind of bug:
it passes code review, it passes casual testing, and it shows up three
weeks later when a slice that always happened to have zero spare
capacity in dev suddenly has some in production. Understanding the
mechanism (header vs. backing array, capacity as the trigger) turns that
from "mysterious flaky bug" into "oh, I know exactly what to check."

**Why value-vs-pointer semantics on structs isn't just a style choice.**
Once you get to Topic 8's handler/service/repository pattern and
dependency injection, "does this struct hold shared, mutable state, or
is it a plain data value?" becomes an architectural decision, not a
syntax one. A `Repository` you inject probably needs to be a pointer
(one shared instance, holding a DB connection everyone should reuse). A
`Book` record flowing through your code probably wants value semantics
(each function gets its own copy, nobody's surprised by action at a
distance). Getting this instinct right now, on small toy examples, is
what makes those later architectural decisions feel obvious instead of
arbitrary.

**Why it matters even more once concurrency shows up (Topic 7).** A
slice or map shared between goroutines is exactly the kind of aliased,
mutable state that makes concurrent code dangerous — two goroutines
`append`-ing to views of the same backing array, or writing to the same
map without synchronization, is a data race, and Go's race detector will
catch some of it but not all of it if you don't understand what's
aliased in the first place. Today's mental model — "who else has a
window onto this same memory?" — is the exact question you'll need to
ask again, under time pressure, once goroutines are involved.

**The honest pitch to give the room, if asked "why not just make slices
copy like Python?"** Performance and predictability at scale. Copying a
slice's *data* on every assignment or function call would make Go slices
behave like Python lists — safe by default, but silently expensive the
bigger they get, and expensive in a way that's invisible at the call
site. Go chose the opposite trade: cheap, constant-time copies of the
*header* always, with the sharp edge of aliasing as the cost of that
cheapness — and then gives you `copy()` and `append([]T{}, src...)` as
explicit, visible escape hatches for the (much rarer) times you actually
want an independent copy. Explicit is slower to write and faster to
run — and once you know the rule, it's also easier to reason about,
because nothing is copying data behind your back that you didn't ask
for.
