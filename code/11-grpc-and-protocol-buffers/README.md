# Topic 11 sample code -- gRPC and Protocol Buffers

## What's real in this bundle, and what isn't

This sandbox has no `protoc`, no Go toolchain, and no internet access, so
this bundle is honest about the split:

| File | Status |
|---|---|
| `booking.proto` | **Real.** Plain-text protobuf, valid, no toolchain needed to write or read it. This is the actual source-of-truth teaching artifact. |
| `bookingpb/booking.go` | **Illustrative, hand-written.** A simplified approximation of what `protoc --go_out --go-grpc_out` would generate. It is NOT generator output -- see the warning comment at the top of the file for exactly what's missing (proto.Message plumbing, wire encoding, the real grpc-go service descriptor). |
| `cmd/server/main.go` | **Illustrative, hand-written.** Shows the shape of a real server implementation (the mutex-protected in-memory store is real, working Go, carried over from Topic 10). The actual network wiring (`net.Listen`, `grpc.NewServer`, `grpcServer.Serve`) is commented out, not executed, because `google.golang.org/grpc` isn't available here. |
| `cmd/client/main.go` | **Illustrative, hand-written.** Shows the call shapes (`client.GetBooking(ctx, req)`, stream `Recv()` until `io.EOF`) exactly as they'd look against a real generated client. The "connection" underneath is a fake in-process stand-in, not a real `grpc.ClientConn`. |
| `go.mod` | Present for directory shape; the `google.golang.org/grpc` and `google.golang.org/protobuf` dependencies a real version needs are commented out, not fetched. |

**Before using any of this in class or in a real project: delete
`bookingpb/`, `cmd/server/`, and `cmd/client/` and regenerate/rewrite them
against a real toolchain.** The `.proto` file is the one artifact here you
can trust and reuse directly.

## Setting up the real toolchain (in delegates' own environment)

1. Install `protoc`, the Protocol Buffer compiler itself. It's not a Go
   tool -- it's a standalone C++ binary. On macOS: `brew install protobuf`.
   On Linux, grab a release from the `protocolbuffers/protobuf` GitHub
   releases page, or use your package manager. Confirm with `protoc --version`
   (aim for a reasonably recent 3.x or later).

2. Install the two Go plugins `protoc` needs to know how to emit Go code
   and gRPC service code:

   ```
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

   These install into `$HOME/go/bin` -- make sure that's on your `PATH`,
   or `protoc` won't find them.

3. Generate real code from `booking.proto`:

   ```
   protoc --go_out=. --go_opt=paths=source_relative \
          --go-grpc_out=. --go-grpc_opt=paths=source_relative \
          booking.proto
   ```

   This produces `booking.pb.go` (message types) and
   `booking_grpc.pb.go` (client/server interfaces) -- replace the
   illustrative `bookingpb/booking.go` in this bundle with those.

4. Add the runtime dependencies to `go.mod`:

   ```
   go get google.golang.org/grpc
   go get google.golang.org/protobuf
   ```

5. Rewrite `cmd/server/main.go` and `cmd/client/main.go` against the real
   generated types -- the commented-out blocks in both files show you
   exactly what that wiring looks like; uncomment and adapt them.

**Version mismatches between `protoc` and the Go plugins are the most
common source of confusing generation errors.** If `protoc` produces
something that doesn't compile, check plugin versions against the
`protoc` version before assuming the `.proto` file itself is wrong.

## The domain

`Booking` here is deliberately the same shape as Topic 10's REST
`Booking` struct (`id`, `guest`, `nights`) so the two topics are directly
comparable -- see lab exercise 6.

## Never hand-edit generated files

Once real `.pb.go` files exist, treat them exactly like a compiled
`.class` file or a minified JS bundle: readable, but not something you
edit directly. The `.proto` file is the source of truth. Regenerating
from it silently overwrites any direct edits made to the generated
output.
