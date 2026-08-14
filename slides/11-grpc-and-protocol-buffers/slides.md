---
title: "**gRPC and Protocol Buffers**"
sub_title: Go Programming — Topic 11
author: Kevin Cunningham
---

## Opening scenario

Topic 10's REST service had a struct tag typo -- `json:"geust"` instead
of `json:"guest"`. Nothing crashed. Nothing errored. The `guest` field
just silently came back empty on every response, and nobody noticed
until a support ticket showed up days later.

**What would have caught that before it ever shipped?**

We'll come back to this once you've seen the alternative.

<!--
speaker_note: |
  Let a few answers land - "tests," "code review," "a schema" are all
  fair, partial answers. Don't resolve it yet. The honest answer this
  topic gives is "a compiler that refuses to build the wrong thing,"
  which is a stronger guarantee than any of the partial answers alone.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Why protobuf over JSON?
===

<!-- end_slide -->

## The client asked for this section specifically

Before any gRPC mechanics: **why would you choose protobuf as your
data format at all**, instead of the JSON you already know from Topic
10?

This isn't a one-line answer. There are four genuinely separate
reasons, and they don't all carry the same weight for every team.

<!-- pause -->

Go through all four before deciding which one matters most to you.

<!--
speaker_note: |
  Signpost this explicitly as its own block - the client flagged that
  earlier material treated this too thinly. Take your time on the next
  four slides; don't compress them to hit a schedule.
-->

<!-- end_slide -->

## Reason 1: binary wire format

JSON is text. Every single message repeats the field *names*, spelled
out, every time:

```json
{"id":"abc123","guest":"A. Ortiz","nights":3}
```

`"guest":` costs 9 bytes -- on every message, forever, for every
booking your service ever sends.

<!-- pause -->

Protobuf's binary encoding puts the field **number** on the wire, not
the name, and packs values compactly (varint encoding for small
integers, no repeated key strings at all). The same data is
meaningfully smaller and faster to parse.

**This matters at scale**: high request volume, bandwidth-constrained
mobile clients, or heavy internal service-to-service traffic, where
"smaller and faster" compounds into real infrastructure cost.

<!-- end_slide -->

## Reason 2: a schema the compiler actually checks

JSON has no schema unless you bolt one on separately -- JSON Schema, or
just hope, which is exactly what let Topic 10's `"geust"` typo through
unnoticed.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// JSON: typo compiles fine,
// field silently stays empty
type Booking struct {
    Guest string `json:"geust"`
}
```

<!-- column: 1 -->

```go
// protobuf: Nights is int32.
// Passing a string here is a
// Go COMPILE error, full stop.
b := &Booking{
    Nights: "three", // won't build
}
```

<!-- reset_layout -->

**A `.proto` file IS the contract.** `protoc` generates real Go types
from it -- get the type wrong, and you find out at `go build`, not from
a support ticket days later.

<!--
speaker_note: |
  This is the direct payoff of the opening provocation - name that
  connection explicitly here if nobody's made it yet.
-->

<!-- end_slide -->

## Reason 3: cross-language by design

The same `.proto` file generates idiomatic client and server code in
Go, Java, Python, C++, and more -- from one source file, not one
per language.

<!-- pause -->

Genuinely useful the moment your Go backend needs to talk to a Python
data-science service, or a mobile app's Java or Swift client. One
source-of-truth contract, instead of every team hand-maintaining its
own JSON (de)serialization logic and hoping everyone stays in sync.

**JSON gives you no equivalent** -- every language team writes and
maintains its own struct/class definitions independently, by
convention, with nothing enforcing agreement between them.

<!-- end_slide -->

## Reason 4: disciplined, backward-compatible evolution

This is the strongest, most underrated reason -- and it's a direct
consequence of the field-number rule you'll meet properly later this
topic.

<!-- pause -->

Field numbers are permanent contracts. That's not a limitation, it's
the actual mechanism that lets you:

<!-- incremental_lists: true -->

- Add a new field to a `.proto`
- Deploy new servers gradually, not all at once
- Have **old clients keep working, unaffected** -- they simply don't
  know about the new field number and safely ignore it

<!-- incremental_lists: false -->

**JSON gives you no structural equivalent.** Nothing stops one team
from silently renaming a JSON key and quietly breaking every consumer,
with no warning at all.

<!-- end_slide -->

## Side by side, and a genuine debate

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**JSON (text)**

```json
{"id":"abc123","guest":"A. Ortiz","nights":3}
```

~46 bytes. Field names spelled out, every time, on every message.

<!-- column: 1 -->

**Protobuf (binary, conceptually)**

`[0x0a][len]abc123 [0x12][len]A. Ortiz [0x18]3`

Roughly half the size or less -- field *numbers*, not names, packed
values, no repeated key text.

<!-- reset_layout -->

**Type in chat:** for an internal service handling millions of
requests a day, does the **size difference** matter more, or does the
**compile-time-checked contract** matter more?

<!--
speaker_note: |
  Genuinely debatable - let both sides land before moving on. Size
  matters more when bandwidth or request volume dominates cost
  (mobile, very high QPS internal services). The contract matters more
  when the real cost driver is engineering time lost to integration
  bugs and silent breakage. Both are real reasons teams choose
  protobuf - don't declare a winner.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

What Go disagrees with: gRPC mechanics
===

<!-- end_slide -->

## Protocol Buffers, mechanically

Define messages and services in a `.proto` file -- a separate format,
not Go -- then run a compiler plus a Go plugin to generate matching Go
types and interfaces.

```protobuf
// booking.proto
syntax = "proto3";
package booking;

