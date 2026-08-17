---
title: "**Additional Design Patterns**"
sub_title: Go Programming — Topic 9
author: Kevin Cunningham
---

## Opening scenario

You need to swap your production database for a fake one in 200 tests —
without touching a single line of the service code that talks to it.

In Java, that means a DI framework: Spring, an `@Autowired` field, a test
`@Configuration` class, maybe a mocking library.

**Type in chat: where would you even start in Go? Is there a Go
equivalent of `@Autowired`, or is the honest answer "there isn't one"?**

We'll come back to this once we've seen how Go wires dependencies by
hand.

<!--
speaker_note: |
  Let guesses land - some delegates will say "there must be a DI
  library," others will already suspect it's just constructor arguments
  from Topic 8's preview. Both reactions are useful, don't resolve yet.
-->

<!-- end_slide -->

## Where we left off in Topic 8

Topic 8 introduced Singleton, a pros/cons framework for evaluating any
pattern, and a first pass at dependency injection: one service, one
constructor argument, a `Handler -> Service -> Repository` layering
preview.

<!-- pause -->

Today: four more classic patterns, evaluated the same way — plus a much
harder dependency graph than "one service, one dependency."

**Same lens as Topic 8, every time:** testability, coupling,
discoverability, and "is there a simpler Go-native way to get this?"

<!-- end_slide -->

<!-- jump_to_middle -->

Strategy
===

<!-- end_slide -->

## Strategy: the Java shape

In Java, an interchangeable algorithm means an interface plus one class
per implementation — even for something as small as "how do I calculate
a discount."

```java
interface DiscountStrategy {
    double apply(double price);
}
class TenPercentOff implements DiscountStrategy {
    public double apply(double price) { return price * 0.9; }
}
```

<!-- pause -->

A class, just to hold one method. Go has a shorter path.

<!-- end_slide -->

## Strategy: the Go shape

A function type. No interface, no classes, no boilerplate constructor.

```go
type DiscountStrategy func(price float64) float64

func NoDiscount(price float64) float64    { return price }
func TenPercentOff(price float64) float64 { return price * 0.9 }

func ApplyDiscount(price float64, strategy DiscountStrategy) float64 {
    return strategy(price)
}

ApplyDiscount(100, TenPercentOff) // 90
```

<!-- pause -->

**`ApplyDiscount` doesn't know or care what's inside `strategy`** — only
that it has the right shape. That's the whole pattern, in four lines.

<!--
speaker_note: |
  This ordering is deliberate - function type first, interface-based
  version later. Showing the lightweight version before the Java-shaped
  one makes the contrast land harder than doing it the other way round.
-->

<!-- end_slide -->

## Strategy is a sticky note, not a rulebook

Picture handing a cashier a sticky note: "if it's a member, take 15%
off." That's the whole strategy — no binder, no procedure manual, no
training session.

<!-- incremental_lists: true -->

- The cashier doesn't need to know why the discount exists
- Swapping the note swaps the behavior completely
- A whole rulebook binder only makes sense once the rule itself needs
  multiple steps or its own memory

<!-- incremental_lists: false -->

A Java `Strategy` interface is the rulebook binder. A Go function value
is the sticky note.

<!-- end_slide -->

## Strategy: pros and cons

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->
# Pros

- Minimal ceremony — a function type + functions, no class hierarchy.
- Trivial to test: call the function directly, no mock framework.
- Composes directly with functional options from Topic 6. 
<!-- column: 1 -->

<!-- pause -->

# Cons


- If the "strategy" genuinely needs more than one method, or state a closure can't hold, a bare func type stops being enough — you're back to an interface + structs anyway. 


<!-- pause -->

<!-- reset_layout -->

**In breakout rooms (5 minutes):** Take the `DiscountStrategy` idea and
add a requirement — the strategy now needs to log every application to
an audit trail with its own retry counter. Does a function type still
work? What does the interface version look like once it does?

<!--
speaker_note: |
  The expected answer: once the strategy needs internal state that
  survives across calls (a retry counter) plus a second method (flush
  the audit log), a closure captures state fine but can't expose a
  second method - that's the real boundary where you go back to an
  interface + struct.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Abstract Factory
===

<!-- end_slide -->

## Abstract Factory: families of related objects

Sometimes one factory method isn't the shape you need — you need a whole
*coordinated bundle* of related objects, all consistent with each other.

```go
type ConnectorFactory interface {
    NewConnection() Connection
    NewQueryBuilder() QueryBuilder
}

type MySQLFactory struct{}
func (MySQLFactory) NewConnection() Connection     { return &mysqlConn{} }
func (MySQLFactory) NewQueryBuilder() QueryBuilder { return &mysqlQueryBuilder{} }

type PostgresFactory struct{}
func (PostgresFactory) NewConnection() Connection     { return &pgConn{} }
func (PostgresFactory) NewQueryBuilder() QueryBuilder { return &pgQueryBuilder{} }
```

