# Topic 2 sample code — Core Language Features

Four small programs, one module. Run everything from this directory
(`code/`).

## `cmd/errdemo`

No exceptions, no try/catch — a function that can fail returns `error`
as its last value, and the caller checks it immediately. Shows the happy
path, the failure path, and wrapping an error with `fmt.Errorf`'s `%w`.

```
go run ./cmd/errdemo
```

This one also ships its first test, `errdemo_test.go` — the "your first
Go test" example from the lecture. Run it, then break `divide` on
purpose (return the wrong value, or drop the zero check) and watch the
failure message change:

```
go test ./cmd/errdemo/...
```

## `cmd/zerodemo`

Declares a variable of several basic types — `int`, `float64`, `string`,
`bool`, a pointer, a slice, a map, and a struct — without assigning any
of them, then prints each. Compare the output against what `None` in
Python or `null` in Java would give you: every zero value here is
deterministic, and the zero-value `Point` is immediately usable with no
constructor.

```
go run ./cmd/zerodemo
```

## `cmd/deferdemo`

Three functions, three angles on `defer`:

- `orderedDefers` — three deferred calls, printed in LIFO order
- `loopGotcha` — the classic bug: `defer fmt.Println(i)` inside a loop
- `loopFixed` — the same loop, fixed by passing `i` as an explicit
  parameter to the deferred closure

```
go run ./cmd/deferdemo
```

Predict the output of `loopGotcha` before running it — most people guess
`0, 1, 2` and get `2, 1, 0` instead.

## `cmd/fmtdemo`

The `fmt` primer as a runnable program. Walks through `Print`/`Println`/
`Printf`, every verb from the lecture (`%v`, `%+v`, `%#v`, `%T`, `%d`,
`%s`, `%q`, `%f`, `%.2f`, `%t`, `%p`), `Sprintf`, `Errorf` with `%w`,
`Fprintf`/`Fprintln` against `os.Stderr` and an in-memory `bytes.Buffer`,
and finishes with a `Stringer` implementation showing how one method
takes over `%v` and `Println` output for a type.

```
go run ./cmd/fmtdemo
```

Read the output next to the source — every line is labeled with the verb
that produced it, so you can match theory to output directly.

## Formatting

```
gofmt -l .    # lists files that don't match gofmt's formatting
gofmt -w .    # rewrites them in place
```