message Booking {
  string id = 1;
  string guest = 2;
  int32 nights = 3;
}

message GetBookingRequest { string id = 1; }

service BookingService {
  rpc GetBooking(GetBookingRequest) returns (Booking);
}
```

```
protoc --go_out=. --go-grpc_out=. booking.proto
```

<!-- end_slide -->

## This is the most "new tooling" topic in the course

Everything up to now has been "here's another Go language feature."
This is genuinely different.

<!-- incremental_lists: true -->

- A separate file format (`.proto`, not Go)
- An external compiler (`protoc`) -- not part of the Go toolchain
- Generated `.pb.go` files as a real build step, **before** `go build`
  even gets involved

<!-- incremental_lists: false -->

If you've used Avro or protobuf already in a Java shop, the *concept*
feels familiar. If not, expect this to feel like more setup friction
than anything else in the course -- because it genuinely is.

<!-- end_slide -->

## Say this before anyone's tempted: don't hand-edit generated files

`.pb.go` files are readable, syntactically normal Go. That makes them
tempting to "just fix" directly when something looks wrong.

**Don't.**

<!-- pause -->

The `.proto` file is the actual source of truth. Regenerating from it
**silently overwrites** any direct edit made to the generated code.

Same discipline as not hand-editing a compiled `.class` file or a
minified JS bundle -- it just doesn't *look* like a build artifact the
way those obviously do.

<!--
speaker_note: |
  Say this now, before the lab, not after someone's already lost an
  edit to a regeneration. The lesson lands either way but costs real
  frustration if learned the hard way.
-->

<!-- end_slide -->

## What gets generated

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// server: implement the
// generated interface
type server struct {
    booking.
      UnimplementedBookingServiceServer
    store map[string]*booking.Booking
    mu    sync.RWMutex
}

func (s *server) GetBooking(
    ctx context.Context,
    req *booking.GetBookingRequest,
) (*booking.Booking, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    b, ok := s.store[req.Id]
    if !ok {
        return nil, status.Errorf(
            codes.NotFound,
            "booking %s not found", req.Id)
    }
    return b, nil
}
```

<!-- column: 1 -->

```go
// client: looks like an
// ordinary function call
conn, _ := grpc.Dial(
    "localhost:50051",
    grpc.WithTransportCredentials(
        insecure.NewCredentials()))

client := booking.
    NewBookingServiceClient(conn)

resp, err := client.GetBooking(
    ctx,
    &booking.GetBookingRequest{
        Id: "abc123",
    })
```

<!-- reset_layout -->

**Demo:** show the illustrative `code/cmd/server` and `code/cmd/client`
files -- clearly marked as hand-written stand-ins, not generator
output, since this sandbox has no `protoc`.

<!-- end_slide -->

## Every generated method takes context.Context first

Look again at `GetBooking(ctx context.Context, req *booking.GetBookingRequest)`.

That first parameter isn't a gRPC quirk.

<!-- pause -->

**It's a Go-wide idiom** you're meeting here in its most visible form:
`context.Context` carries cancellation signals, deadlines, and
request-scoped values through a call chain. Any Go function that does
I/O, or might need to be cancelled partway through, conventionally
takes `ctx context.Context` as its first argument -- gRPC just makes
you feel this immediately, because literally every generated client
and server method is built around it.

<!-- end_slide -->

## Four call shapes

gRPC sits on top of protobuf and HTTP/2, and supports four shapes:

<!-- incremental_lists: true -->

