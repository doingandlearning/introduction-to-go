# Lab 7: Concurrency — Riverside Courier Co.

Riverside Courier Co. dispatches couriers, tracks deliveries, and needs a
small Go system to coordinate all of it concurrently. Starter code is in
`starter/` (TODOs to fill in). A complete reference is in `solution/` —
don't look until you've had a go.

Exercises 2, 3, 4, and 8 ship a **pre-written, already-failing test** in
`starter/`. Run `go test ./...` (or the scoped command each exercise
gives you) before you write anything — read the failure, then implement
until it's green. Exercises 1, 5, 6, 7, and 9 don't get a test: 1 and 5
are entirely inline in `main()` with nothing separable to call from a
test; 6 and 7 are deliberately-broken standalone programs meant to be
**run and observed**, not completed, in `standalone/deadlock/` and
`standalone/race/`; 9 is optional/stretch and depends on real elapsed
time, which doesn't lend itself to a fast unit test.

One concurrency-specific wrinkle worth knowing going in: if a producer
function forgets to send or forgets to close its channel, the pre-written
test for it can't fail with a normal assertion — nothing ever arrives.
Those tests use `select` with a short timeout instead, so a stuck
implementation fails fast with a clear message rather than hanging until
`go test`'s own multi-minute default timeout kills the whole run. That's
not incidental — it's the exact pattern Exercise 5 teaches you to write
by hand a little further down.

By the end of this lab you will have launched goroutines, used `defer` to
guarantee per-goroutine cleanup, coordinated goroutines with `WaitGroup` and
channels (both unbuffered and buffered), used `select` with a timeout,
watched a real deadlock and a real data race happen, fixed the race two
different ways, and built a small worker pool.

This lab runs across two sessions:

- **Part 1 (today)** — Exercises 1–2, goroutines and `defer`.
- **Part 2 (next session)** — Exercises 3–9, channels onward.

---

# Part 1 — Today: goroutines & defer

## Exercise 1: Couriers checking in

**Objective:** Launch goroutines and use `sync.WaitGroup` to wait for all
of them.

**Context:** `starter/cmd/dispatch/main.go` has a TODO where five couriers
should each report in concurrently before the dispatcher can start the
day. The whole exercise lives inline in `main()`, so there's no separate
function for a pre-written test to call — you'll verify this one by
running it, not `go test`.

**Tasks:**

1. In the TODO, launch 5 goroutines. Each one should call
   `checkIn(courierID)`, which prints `"courier <n> checked in"`.
2. Use a `sync.WaitGroup` so `main` waits for all 5 to finish before
   printing `"all couriers checked in, starting the day"`.
3. Run `go run ./cmd/dispatch` five or six times in a row. Note that the
   order of the check-in lines changes between runs.

**Key Learning:** `WaitGroup` guarantees every goroutine finished before
`Wait()` returns — it guarantees nothing about the order they finished
in. If you need a specific order, `WaitGroup` is the wrong tool.

---

## Exercise 2: Checking a courier back out, safely

**Objective:** Use `defer` inside a goroutine to guarantee cleanup runs —
in LIFO order, and even when the goroutine returns early.

**Context:** `starter/cmd/shiftend/main.go` has a TODO in
`checkOutCourier`. Each courier goroutine "opens" two resources in
sequence — a radio channel, then a route log — and both need to be closed
out however the function returns, including the early-return path for a
courier whose route got blocked. `checkOutCourier` takes an `io.Writer`
instead of printing straight to stdout, so `starter/cmd/shiftend/main_test.go`
can capture what it prints with a `bytes.Buffer` and check it — `main`
already passes `os.Stdout` for you.

**Tasks:**

1. Run `go test ./cmd/shiftend/...`. Both tests fail — read why.
2. Implement `checkOutCourier(w io.Writer, courierID int, wg *sync.WaitGroup)`:
   - `defer wg.Done()` first, same pairing as Exercise 1.
   - `Fprintf(w, "courier <id>: radio open")`, then immediately `defer`
     an `Fprintf` of `"courier <id>: radio closed"`.
   - `Fprintf(w, "courier <id>: route log started")`, then immediately
     `defer` an `Fprintf` of `"courier <id>: route log filed"`.
   - If `courierID == 3`, `Fprintf(w, "courier 3: route blocked, aborting")`
     and `return` right there — before the final delivery print.
   - Otherwise, `Fprintf(w, "courier <id>: delivering final parcel")`.
3. Re-run `go test ./cmd/shiftend/...` and confirm both tests pass.
4. Before running the real program, predict on paper the order of the
   five lines a single *non*-blocked courier prints. Then run
   `go run ./cmd/shiftend` and check your prediction against the
   interleaved output for courier 3 specifically.