<!-- pause -->

Swap `MySQLFactory{}` for `PostgresFactory{}` at the composition root,
and every connection and query builder downstream changes together.

<!-- end_slide -->

## A furniture collection, not five mismatched stores

Abstract Factory is buying an entire coordinated furniture collection
from one catalog — sofa, table, and lamp all in the same finish —
instead of picking pieces from five different stores that happen to
each sell "a sofa."

<!-- pause -->

Java reaches for this constantly, partly to work around verbose
`new ConcreteThing()` calls everywhere and the lack of first-class
functions. Go's lighter interfaces and lack of constructors soften both
of those reasons — the pattern is real, but it's reached for less often.

<!-- end_slide -->

## Abstract Factory: pros and cons

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->
# Pros

- Swapping an entire coordinated family (all-MySQL vs all-Postgres) is one line at the composition root.
- Callers never touch a concrete type. 
<!-- pause -->

<!-- column: 1 -->
# Cons


- Genuinely more ceremony than Go culture usually wants for "just pick one DB driver." 
- Easy to over-apply — if you only ever have one real implementation, it's premature abstraction. 

<!-- pause -->

<!-- reset_layout -->


**Type in chat: using Topic 8's rule — "extract an interface when you
actually have two implementations" — would you reach for
`ConnectorFactory` on day one of a new project, or wait?**

<!--
speaker_note: |
  Steer toward "wait." Most teams start with one database and never
  need the second family - the interface (and definitely the abstract
  factory) should show up when the second real implementation does, not
  speculatively.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Decorator
===

<!-- end_slide -->

## Decorator: wrap something, keep the same interface

Decorator wraps a thing to add behavior, while looking exactly like the
thing it wraps to everyone downstream.

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("request:", r.URL.Path)
        next.ServeHTTP(w, r)
        log.Println("done")
    })
}

handler := LoggingMiddleware(AuthMiddleware(baseHandler))
```

<!-- pause -->

**This is not a textbook example you'll never use.** This exact shape is
the load-bearing idiom behind every piece of Go HTTP middleware you will
ever write — including Topic 10's REST service, next.

<!-- end_slide -->

## Airport security, not a design pattern

Each security checkpoint wraps the "get to your gate" journey. You pass
through one, then the next, then the next.

<!-- incremental_lists: true -->

- The gate agent at the end doesn't know or care how many checkpoints
  you passed through
- Each checkpoint can be added, removed, or reordered independently
- Stand in a different queue order, and the trip changes — same
  checkpoints, different experience

<!-- incremental_lists: false -->

That last point is the part worth sitting with. Ordering isn't
cosmetic.

<!-- end_slide -->

## Ordering is a real decision, not an implementation detail

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```go
// auth first
handler := AuthMiddleware(
    LoggingMiddleware(base),
)
```

A request that fails auth
**never gets logged.**

<!-- column: 1 -->

```go
// logging first
handler := LoggingMiddleware(
    AuthMiddleware(base),
)
```

A request that fails auth
**gets logged anyway.**

<!-- reset_layout -->

**In breakout rooms (5 minutes):** Which ordering would your team want
for a production API, and why might the "wrong" answer actually be
correct for a security-auditing use case?

<!--
speaker_note: |
  There's no universal right answer - security/audit teams often want
  every attempt logged, including failed auth, which argues for logging
  outermost. A team optimizing for clean logs might want the opposite.
  Let the room actually disagree.
-->

<!-- end_slide -->

## Decorator: pros and cons

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

# Pros

- Composable — stack as many as you want, in any order, each independently testable and reusable across handlers. 
- Falls directly out of interfaces + first-class functions, zero extra language machinery. 
<!-- pause -->

<!-- column: 1 -->
# Cons

- Ordering matters and isn't always obvious from the call site. 
- A long chain gets hard to trace — "which of these six wrappers actually set this header?" 

<!-- pause -->

<!-- reset_layout -->


**Demo:** chain `LoggingMiddleware` and a header-setting middleware
around one handler two different ways, and compare what each ordering
actually produces.

<!--
speaker_note: |
  Run the code/cmd/decorator demo live here, then swap the nesting order
  and rerun so the room sees, not just hears, that ordering is
  observable. Good moment to preview that Topic 10's whole middleware
  stack is built from exactly this.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Builder
===

<!-- end_slide -->

## Builder: step-by-step assembly with validation

```go
type QueryBuilder struct { cols, table, where string }

