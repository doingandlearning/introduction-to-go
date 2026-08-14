# Topic 5 sample code — Inheritance and Interfaces

Three small programs, one module. Run everything from this directory
(`code/`).

## `cmd/embeddemo`

Struct embedding: `Employee` embeds `Person` with no field name, which
promotes `Person`'s field and method onto `Employee`. Also shows
`Employee` shadowing `Greet()` and calling through to the embedded
version explicitly - the manual "super call."

```
go run ./cmd/embeddemo
```

Try deleting the `Role` field from `Employee` and re-running - the
promotion behavior doesn't change, because it has nothing to do with
what else is on the struct.

## `cmd/interfacedemo`

Three unrelated types - `Person`, `Robot`, `Parrot` - each independently
satisfy a one-method `Greeter` interface with no declared relationship
to it or to each other. Looped over in a `[]Greeter` and called
polymorphically. Also includes an `any` + type switch example.

```
go run ./cmd/interfacedemo
```

Search this file for the word "implements." It isn't there - that's the
whole point.

## `cmd/nilinterfacedemo`

Reproduces the classic nil-pointer-in-an-interface gotcha: a nil
`*MyError` returned as an `error` compares `!= nil`, because the
interface value is a `(type, value)` pair and the type half is non-nil.
Then shows the fix - return a literal `nil` in the success path instead
of a variable declared as a concrete pointer type.

```
go run ./cmd/nilinterfacedemo
```

Expected output:

```
leaky():  err == nil -> false
fixed():  err == nil -> true
```

## Formatting

```
gofmt -l .    # lists files that don't match gofmt's formatting
gofmt -w .    # rewrites them in place
```