5. Confirm from that output that courier 3 still prints both `"route log
   filed"` and `"radio closed"` even though it returned early and never
   reached the delivery line — the same thing `TestCheckOutCourierEarlyReturn`
   already proved for you.
6. In a comment above `checkOutCourier`, answer: with 5 goroutines running
   this function at once, what would go wrong if you deleted both `defer`
   statements and instead put the two "closed"/"filed" prints at the very
   end of the function, after the delivery print?

**Key Learning:** Each goroutine has its own independent defer stack —
`defer` doesn't just guarantee cleanup runs on early return in a
single-threaded function (Topic 2), it guarantees it *per goroutine*,
regardless of which of several concurrently-running goroutines happens to
take the early-return path. Put cleanup at the end of the function instead
and any goroutine that returns early skips it — with 5 goroutines running
at once, that's 5 independent chances to leak a resource, not one.

---

# Part 2 — Next session: channels onward

## Exercise 3: One order at a time

**Objective:** Build a producer/consumer pair over an **unbuffered**
channel.

**Context:** `starter/cmd/courier/main.go` has a TODO. A dispatcher
goroutine needs to send 5 order IDs down a channel, one at a time, to a
courier (main) that prints each as it arrives.

**Tasks:**

1. Run `go test ./cmd/courier/...`. It fails — right now it times out
   after 2 seconds rather than failing instantly, because nothing has
   sent or closed the channel yet. That's expected: a producer that
   never sends can't be caught by a normal assertion, only by a timeout.
2. Implement `dispatchOrders(ch chan<- int)` so it sends order IDs `1`
   through `5` into `ch`, then closes the channel.
3. In `main`, add `"fmt"` to the import block, then receive from the
   channel with `range` and print `"delivering order <n>"` for each.
4. Re-run `go test ./cmd/courier/...` and confirm it passes immediately
   now, no timeout.
5. Run `go run ./cmd/courier`. Confirm all 5 orders print, in order, and
   the program exits cleanly (it wouldn't if you forgot to close the
   channel).

**Key Learning:** An unbuffered channel is a synchronous handoff — the
sender's `ch <- n` doesn't complete until the receiver is ready for it.
`close()` is the sender's job, and `range` exits automatically once the
channel is closed and drained.

---

## Exercise 4: Give the dispatcher some slack

**Objective:** Do the same job with a **buffered** channel, and reason
about exactly when it starts blocking.

**Context:** `starter/cmd/courierbuffered/main.go` has a TODO. Same
scenario as Exercise 3, but the dispatcher should be able to get a little
ahead of the courier. The pre-written test only checks the values that
arrive and that the channel closes — it can't (and doesn't try to) assert
anything about print timing; that's what task 4 below is for.

**Tasks:**

1. Run `go test ./cmd/courierbuffered/...`. Same timeout-shaped failure
   as Exercise 3, for the same reason.
2. Add `"fmt"` to the import block, then implement the same send-5-orders
   logic as Exercise 3, but this time over a channel created with
   capacity 3: `make(chan int, 3)`. Print `"about to send order <n>"`
   immediately before each send, then close the channel when done.
3. Re-run `go test ./cmd/courierbuffered/...` and confirm it passes.
4. Run `go run ./cmd/courierbuffered`. Compare the printed order against
   Exercise 3 — in this version, how many "about to send" lines print
   before the first "delivering" line has a chance to appear?
5. Answer in a comment above `main`: which send is the first one that
   would block if the courier were slow to start receiving, and why?

**Key Learning:** A buffered channel lets the sender run ahead of the
receiver, up to the buffer's capacity — after that, sends block exactly
like an unbuffered channel would on every send.

---

## Exercise 5: Don't wait forever for a courier

**Objective:** Use `select` with a timeout.

**Context:** `starter/cmd/radio/main.go` has a TODO. Dispatch is waiting
on a status update from either of two couriers, but if neither responds
within 2 seconds, it should give up and log a timeout instead of hanging.
Like Exercise 1, this TODO is entirely inline in `main()` — no test file
here, you're building and observing this one directly. (It's also the
exact pattern the pre-written tests in Exercises 3, 4, and 8 were already
using behind the scenes to protect themselves against a hung producer.)

**Tasks:**

1. Create two `chan string` values, `courierA` and `courierB`. Do **not**
   send on either of them anywhere in the program.
2. Implement a `select` with a case for each channel and a
   `case <-time.After(2 * time.Second):` case that prints
   `"no response from either courier, escalating"`.