- **Unary** -- one request, one response. The REST-like default.
- **Server streaming** -- one request, a stream of responses.
- **Client streaming** -- a stream of requests, one response.
- **Bidirectional streaming** -- both sides stream, independently.

<!-- incremental_lists: false -->

Today: unary (just seen) and server streaming, in depth.

<!-- end_slide -->

## Server streaming, and why REST can't do this cleanly

```protobuf
service BookingService {
  rpc GetBooking(GetBookingRequest) returns (Booking);
  rpc ListBookings(ListBookingsRequest) returns (stream Booking);
}
```

<!-- pause -->

**REST is fundamentally request-then-response.** Getting anything
stream-shaped out of it means reaching for WebSockets, Server-Sent
Events, or long-polling -- none of which are REST itself, all bolted on
after the fact.

gRPC has streaming as a **first-class, typed feature**, expressed
directly in the `.proto` service definition. Anything genuinely
continuous -- live updates, large result sets you don't want to buffer
entirely in memory, chat-like exchanges -- fits this model naturally.

<!-- end_slide -->

## Exchanging letters vs. keeping a line open

REST/JSON mails one letter, waits for one reply letter, repeats --
every exchange is its own complete, self-contained transaction.

<!-- pause -->

gRPC streaming opens a line and keeps it open: either side can keep
talking for as long as the call stays connected.

**Demo:** run the illustrative server-streaming client in
`code/cmd/client/main.go` -- watch it call `stream.Recv()` in a loop
until `io.EOF`, exactly the pattern you'll write for real once you've
generated real code with `protoc`.

<!--
speaker_note: |
  This is exactly why a live stock ticker or chat feed fits naturally
  on a phone line but feels forced as a sequence of mailed letters -
  worth saying out loud even though it's not on the slide.
-->

<!-- end_slide -->

## Field numbers are permanent contracts

Protobuf fields are numbered, typed, and that numbering is **permanent**
once shipped.

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```protobuf
// Dangerous: reusing 3
message Booking {
  string id = 1;
  string guest = 2;
  string status = 3; // was
                      // "nights"!
}
```

<!-- column: 1 -->

```protobuf
// Correct: reserve, don't reuse
message Booking {
  string id = 1;
  string guest = 2;
  reserved 3;
  reserved "nights";
  string status = 4;
}
```

<!-- reset_layout -->

Reusing number 3 for a *different type* means old clients still
running against the old schema will misinterpret bytes on the wire --
often silently, not with a clean error.

<!-- end_slide -->

## Back to the opening scenario

That REST struct-tag typo (`json:"geust"`) compiled fine, shipped, and
failed silently. Nothing in JSON or `encoding/json` was ever going to
catch it before a human noticed.

<!-- pause -->

**A `.proto` contract with a mismatched field type is a Go compile
error.** You find out at `go build`, from your own machine, not from a
support ticket days later.

That's the resolve to the question we opened with: a real, generated,
compiler-checked schema -- not a test you remembered to write, not a
reviewer who happened to notice.

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Protobuf's binary wire format** is smaller and faster than JSON --
   field numbers, not names, on the wire, with no repeated key text
2. **A `.proto` file is a real, compiler-checked contract** -- type
   mismatches are Go compile errors, not silent runtime surprises
3. **One `.proto` generates code for many languages** -- a shared
   contract instead of each team hand-maintaining its own JSON logic
4. **Field numbers are permanent, backward-compatible contracts** --
   `reserved` protects old clients when you evolve a schema
5. **gRPC adds four call shapes on top**, including streaming, which
   has no clean REST equivalent
6. **`context.Context` first, always** -- a Go-wide idiom, not a
   gRPC-specific one
7. **Never hand-edit `.pb.go` files** -- the `.proto` is the source of
   truth, regeneration silently overwrites direct edits
8. **A gRPC method tests like any other Go method** -- because
   `context.Context` comes first, you call it directly in a test with
   zero network setup, no `grpc.Dial` required

<!-- end_slide -->

## Bridge to Topic 12

**We've established:**

<!-- incremental_lists: true -->

- Protobuf trades JSON's flexibility for a real, compiler-checked
  contract and a smaller, faster wire format
- gRPC adds typed streaming that REST has no clean equivalent for
- Generated code is a build artifact -- the `.proto` is the source of
  truth, never the other way around
- You already proved `GetBooking` with a direct function-call test in
  this lab's Exercise 7 -- no network, no server, no `grpc.Dial`

<!-- incremental_lists: false -->

**Topic 12: Testing** -- Go's built-in testing toolchain, `_test.go`
files, table-driven tests, and how to test the services you've been
building across the last two topics, REST and gRPC alike.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