func NewQueryBuilder() *QueryBuilder { return &QueryBuilder{} }
func (q *QueryBuilder) Select(cols string) *QueryBuilder { q.cols = cols; return q }
func (q *QueryBuilder) From(table string) *QueryBuilder  { q.table = table; return q }
func (q *QueryBuilder) Where(cond string) *QueryBuilder  { q.where = cond; return q }

func (q *QueryBuilder) Build() (string, error) {
    if q.cols == "" || q.table == "" {
        return "", fmt.Errorf("select and from are required")
    }
    return fmt.Sprintf("SELECT %s FROM %s WHERE %s", q.cols, q.table, q.where), nil
}
```

<!-- pause -->

```go
sql, err := NewQueryBuilder().Select("id, name").From("guests").Where("active = true").Build()
```

<!-- end_slide -->

## The deli counter, not the coffee counter

A deli counter insists on an order: bread, then protein, then toppings,
then sauce. Try to add sauce before there's bread underneath it and the
person behind the counter stops you.

<!-- pause -->

Contrast Topic 6's coffee shop modifier slips — oat milk, extra shot,
decaf, in **any order**, because none of them depends on another.
`Build()` is the deli counter's insistence made structural: it can
refuse to hand you a sandwich until the required layers exist.

<!-- end_slide -->

## Builder vs. functional options — a real judgment call

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Builder wins when:**

Steps have a required order, or some steps are mandatory before
`Build()` can succeed — validation is a natural, structural part of the
chain.

<!-- column: 1 -->

**Functional options win when:**

Options are genuinely independent and order truly doesn't matter — no
natural hook exists for "you must call `WithX` before `WithY`."

<!-- reset_layout -->

**Type in chat: for a `Server` with 6 independent settings (Topic 6) vs.
a `QueryBuilder` where `Select` must precede `Build`, which pattern
fits which — and is there a case where you're honestly not sure?**

<!--
speaker_note: |
  The honest answer is "it depends," and that's fine to say plainly -
  don't force a false resolution here. If nobody raises an ambiguous
  case, offer one: a builder where two of five steps are optional and
  three are required is a genuine gray area.
-->

<!-- end_slide -->

## Builder: pros and cons

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->
# Pros

- Natural place for step-by-step validation and required ordering. 
- Reads clearly at the call site — chained, self-documenting. 


<!-- pause -->

<!-- column: 1 -->
# Cons

- More boilerplate than functional options for the common "some independent optional knobs" case. 
- Picking the wrong one for the job is a real, common mistake. 
<!-- pause -->

<!-- reset_layout -->

Exercise 5 in the lab makes you build the same thing both ways and feel
the seam yourself — the "it depends" isn't a cop-out, it's the actual
skill.

<!-- end_slide -->

<!-- jump_to_middle -->

Dependency Injection, at real scale
===

<!-- end_slide -->

## Back to the opening scenario

You asked: is there a Go `@Autowired`? **No — and that's a deliberate
cultural choice, not a missing feature.**

<!-- pause -->

Topic 8 showed constructor injection with one dependency. Real systems
have graphs, not single edges. Let's wire one.

<!-- end_slide -->

## A dependency graph, not a dependency

`OrderHandler` depends on `OrderService`. `OrderService` depends on
**two** interfaces — `Repository` and `Notifier` — not one.

```go
func main() {
    repo := NewInMemoryRepository()
    notifier := NewConsoleNotifier()
    svc := NewOrderService(repo, notifier)
    handler := NewOrderHandler(svc)
    // register handler with an http.ServeMux — previews Topic 10
}
```

<!-- pause -->

**Read this top to bottom and you have the entire object graph of the
program.** Nothing is hidden, nothing is resolved for you at a time you
don't control.

<!-- end_slide -->

## Spring vs. `main` — the same job, two philosophies

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Spring**

```java
@Service
class OrderService {
    @Autowired
    Repository repo;
    @Autowired
    Notifier notifier;
}
```

Resolved implicitly, at runtime,
by a container doing component
scanning.

<!-- column: 1 -->

**Go**

```go
svc := NewOrderService(
    repo,
    notifier,
)
```

Resolved explicitly, by you,
reading top to bottom in `main`.

<!-- reset_layout -->

Neither is "wrong." Go's bet is that explicitness at the wiring point is
worth more than the convenience of not writing it.

<!-- end_slide -->

## This only works because of interfaces

`Handler -> Service -> Repository` from Topic 8, Strategy, Abstract
Factory, and manual DI in this topic are **the same underlying idea**
wearing four different outfits: depend on behavior, not on a concrete
implementation.

<!-- incremental_lists: true -->

- Strategy: depend on a function shape, not a specific algorithm
- Abstract Factory: depend on a factory interface, not a concrete family
- Decorator: depend on the same interface the thing you're wrapping has
- DI: depend on `Repository`/`Notifier` interfaces, not
  `InMemoryRepository`/`ConsoleNotifier`

<!-- incremental_lists: false -->

Swap any concrete type for a test double, and nothing above it in the
call chain notices.

<!-- end_slide -->

## When hand-wiring gets annoying

At real scale — dozens of services, a graph several layers deep —
constructing everything by hand in `main` gets verbose, and reordering
it by hand is error-prone.

<!-- pause -->

The standard answer, if this pain is real for your team: **`google/wire`**.

<!-- incremental_lists: true -->

- Compile-time code generation, not a runtime container
- No reflection — it still produces plain `NewX(...)` calls, just
  generated for you instead of typed by hand
- Not part of the standard toolchain, and most small-to-medium Go
  services never need it

<!-- incremental_lists: false -->

One thing to know exists, not a tool to master today.

<!--
speaker_note: |
  Keep this to the single slide as scoped - it's a name-drop for teams
  that hit real pain, not a deep dive. If asked for more, the short
  answer is "it reads your constructor signatures and generates the
  main-style wiring code you'd have written by hand."
-->

<!-- end_slide -->

## Manual DI: pros and cons

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->
# Pros

- Fully explicit dependency graph — readable top to bottom in `main`, no "magic" injection. 
- Trivial to substitute test doubles — pass a fake `Repository`/`Notifier` in a test. Zero framework dependency or startup cost. 

<!-- pause -->

<!-- column: 1 -->
# Cons

- Genuinely more typing than a DI container for large graphs. 
- Easy for the composition root to become one large, unwieldy `main` if not organized. 

<!-- pause -->

<!-- reset_layout -->



<!-- pause -->

**Mitigation, not a different tool:** split wiring into a dedicated
`wire.go` or `setup.go` file once `main` gets crowded — still 100% hand-
written, just organized. That's a different move from actually adopting
the `wire` tool.

**In breakout rooms (5 minutes):** At what point — how many services,
how many layers — would your team stop hand-wiring `main` and reach for
`google/wire`? Defend a specific number.

<!--
speaker_note: |
  There's no correct number - the point is forcing a concrete answer
  instead of "it depends" as a dodge. Push back if a room says "never"
  or "immediately" without reasoning.
-->

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Strategy**: a function type usually beats an interface — reach for
   the interface only once state or multiple methods are needed
2. **Abstract Factory**: earns its place for real swappable families,
   easy to over-apply for "just pick one implementation"
3. **Decorator**: the actual mechanism behind Go HTTP middleware —
   ordering is a real decision, not cosmetic
4. **Builder vs. functional options**: required ordering and validation
   favor Builder; independent optional settings favor functional options
5. **Manual DI**: explicit, hand-wired constructor calls at the
   composition root — `google/wire` exists for when that gets painful at
   scale
6. Exercise 7 puts your `_test.go` habit from Topic 2 to work on today's
   patterns: a Strategy function tests for free, a Builder's validation
   is exactly what you want CI to catch, and a two-dependency DI graph
   is still trivial to fake

<!-- end_slide -->

## Back to the opening scenario

You asked where to even start, swapping a real database for a fake one
across 200 tests, with no framework in Go.

**You start exactly where Exercise 6 puts you:** `Repository` and
`Notifier` are interfaces. Every test constructs `NewOrderService` with
fakes instead of real implementations. `OrderService` itself never
changes.

<!-- pause -->

**Type in chat: now that you've seen the wiring, does "no `@Autowired`"
feel like a missing feature, or a design choice you'd defend?**

<!--
speaker_note: |
  Most delegates land on "design choice" once they've seen how little
  code manual DI actually costs for realistic graph sizes - let a few
  dissent if the graph-at-scale argument resonated with them instead.
-->

<!-- end_slide -->

## Bridge to Topic 10

**We've established:**

<!-- incremental_lists: true -->

- Strategy, Abstract Factory, Decorator, and Builder all trade Java's
  class-heavy ceremony for something lighter in Go — with a real
  boundary where the lighter version stops being enough
- Decorator, specifically, is the exact mechanism behind HTTP middleware
- Manual DI scales from one dependency to a real graph without a
  framework, and `google/wire` exists for when it doesn't
- Every pattern today paid for itself in Exercise 7 too — Strategy,
  Builder, and DI all stayed just as easy to unit test as they were to
  design

<!-- incremental_lists: false -->

**Topic 10: REST Services** — we'll build the real
`Handler -> Service -> Repository` stack, wire it with everything from
today, and wrap it in the exact middleware chain from the Decorator
section.

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===
