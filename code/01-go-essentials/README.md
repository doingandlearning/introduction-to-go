# Topic 1 sample code — Go Essentials

Three small programs, one module. Run everything from this directory
(`code/01-go-essentials/`).

## `cmd/hello`

The smallest useful Go program. Demonstrates the four ways to turn source
into a running program:

```
go run ./cmd/hello              # compile + run, throws away the binary
go build ./cmd/hello && ./hello  # produces a standalone binary
go install ./cmd/hello           # builds, places binary in $HOME/go/bin
```

## `cmd/greetdemo` + `internal/greeting`

Two packages, one import. Shows Go's entire visibility mechanism:
capitalize a name to export it, lowercase it to keep it package-private.

```
go run ./cmd/greetdemo
```

Try the experiment in the comment at the top of `cmd/greetdemo/main.go` —
rename `greetingPrefix` to `GreetingPrefix`, call it from `main.go`, then
revert the rename but not the call. Watch the build fail.

## `cmd/vetdemo`

Compiles cleanly. Has a real bug the compiler won't catch.

```
go build ./cmd/vetdemo   # succeeds
go vet ./cmd/vetdemo     # flags the Printf/argument count mismatch
```

## Formatting

```
gofmt -l .    # lists files that don't match gofmt's formatting
gofmt -w .    # rewrites them in place
```

Deliberately mis-indent something in any file here first, then run both —
the diff is the whole point.
