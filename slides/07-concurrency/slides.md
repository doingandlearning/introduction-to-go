---
title: "**Concurrency**"
sub_title: "Go Programming — Topic 7"
author: Kevin Cunningham
---

## Opening scenario

You need to fetch data from 10 independent sources — APIs, databases,
whatever — as fast as possible. They don't depend on each other.

In Java, before you write a line of fetching code, you'd think hard about
a thread pool: how many threads, what size, what happens if you get the
number wrong.

**Type in chat: what do you think Go makes you think about instead?**

We'll come back to this once you've seen goroutines and channels in
action.

<!--
speaker_note: |
  Let a few guesses land - "nothing," "less," "the data flow" are all
  reasonable. Don't confirm or deny yet. The honest answer is "almost
  nothing about threads, quite a lot about how results get back to you"
  - that's the payoff slide later.
-->

<!-- end_slide -->

## Why concurrency gets extra time this week

Every other topic in this course has an equivalent you already know:
interfaces, generics, error handling all map onto something from Java,
Python, or TypeScript.

<!-- pause -->

**Concurrency doesn't.** It's the single feature Go was designed around
from day one, and it's the reason companies pick Go for network services,
CLIs, and infrastructure tools in the first place.

We're going to spend real time here — more than the topic slot implies —
because this is what makes Go, Go.

<!-- end_slide -->

<!-- jump_to_middle -->

Goroutines
===

<!-- end_slide -->

## Launching a goroutine

A goroutine is a function running concurrently with the rest of your
program. You launch one with the `go` keyword.

```go
func sayHi() {
    fmt.Println("hi")
}

func main() {
    go sayHi()          // runs concurrently
    fmt.Println("main") // main doesn't wait for sayHi
}
```

<!-- pause -->

**Demo:** run this. `"main"` and `"hi"` can print in either order, or
`"hi"` might not print at all if `main` exits first — nothing waits for
a bare `go` call by default.

<!--
speaker_note: |
  Run it a few times in a row live if you can - the nondeterminism is
  the whole lesson. Some runs print both lines, some only print "main".
  That surprise is exactly what WaitGroup fixes two slides from now.
-->

<!-- end_slide -->

## Not an OS thread. Not quite a "green thread" either.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Java `Thread`**

- 1:1 with an OS thread
- ~1MB stack, fixed
- OS-scheduled
- Expensive enough that pool sizing is a real design decision

<!-- column: 1 -->

**Go goroutine**

- Scheduled by the Go runtime, M:N onto far fewer OS threads
- Starts at a couple of KB, grows as needed
- Practically free to create

<!-- reset_layout -->

<!-- pause -->

**"Launch ten thousand of them" is a normal sentence in Go.** Where you'd
hesitate before spinning up a Java thread pool, `go doSomething()` inside
a loop of thousands is unremarkable.

<!--
speaker_note: |
  GOMAXPROCS controls how many OS threads the runtime uses to execute
  goroutines - worth naming if someone asks "so how many threads
  actually run this," but don't get pulled into scheduler internals
  unless there's time and appetite for it.
-->

<!-- end_slide -->

## The sticky note, not the new hire

Hiring a full-time employee with their own office is slow and expensive —
that's an OS thread. A goroutine is a sticky note handed to whoever's
free at a shared, open-plan desk.

<!-- incremental_lists: true -->

- Writing a sticky note (launching a goroutine) costs almost nothing
- The office manager (the Go runtime scheduler) decides who picks it up and when
- You can hand out thousands of sticky notes without hiring anyone new

<!-- incremental_lists: false -->

<!-- end_slide -->

## No async/await

JavaScript, Python's `asyncio`, and TypeScript all mark concurrency in
the function signature: `async function`, `await somePromise()`. You can
see it in the type.

<!-- pause -->

Go goroutines look like **ordinary synchronous calls**. The only marker
that something runs concurrently is the `go` keyword at the call site.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```js
// JS: the signature tells you
async function fetchAll() {
  const a = await fetchA();
  const b = await fetchB();
}
```

<!-- column: 1 -->

```go
// Go: nothing in fetchAll's
// signature says "concurrent"
func fetchAll() {
    go fetchA()
    go fetchB()
}
```

<!-- reset_layout -->

**No "async all the way up" tax** — you never have to make a caller
`async` just because it called something `async`. But you lose the
visual breadcrumb. Train your eye to spot `go` statements and channel
operations instead of a keyword in every signature.

<!-- end_slide -->

## Waiting for goroutines: `sync.WaitGroup`

`go sayHi()` doesn't wait. If you need to know when a batch of goroutines
has finished, use a `WaitGroup`.

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println("worker", n)
    }(i)
}

