# Topic 4 sample code — Object-Oriented Programming

Three small programs, one module. Run everything from this directory
(`code/`).

## `cmd/receiverdemo`

A struct (`Item`) plus two methods — one value receiver, one pointer
receiver. Shows that a method is just a function with a receiver, and
that the receiver kind decides whether a mutation sticks.

```
go run ./cmd/receiverdemo
```

## `cmd/itemctor`

Go has no constructor syntax. `NewItem` is the conventional
"validate and build" function — a plain function returning `(*Item,
error)`. Shows both the success path and the rejected-input path, then
prints a zero-value `Item` to show it's usable without ever calling
`NewItem` at all.

```
go run ./cmd/itemctor
```

## `cmd/mapgotcha`

The map-of-structs gotcha: `map[string]Item` values aren't addressable,
so `catalog["widget"].ApplyDiscount(10)` won't compile.

```
go run ./cmd/mapgotcha           # runs the working read-mutate-write-back fix
```

To see the compiler reject the direct call, uncomment the line marked
"DOES NOT COMPILE" in `cmd/mapgotcha/main.go` and run:

```
go build ./cmd/mapgotcha
```

Read the error, then re-comment the line and re-run `go run
./cmd/mapgotcha` to see the fix that follows it.
