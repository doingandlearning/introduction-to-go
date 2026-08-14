# Topic 7 sample code — Concurrency

Seven small programs, one module. Run everything from this directory
(`code/`).

Five are safe to run any number of times. **Two are deliberately broken**
— `cmd/deadlock` and `cmd/race` — built to demonstrate specific failure
modes live. Don't wire either into CI; they exist to be run once,
on purpose, in front of the room.

## `cmd/waitgroup`

Five goroutines, a `sync.WaitGroup`. Demonstrates that `WaitGroup`
guarantees *completion*, not *order*.

```
go run ./cmd/waitgroup
```

Run it several times back to back — the "worker N done" lines print in a
different order almost every run.

## `cmd/producer`

Producer/consumer over an **unbuffered** channel. One goroutine sends five
numbers, main receives and prints them with `range`, the producer closes
the channel when it's done.

```
go run ./cmd/producer
```

## `cmd/buffered`

Same shape as `cmd/producer`, but the channel has capacity 3. Prints the
send count immediately before each send so you can see exactly which send
is the first to block.

```
go run ./cmd/buffered
```

Compare its output against `cmd/producer` — the first three sends return
immediately here; in the unbuffered version, every single send waits for
a receive.

## `cmd/selectdemo`

`select` waiting on two input channels plus a `time.After` timeout case.
Neither input channel is ever sent on, so the timeout always fires after
~2 seconds.

```
go run ./cmd/selectdemo
```

## `cmd/workerpool`

A fixed pool of 3 worker goroutines pulling jobs off a shared `jobs`
channel and pushing squared results onto a shared `results` channel. A
`sync.WaitGroup` tracks when every worker has finished, which is what
makes it safe to close `results`.

```
go run ./cmd/workerpool
```

## `cmd/deadlock` — INTENTIONALLY BROKEN

Sends on an unbuffered channel that nobody ever receives from. The Go
runtime detects that every goroutine in the program is asleep and crashes
loudly instead of hanging silently.

```
go run ./cmd/deadlock
```

Expected output:

```
fatal error: all goroutines are asleep - deadlock!
```

This only fires when the *whole* program is stuck. A deadlock between a
subset of goroutines, while others stay busy, is not detected — that
version is silent and much harder to spot. See the comment at the top of
`main.go` for more.

## `cmd/race` — INTENTIONALLY BROKEN

Two goroutines increment a shared `counter` a thousand times each with no
synchronization. Run it plain first, then with the race detector:

```
go run ./cmd/race            # "succeeds" with a wrong number, usually
go run -race ./cmd/race      # reports the exact conflicting accesses
```

`go run -race` / `go test -race` is one of Go's best practical tools —
there's no equivalent shipped this conveniently in Java or Python's
standard tooling. The two idiomatic fixes (a `sync.Mutex`, and a
redesign using a channel + single owning goroutine) are kept as commented
code at the bottom of `main.go` so you can uncomment and compare rather
than needing separate binaries.

## Formatting

```
gofmt -l .    # lists files that don't match gofmt's formatting
gofmt -w .    # rewrites them in place
```