wg.Wait() // blocks until all 5 have called Done
```

<!-- pause -->

`Add(1)` before launching, `Done()` when a goroutine finishes (almost
always via `defer`), `Wait()` blocks the caller until the count hits
zero.

**Demo:** run this several times — the worker numbers print in a
different order almost every run. `WaitGroup` guarantees *completion*,
never *order*.

<!--
speaker_note: |
  Have delegates run it 4-5 times back to back and shout out the order
  they see. It's a cheap, memorable way to make "no ordering guarantee"
  land as something they watched rather than something you told them.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Channels
===

<!-- end_slide -->

## Channels: typed conduits between goroutines

A channel is a typed pipe. One goroutine sends a value in, another
receives it out.

```go
ch := make(chan int)  // unbuffered channel of int

ch <- 42     // send: blocks until someone receives
v := <-ch    // receive: blocks until someone sends
```

<!-- pause -->

An **unbuffered** channel is a synchronous rendezvous: the send doesn't
complete until a receive is ready for it, at the same moment. Neither
side proceeds alone.

<!-- end_slide -->

## The relay baton

A channel is the baton at a track relay — the *only* legitimate way data
moves from one runner to the next. Nobody shouts the result across the
field (that's shared memory without synchronization, and it's how you
get race conditions).

<!-- pause -->

**Unbuffered = a strict handoff.** Both runners have to be at the
exchange zone at the same moment, or the handoff doesn't happen. One
runner arriving early just waits.

<!-- end_slide -->

## Share memory by communicating

This is an actual Go proverb, and it's the philosophy behind why channels
exist:

> Don't communicate by sharing memory; share memory by communicating.

<!-- pause -->

Java's instinct is lock-first: shared variable, `synchronized`, mutex.
Go's idiomatic default reach for coordinating goroutines is a **channel**
— pass the data to whoever needs it next, rather than letting many
goroutines poke the same variable under a lock.

<!-- pause -->

**Be honest about this: it's a default, not a law.** For small, simple
shared state — a counter, a cache entry — a `sync.Mutex` is still
completely correct and often simpler than wiring up a channel. Don't let
"share memory by communicating" calcify into "channels for everything."
We'll see both approaches to the same problem shortly.

<!-- end_slide -->

## Producer / consumer over an unbuffered channel

```go
func produce(ch chan<- int) {
    for i := 1; i <= 5; i++ {
        ch <- i // blocks until main receives
    }
    close(ch) // signal: no more values coming
}

func main() {
    ch := make(chan int)
    go produce(ch)

    for v := range ch { // exits automatically on close
        fmt.Println("got", v)
    }
}
```

<!-- pause -->

`range` over a channel receives values until the channel is closed, then
exits the loop cleanly — no manual "am I done yet" check needed.

<!-- end_slide -->

## Closing a channel: a signal, not cleanup

Closing a channel means **"no more values are coming,"** nothing more.
It is not garbage collection, not "free the resource" — a closed channel
is still a channel, and previously-sent values can still be drained.

<!-- incremental_lists: true -->

- Close from the **sending** side. Never the receiving side.
- Never close a channel twice — both mistakes **panic** at runtime.
- Detect closure on receive with comma-ok: `v, ok := <-ch` — once the channel is drained and closed, `ok` is `false` and `v` is the zero value.
- `range` over a channel does the comma-ok check for you and exits automatically.

<!-- incremental_lists: false -->

**Demo:** close a channel twice on purpose and read the panic message —
it names exactly what you did.

<!-- end_slide -->

## Goroutines don't get garbage collected

A goroutine blocked forever on a channel nobody will ever write to (or
read from) doesn't disappear just because nothing else references it. It
sits there, holding its small stack, for the **lifetime of the process**.

<!-- pause -->

This is a **goroutine leak**. Unlike a heap leak, it doesn't show up as
steadily climbing memory in the same obvious way — a few thousand
leaked, blocked goroutines are cheap enough in memory that this can run
unnoticed for a long time before it's a real problem.

**The fix is almost always the same shape:** make sure every goroutine
you launch has a guaranteed way to finish — a closed channel, a
cancellable context, a value that's certain to arrive.

<!-- end_slide -->

<!-- jump_to_middle -->

Buffered channels and `select`
===

<!-- end_slide -->

## Buffered channels: slack in the handoff

`make(chan int, 3)` creates a channel with room for 3 values in flight
before a send blocks.

```go
buf := make(chan int, 3)
buf <- 1 // doesn't block — room in the buffer
buf <- 2 // doesn't block
buf <- 3 // doesn't block — buffer is now full

// buf <- 4 would block here until something receives
```

<!-- pause -->

**The bucket on the conveyor belt.** With an unbuffered channel, the
sender and receiver must meet at exactly the same moment. A buffered
channel is a bucket between them — the sender can drop values in and
keep moving, right up until the bucket's full.

<!-- end_slide -->

## Unbuffered vs. buffered, side by side

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Unbuffered**

```go
ch := make(chan int)
ch <- 1 // blocks until
        // someone receives