3. Run `go run ./cmd/radio` and confirm it waits roughly 2 seconds before
   printing the timeout message.

**Key Learning:** `select` with a `time.After` case turns "wait for a
channel" into "wait for a channel, but not forever" — a pattern you'll
reach for constantly once you're writing anything that talks to the
outside world (including, as you've just seen, your own tests).

---

## Exercise 6: Watch a deadlock happen

**Objective:** Trigger a real whole-program deadlock and read the exact
runtime error.

**Context:** `standalone/deadlock/main.go` is already complete — it's
built to fail. Main sends on an unbuffered channel that nothing ever
receives from.

**Tasks:**

1. Run `go run ./standalone/deadlock`.
2. Read the error output line by line. Identify: which goroutine is
   reported stuck, what operation it's stuck on (`chan send`), and which
   line of source it names.
3. In a comment at the bottom of the file, write one sentence describing
   a scenario where a deadlock like this would happen between *some* of
   your goroutines, but not all of them — and explain why the runtime
   would *not* report that version.

**Key Learning:** The Go runtime detects and crashes loudly on a
whole-program deadlock — but a deadlock affecting only a subset of your
goroutines, while others stay busy, is invisible to this detector. It's
one of the most important "silent failure" shapes to recognize by
symptom (a hang, no error) even when you can't see it directly.

---

## Exercise 7: Watch a race happen, then fix it twice

**Objective:** Trigger a data race, read the race detector's report, and
fix it two different ways.

**Context:** `standalone/race/main.go` is already complete. Two
goroutines both increment a shared `deliveryCount` variable a thousand
times each, with no synchronization.

**Tasks:**

1. Run `go run ./standalone/race` a few times. Note that the final count
   printed is rarely `2000` and changes between runs.
2. Run `go run -race ./standalone/race`. Read the `WARNING: DATA RACE`
   report — find the two goroutine IDs and the two line numbers it
   names.
3. Copy the file to `standalone/race/main_mutex.go` (adjust the package
   as needed, or work in a scratch copy) and fix the race using a
   `sync.Mutex` around the increment.
4. In a second copy, fix it by redesigning: each goroutine sends `1` down
   a channel instead of touching the shared variable directly, and a
   single owning goroutine receives from that channel and keeps the
   running total.
5. Run `go run -race` against both fixes and confirm neither reports a
   race.

**Key Learning:** `go run -race` / `go test -race` catches a whole class
of bug that compiles cleanly and often "looks fine" in casual testing.
Both the mutex fix and the channel fix are idiomatic Go — the mutex is
frequently the simpler choice for small shared state like a single
counter, which is worth saying out loud given how often "use a channel"
gets treated as the only correct answer.

---

## Exercise 8: The dispatch pool

**Objective:** Build a small worker pool — the main exercise this lab is
built around — and prove its result with a test that asserts an
*aggregate outcome*, not a goroutine execution order.

**Context:** `starter/cmd/workerpool/main.go` has TODOs, and
`starter/cmd/workerpool/workerpool_test.go` already has a full,
pre-written test waiting for you. Riverside wants a fixed crew of 3
dispatcher goroutines pulling delivery jobs off a shared `jobs` channel
and pushing completion records onto a shared `results` channel. Testing
concurrent code is different from testing `divide` back in Topic 2: you
cannot assert *which* worker processed job `7`, or the exact order
results arrive in `results` — the scheduler doesn't guarantee either, and
a test that assumes one will flake sooner or later. What you *can*
assert, every single run, is the aggregate: exactly 12 results come back,
and the complete set of fees is `{2, 4, 6, ..., 24}` — that's true no
matter which worker did the work or what order it finished in.

**Tasks:**

1. Run `go test ./cmd/workerpool/...`. It fails — right now it times out
   rather than failing instantly, for the same reason as Exercises 3 and
   4: a worker pool that never finishes can't be caught by a normal
   assertion, only a timeout.
2. Implement `worker(id int, jobs <-chan int, results chan<- int, wg
   *sync.WaitGroup)` so each worker ranges over `jobs`, and for each job
   `j` it receives, sends `j * 2` (a stand-in "delivery fee") to
   `results`, then calls `wg.Done()` via `defer` when the channel closes
   and the range loop ends.
3. In `main`, launch exactly 3 workers, all sharing the same `jobs` and
   `results` channels.
4. Send job IDs `1` through `12` into `jobs`, then close it.
5. Use the `WaitGroup` to know when every worker has finished, **then**
   close `results` — not before.
6. Range over `results` and print each one.
7. Re-run `go test ./cmd/workerpool/...` and confirm it passes.
8. Run `go run ./cmd/workerpool` a few times. Confirm you always get 12
   results, even though the order changes between runs.
9. Optional regression check: deliberately break the *aggregate logic*,
   not the timing — change `fee := j * 2` to `fee := (j - 1) * 2`, an
   off-by-one in the fee calculation. Re-run `go test ./cmd/workerpool/...`
   and confirm your test fails and reports exactly which value is wrong.
   Fix it back.

**Key Learning:** The `WaitGroup` isn't decoration here — it's the
mechanism that makes closing `results` safe. Close it while a worker is
still mid-send and that worker panics. This pattern (fixed workers, one
input channel, one output channel, a `WaitGroup` gating the close) is the
backbone of most real concurrent Go services. And testing it means
testing outcomes that don't depend on scheduling order — a count, a sum,
or a complete result set — never an assumption about which goroutine ran
first or which order values arrived in. A test that asserts a specific
interleaving isn't testing your program, it's testing today's scheduler,
and it will flake the moment that changes.

---

## Exercise 9 (optional/stretch): Throttle the pool and count deliveries lock-free

**Objective:** Extend the worker pool with two more concurrency
primitives that don't come up in the required path: `sync/atomic` for a
lock-free counter, and `time.Ticker` for a fixed-rate limiter.

**Context:** `starter/cmd/throttledpool/main.go` is a copy of the
Exercise 8 worker pool, already wired up, with TODOs only for the two new
parts. Riverside's dispatch radio can only take 5 status updates a second
before the base station starts dropping packets, and ops wants a live
count of total deliveries processed across all workers without adding
another mutex. There's no pre-written test for this one — it genuinely
takes over 2 real seconds to run by design, which doesn't suit a fast
unit test, so verify it by running it.

**Tasks:**

1. Add an `*atomic.Int64` parameter to `worker`, and call `.Add(1)` on it
   once per completed job. In `main`, declare `var delivered atomic.Int64`
   and pass `&delivered` into each worker.
2. Before the loop that sends jobs, create a `*time.Ticker` firing every
   `time.Second / 5` and `defer` its `Stop()`. Receive from `ticker.C`
   once immediately before each send, so jobs reach the workers no faster
   than 5 per second, however fast the workers could otherwise drain
   them.
3. After the pool finishes, print `delivered.Load()` alongside the
   elapsed time. Run `go run ./cmd/throttledpool` and confirm it takes at
   least ~2.2 seconds for 12 jobs at 5/sec, and that the delivered count
   is 12.
4. In a comment above `main`, answer: the package also has a
   variable-free shortcut, `time.Tick(...)`, that skips creating a
   `*Ticker` you can `Stop()`. Why would using `time.Tick` instead be a
   real problem in a long-running service, even though it works fine in
   this short-lived program?

**Key Learning:** `sync/atomic` gives lock-free updates for simple shared
state — cheaper than a `sync.Mutex` when the shared state truly is "just
a number," the same trade-off from Exercise 8 but leaning the other way.
A `*time.Ticker`, stopped explicitly, turns "as fast as possible" into
"at most N per second" — `time.Tick`'s underlying timer is never garbage
collected and never stoppable, so using it inside anything longer-lived
than a `main` function that exits on its own is a real, if slow, resource
leak.

---

## Summary

By the end of this lab you should be able to:

- Launch goroutines and use `sync.WaitGroup` to wait for a batch of them,
  without assuming anything about their completion order
- Use `defer` inside a goroutine to guarantee per-goroutine cleanup runs
  in LIFO order, however that goroutine returns — including early
- Build a producer/consumer pair over both unbuffered and buffered
  channels, and explain exactly when a buffered channel starts blocking
- Use `select` with `time.After` to bound how long you wait on a channel
  — including inside a test, to stop a stuck goroutine from hanging the
  whole run
- Recognize a whole-program deadlock from its exact runtime error, and
  explain why a partial deadlock wouldn't be caught the same way
- Trigger and read a `go run -race` report, and fix a data race with
  either a `sync.Mutex` or a channel-and-owning-goroutine redesign
- Build a worker pool where a `WaitGroup` safely gates when the results
  channel can be closed
- Read a pre-written, failing test for concurrent code and make it pass
  by asserting the aggregate outcome (a count, a sum, or a full result
  set) rather than the order goroutines actually ran in
- (optional/stretch) Use `sync/atomic` for a lock-free counter and a
  `*time.Ticker` to rate-limit work, and explain why `time.Tick` alone is
  a leak risk outside a short-lived program
