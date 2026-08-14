# Lab 11: gRPC and Protocol Buffers -- a Booking service

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` -- don't look until you've had a go.

**Important, read this first:** this lab genuinely needs `protoc`,
`protoc-gen-go`, and `protoc-gen-go-grpc` installed on your own machine --
none of that toolchain is available in the environment these materials
were authored in. Exercise 1 has you install it as its first step. The
Go files under `starter/` and `solution/` that stand in for generated
code and the client/server that would use it are clearly marked
**ILLUSTRATIVE -- NOT GENERATOR OUTPUT** at the top of each file: they
show you the *shape* of what you'll get, but you must regenerate the
real thing with `protoc` before any of it will actually run. See
`code/README.md` (one directory up) for the full toolchain setup.

---

## Exercise 1: Install the toolchain and generate real code

**Objective:** Get a working `protoc` + Go plugin setup and generate
your first real `.pb.go` files.

**Context:** Everything from here on depends on this working. Budget
real time for it -- version mismatches between `protoc` and the plugins
are the most common source of confusing errors, and this is the one
topic in the course where "install a compiler and two plugins, get the
PATH right" can eat a meaningful chunk of a session.

**Tasks:**

1. Install `protoc` (the Protocol Buffer compiler -- a standalone
   binary, not a Go tool) and confirm with `protoc --version`.
2. Install the two Go plugins: `go install
   google.golang.org/protobuf/cmd/protoc-gen-go@latest` and `go install
   google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`. Confirm
   `$HOME/go/bin` is on your `PATH`.
3. Copy `starter/booking.proto` into your own working directory. Run:
   ```
   protoc --go_out=. --go_opt=paths=source_relative \
          --go-grpc_out=. --go-grpc_opt=paths=source_relative \
          booking.proto
   ```
4. Open the generated `booking.pb.go`. Find the generated `Booking`
   struct. Compare its shape to the hand-written, illustrative version
   in `starter/bookingpb/booking.go` -- what's different, and what did
   the illustrative version simplify away?

**Key Learning:** `protoc` and its Go plugins are a real, separate build
step before `go build` even runs -- this is the most "new tooling," not
just "new syntax," topic in the whole course. A `.proto` file compiles
to plain Go, but you need the external compiler to get there.

---

## Exercise 2: Implement the service, in memory

**Objective:** Implement the generated `BookingServiceServer` interface
against a mutex-protected in-memory store.

**Context:** This is a direct continuation of the store pattern from
Topic 10 -- same `sync.RWMutex`-protected `map`, same discipline, now
backing a gRPC server instead of an HTTP one.

**Tasks:**

1. Using your real generated types from Exercise 1 (not the
   illustrative stand-ins), write a `server` struct with a `map[string]*bookingpb.Booking`
   and a `sync.RWMutex`, seeded with two or three bookings.
2. Implement `GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest) (*bookingpb.Booking, error)`
   using `RLock`/`RUnlock` for the lookup, returning a `status.Errorf(codes.NotFound, ...)`
   error (from `google.golang.org/grpc/status` and `.../codes`) when the ID isn't found.
3. Register your implementation with a real `*grpc.Server` and listen
   on `:50051`. Use `cmd/server/main.go`'s commented-out wiring block as
   your starting point.

**Key Learning:** the interface changes (`BookingServiceServer` instead
of an `http.HandlerFunc`), but the concurrency-safety discipline from
Topic 10 doesn't -- a gRPC server fields concurrent calls just like a
REST server does.

---

## Exercise 3: A client that reads like a function call

**Objective:** Call `GetBooking` through the real generated client
stub and confirm it behaves like an ordinary synchronous call.

**Tasks:**

1. Write a small `main.go` that dials your running server with
   `grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))`
   and creates a client with `bookingpb.NewBookingServiceClient(conn)`.
2. Call `resp, err := client.GetBooking(ctx, &bookingpb.GetBookingRequest{Id: "abc123"})`
   and print the result.
3. Time the call (wrap it with `time.Now()`/`time.Since()`). Confirm to
   yourself this is a real network round trip, even though the code
   reads exactly like calling a local function.

**Key Learning:** the generated client stub is what makes gRPC calls
"feel" like ordinary Go function calls -- the network round trip is
real, but syntactically invisible.

---

## Exercise 4: Add a streaming RPC

**Objective:** Add `ListBookings`, a server-streaming RPC, and consume
it from a client until `io.EOF`.

**Tasks:**

1. Add `rpc ListBookings(ListBookingsRequest) returns (stream Booking);`
   to your `.proto` (an empty `ListBookingsRequest` message is fine).
   Regenerate with `protoc`.
2. Implement the streaming server method: loop over your store and call
   `stream.Send(b)` for each booking.
3. In your client, call `ListBookings`, then loop calling `stream.Recv()`
   until you get `io.EOF`, printing each booking as it arrives.

**Key Learning:** streaming RPCs have no clean REST/JSON equivalent --
REST is fundamentally request-then-response. gRPC has streaming as a
first-class, typed feature directly in the `.proto` service definition,
not something bolted on afterward with WebSockets or SSE.

---

## Exercise 5: Break schema evolution on purpose, then fix it

**Objective:** Understand why field numbers are permanent contracts by
deliberately violating that rule, then doing it correctly.

**Tasks:**

1. In your `.proto`, delete the `nights` field entirely (field number
   3) and add a brand-new field, say `string status = 3;`, reusing the
   same number.
2. Regenerate. Reason through (read a little protobuf documentation if
   you want more detail): why is this dangerous for any client still
   running against the *old* generated code, even though your new code
   compiles just fine?
3. Redo it correctly: remove the field, mark its number `reserved 3;`
   (and optionally `reserved "nights";`), and add your new `status`
   field with the *next* unused number instead.

**Key Learning:** protobuf's numbered fields aren't a limitation --
they're the actual mechanism that lets you add fields and deploy new
servers gradually while old clients keep working, unaffected, because
they simply don't know about a new number and safely ignore it. JSON
gives you no equivalent structural guarantee.

---

## Exercise 6: Build the same lookup twice -- REST and gRPC

**Objective:** Implement the identical `GetBooking`-style lookup as
both Topic 10's REST endpoint and this topic's gRPC unary call, and
write down what's genuinely different versus fundamentally the same.

**Tasks:**

1. If you still have it, reuse your Topic 10 `GET /bookings/{id}`
   handler. If not, rebuild a minimal version against the same
   in-memory `Booking` data you're using here.
2. Run both the REST version and the gRPC version side by side against
   the same underlying data.
3. Write down, concretely: what's genuinely *different* between the two
   (wire format, the codegen step, the separate `.proto` contract file)
   versus what's fundamentally the *same* idea (a typed request comes
   in, a typed response goes out).

**Key Learning:** gRPC and REST solve the same fundamental problem
(typed request in, typed response out, over a network) with different
tooling and different guarantees. Neither is "wrong" -- the choice
between them is a real, contextual engineering decision, not a
correctness question.

---

## Exercise 7: Prove it with a test

**Objective:** Write a test for `GetBooking`'s business logic by
calling it directly as a plain Go function -- no server, no network,
no `grpc.Dial`.

**Context:** Every generated gRPC method is just a Go method with
`ctx context.Context` as its first parameter -- Exercise 2 already
showed you that signature, and the "Every generated method takes
context.Context first" slide named it as a Go-wide idiom, not a gRPC
quirk. That means testing `GetBooking` needs nothing gRPC-specific: no
listener, no dialed connection, no real transport -- just construct
the `server`, call `GetBooking` directly, and assert on what comes
back, exactly like testing any other method you've written this
course. This is a genuinely good reason the "test as you go" habit
from Topic 2 is easy to keep here: gRPC's network plumbing lives
entirely in the parts you're *not* testing.

**Tasks:**

1. In `cmd/server/server_test.go`, write `TestGetBooking_Found`:
   construct a server with `newServer()`, call
   `s.GetBooking(context.Background(), &bookingpb.GetBookingRequest{Id: "abc123"})`
   directly, and assert the returned `*bookingpb.Booking` has the
   `Guest` and `Nights` you expect -- `t.Fatalf` if `err != nil`,
   `t.Errorf` if a field doesn't match.
2. Write `TestGetBooking_NotFound`: call `GetBooking` with an ID that
   isn't in the store and assert `err != nil`.
3. Confirm (by reading through the logic, since there's no local
   toolchain here) that both tests pass against
   `solution/cmd/server/main.go`'s implementation.
4. Deliberately break the not-found path -- e.g. change the `if !ok`
   block so it falls through and returns `b, nil` (the nil zero value)
   instead of an error -- and confirm `TestGetBooking_NotFound` would
   now fail. Put the correct code back afterward.

**Key Learning:** the first parameter of
`GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest)` is
what makes this possible -- it's an ordinary Go method, so it's
testable exactly like any other method, with a direct function call
and no gRPC transport involved. This isn't a workaround specific to
the illustrative code in this lab, either: it works because of the
method signature `protoc` generates, so once you've regenerated real
code in Exercise 1, this exact same test still compiles and passes
unchanged.

---

## Summary

By the end of this lab you should be able to:

- Set up and run `protoc` plus the two Go plugins, and explain what
  each generated file is for
- Implement a generated gRPC service interface against a
  mutex-protected in-memory store
- Call a unary RPC through a generated client stub and explain why it
  "reads like a function call"
- Implement and consume a server-streaming RPC, reading until `io.EOF`
- Explain, concretely, why reusing a protobuf field number is dangerous
  and how `reserved` prevents it
- Compare a REST implementation and a gRPC implementation of the same
  operation and name what's different versus what's the same
- Test a gRPC service method directly as a plain Go function call, no
  `grpc.Dial` or running server required, because `context.Context`
  first is all the pattern needs
