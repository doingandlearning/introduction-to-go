# Topic 12 sample code — Testing

One module, one package (`internal/billing`), one demo command. Run
everything from this directory (`code/`).

## `internal/billing`

Reuses the discount/strategy shape from Topic 9 (`TenPercentOff`,
`DiscountStrategy`, `ApplyDiscount`) plus two new functions written
specifically to be tested:

- `TierDiscount` — four branches based on price, good for a
  table-driven test.
- `Divide` — panics instead of returning an error when dividing by
  zero, good for a defer/recover test.

`billing_test.go` demonstrates every concept from the lecture against
this one package:

| Concept | Test |
|---|---|
| Basic test | `TestTenPercentOff`, `TestDivide` |
| Panic test (defer/recover) | `TestDivideByZeroPanics` |
| Table-driven test with `t.Run` | `TestTierDiscount` |
| Testing a function passed as a value | `TestApplyDiscount` |
| Benchmark | `BenchmarkTenPercentOff` |

## Commands

```
go test ./...                                  # run all tests
go test -v ./...                                # run all tests, verbose per-test output

go test -run TestTierDiscount ./...             # run one test function
go test -run TestTierDiscount/boundary ./...    # run one table-driven case by name

go test -cover ./...                            # coverage percentage
go test -coverprofile=cover.out ./...           # write a coverage profile
go tool cover -html=cover.out                   # open the line-by-line HTML report

go test -bench=. ./...                          # run benchmarks
go test -bench=. -benchmem ./...                # ...with allocation counts
```

## `cmd/billingdemo`

Not part of the testing story — just something to `go run` so the
package's behaviour is visible without reading test output:

```
go run ./cmd/billingdemo
```

## A note on `t.Errorf` vs `t.Fatalf`

`Error`/`Errorf` mark the test failed but let the rest of the test
function keep running. `Fatal`/`Fatalf` mark it failed and stop
immediately — useful when a later line would panic or make no sense
given the earlier failure (e.g. a nil pointer you were about to
dereference). Every test above uses `Errorf`; none of them need to
stop early.