```

Strict handoff. Sender and
receiver synchronize on every
single value.

<!-- column: 1 -->

**Buffered, capacity 3**

```go
ch := make(chan int, 3)
ch <- 1 // doesn't block
ch <- 2 // doesn't block
ch <- 3 // doesn't block
ch <- 4 // NOW it blocks
```

Sender gets ahead of the
receiver, up to the buffer size.

<!-- reset_layout -->

**Demo:** print the send count immediately before each send on both
versions of the producer. Watch exactly which send is the one that
blocks.

<!--
speaker_note: |
  This is a good one to actually run rather than describe - add a
  fmt.Println("about to send", i) right before each ch <- i in both
  the unbuffered and buffered producer, then watch the printed count
  freeze at the point it blocks. It makes the abstract "capacity 3"
  concrete.
-->

<!-- end_slide -->

## `select`: waiting on multiple channels

`select` waits on several channel operations at once and runs whichever
one becomes ready first.

```go
select {
case v := <-ch1:
    fmt.Println("from ch1:", v)
case v := <-ch2:
    fmt.Println("from ch2:", v)
case <-time.After(2 * time.Second):
    fmt.Println("timed out")
default:
    fmt.Println("nothing ready yet")
}
```

<!-- pause -->

A `default` case makes `select` **non-blocking** — if nothing else is
ready immediately, it runs instead of waiting. Omit `default` and
`select` blocks until one of the other cases can proceed.

<!-- end_slide -->

## The switchboard operator

`select` is a switchboard operator watching several blinking lines at
once. They pick whichever line lights up first and connect it. With no
`default` case, they do nothing at all until one of the lines lights up
— they don't spin, they don't poll, they just wait.

<!-- pause -->

The `time.After(...)` case is what turns this into a **timeout pattern**:
race your real work against a timer, and if the timer wins, you know the
real work didn't show up in time.

<!-- end_slide -->

## A nil channel blocks forever — on purpose

Reading from or sending on a `nil` channel blocks forever. That sounds
useless until you see it inside a `select`.

```go
var maybeCh chan int // nil until assigned

select {
case v := <-ch1:
    fmt.Println("ch1:", v)
case v := <-maybeCh: // never becomes ready while nil
    fmt.Println("maybe:", v)
}
```

<!-- pause -->

A `select` case reading from a `nil` channel **never fires** — which is a
real, deliberate technique for conditionally disabling a case at
runtime. Set a channel variable to `nil` and its `select` case goes
silent without deleting any code.

<!-- end_slide -->

<!-- jump_to_middle -->

When it goes wrong
===

<!-- end_slide -->

## Deadlock: the whole program stops

A **deadlock** is goroutines waiting on each other in a cycle that can
never resolve. The simplest version: main sends on an unbuffered channel
that nobody ever receives from.

```go
func main() {
    ch := make(chan int)
    ch <- 42 // nobody is receiving — blocks forever
}
```

<!-- pause -->

**Demo — run this live.** When the runtime detects that *every* goroutine
in the program is asleep waiting on something, it doesn't hang silently
— it crashes loudly:

```
fatal error: all goroutines are asleep - deadlock!
```

<!--
speaker_note: |
  Genuinely run this one live rather than showing the text. Students
  retain "I watched the whole program crash with that exact message"
  far better than being told about it. Let it sit on screen for a
  couple of seconds before moving on.
-->

<!-- end_slide -->

## The deadlock the runtime *won't* catch you

The whole-program detector above only fires when **every** goroutine is
stuck. If two goroutines deadlock on each other while a third one is
still busy doing something else, the runtime sees the program as alive
and never reports anything.

<!-- pause -->

**That's silent, and much harder to spot** than the loud crash — no
error, no crash, just two goroutines that will never move again for the
rest of the process's life. This is exactly the same shape as the
goroutine leak two slides back: a subset of your program is dead, and
nothing tells you.

**When something in production "just hangs" with no error, this is one
of the first things to suspect.**

<!-- end_slide -->

## Race conditions are still fully possible

A channel and a mutex both prevent data races when used correctly. Go
doesn't prevent you from writing code that skips both:

```go
var counter int

func increment() {
    for i := 0; i < 1000; i++ {
        counter++ // read-modify-write, not atomic
    }
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 2; i++ {
        wg.Add(1)
        go func() { defer wg.Done(); increment() }()
    }
    wg.Wait()
    fmt.Println(counter) // rarely exactly 2000
}
```

<!-- pause -->

Two goroutines both doing `counter++` — a read, an add, a write — can
interleave. Some increments get lost. The final number is unpredictable
and usually wrong.

<!-- end_slide -->

## `go run -race`: watch it live

Go ships a **race detector** built into the toolchain. This is genuinely
one of Go's best practical tools — nothing shipped this conveniently in
Java or Python's standard tooling catches this class of bug for you.

```
go run -race ./cmd/race
```

<!-- pause -->

**Demo — run it live.** The output names the exact two goroutines, the
exact lines, and whether each access was a read or a write:

```
WARNING: DATA RACE
Write at 0x00c0000140a0 by goroutine 8:
  main.increment.func1() /path/main.go:14
Previous write at 0x00c0000140a0 by goroutine 7:
  main.increment.func1() /path/main.go:14
```

<!--
speaker_note: |
  Run this one live too, right after the deadlock demo - back to back
  they make a strong pair. Point at the two goroutine numbers and two
  line numbers in the real output on screen. Emphasize go test -race
  as something worth wiring into CI, not just a manual debugging tool.
-->

<!-- end_slide -->

## Fixing the race: two valid ways

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Option 1: `sync.Mutex`**

```go
var mu sync.Mutex
var counter int

func increment() {
    for i := 0; i < 1000; i++ {
        mu.Lock()
        counter++
        mu.Unlock()
    }
}
```

Simple, direct, correct for small
shared state like this.

<!-- column: 1 -->

**Option 2: channel + owner**

```go
incs := make(chan int)
go func() { // sole owner
    total := 0
    for range incs {
        total++
    }
}()
// workers send instead
// of touching counter
incs <- 1
```

Nobody but the owning goroutine
ever touches the count.

<!-- reset_layout -->

**Both are idiomatic.** The mutex is honestly the simpler fix here — a
reminder that "share memory by communicating" is a strong default, not a
rule you contort code to satisfy.

<!-- end_slide -->

<!-- jump_to_middle -->

Putting it together: worker pools
===

<!-- end_slide -->

## Worker pool: a fixed crew, a shared queue

A fixed number of worker goroutines pull jobs off one shared channel and
push results onto another. This is the pattern behind most real
concurrent Go services.

```go
jobs := make(chan int, 100)
results := make(chan int, 100)
var wg sync.WaitGroup

for w := 1; w <= 3; w++ { // 3 workers
    wg.Add(1)
    go func() {
        defer wg.Done()
        for j := range jobs {
            results <- j * j
        }
    }()
}

for j := 1; j <= 9; j++ { jobs <- j }
close(jobs) // no more work coming

wg.Wait()      // all workers drained `jobs` and exited
close(results) // safe to close now
```

<!-- pause -->

**The `WaitGroup` is what makes closing `results` safe.** Close it too
early and a worker still writing to it panics; the `wg.Wait()` guarantees
every worker has finished before that happens.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Goroutines are cheap**: KB-sized stacks, scheduled by the Go runtime, not the OS — thousands are normal
2. **No async/await**: concurrency is marked only by `go` and channel operations, not the function signature
3. **Channels are the idiomatic handoff**: unbuffered = synchronous rendezvous, buffered = slack up to capacity
4. **`select` waits on multiple channels**: `default` makes it non-blocking, `time.After` gives you timeouts
5. **Closing is a signal, not cleanup**: close from the sender, detect with comma-ok, never close twice
6. **Two failure modes to know by sight**: whole-program deadlock crashes loudly; races need `-race` to see
7. **Mutexes are still valid**: "share memory by communicating" is the default, not a law
8. **Test outcomes, not order**: a deterministic concurrency test asserts the aggregate result — a count, a sum, a full result set — never which goroutine ran first

<!-- end_slide -->

## Back to the opening scenario

Fetching from 10 independent sources in Java: you'd size a thread pool
before writing the fetch logic.

**In Go:** `go fetch(source)` ten times, and think instead about *how
results get back to you* — a channel to collect them, maybe a
`WaitGroup` to know when every fetch finished, maybe `select` with a
timeout so one slow source doesn't hold up the rest.

<!-- pause -->

**Type in chat: compare that to what you guessed at the start — how
close were you?**

<!--
speaker_note: |
  Read back a couple of the original guesses from the opening slide if
  you noted them. The point to land: Go shifts the design question from
  "how many workers do I provision" to "how does data flow between
  goroutines" - sizing still matters sometimes (worker pools exist!) but
  it's no longer the first thing you reach for.
-->

<!-- end_slide -->

## Bridge to Topic 8

**We've established:**

<!-- incremental_lists: true -->

- Goroutines are near-free to launch and scheduled by the Go runtime, not the OS
- Channels, not shared variables, are the idiomatic way goroutines coordinate
- Deadlocks and races are both real, and Go gives you concrete tools to catch each
- Testing concurrent code follows the same rule: assert the deterministic aggregate — never a specific goroutine execution order

<!-- incremental_lists: false -->

**Topic 8: Intro to Design Patterns** — now that you can write concurrent
Go, we'll look at how to structure larger programs around it: handlers,
controllers, services, and dependency injection.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
