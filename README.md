# Go Programming

A four-day, hands-on introduction to Go, from writing your first program to shipping a concurrent, containerized service. Lectures and labs alternate through 13 topics, closing with a capstone that combines everything into one realistic tool.

## Course overview

This course offers a complete introduction to Go, covering the basics of writing and running Go programs through advanced features like concurrency, design patterns, and containerized deployment. You'll explore Go's core language constructs, data structures, and its distinctive approach to object-oriented and functional programming. Key modules include structuring services with the handler/service/repository pattern and dependency injection, building RESTful services, implementing gRPC with Protocol Buffers, and applying essential design patterns like Singleton and Strategy. The course also covers effective testing, benchmarking, and debugging practices, plus containerizing and running Go services locally with Docker and Docker Compose — equipping you to develop, test, and ship scalable, high-performance applications in Go.

## Course details

- **Duration:** 4 days (approx. 24 hours of lecture and labs)
- **Go version:** 1.26
- **Prerequisites:**
  - Familiarity with core programming concepts — variables, control flow, functions, and basic data structures (arrays, lists, maps) — ideally in a language like Python, Java, or C.
  - Basic knowledge of OOP concepts such as classes, inheritance, and interfaces, which helps in understanding Go's approach to structuring code.
  - No prior Docker experience is required; basic command-line familiarity is helpful for the Docker & Deployment module.

## Learning outcomes

By the end of the course you will be able to:

- Write, run, and organize Go programs using core language features like variables, functions, and packages
- Use loops, decision-making, and data types such as structs, arrays, and maps effectively
- Work with Go's inheritance, interfaces, and functional programming techniques
- Use goroutines and channels to create efficient, parallelized programs
- Structure services with the handler/service/repository pattern and apply dependency injection for testable, loosely-coupled code
- Use patterns like Singleton, Strategy, and Factory for structured, scalable code, and evaluate the trade-offs of each
- Create REST APIs and use gRPC with Protocol Buffers for high-performance services, understanding the benefits Protocol Buffers offer over JSON
- Write unit tests, measure coverage, and benchmark performance to ensure quality and efficiency
- Containerize Go services with multi-stage Docker builds and run them locally with Docker Compose, including guidance for Docker Desktop on WSL

## Repository layout

The course materials live in three parallel sets of folders, numbered 01–13 to match the topics (plus a capstone):

```
slides/   Slide decks for each topic (one slides.md per topic)
code/     Sample code from the lectures (one Go module per topic)
labs/     Hands-on exercises (one lab per topic, plus a capstone)
```

### `slides/`

A `slides.md` file per topic, written in Presenterm markdown.

### `code/`

Standalone Go modules (Go 1.26) demonstrating what each lecture covers. Every module has its own `README.md` explaining what's inside and how to run it — typically with `go run ./...` from that topic's directory.

### `labs/`

One lab per topic. Each lab is a self-contained directory with:

```
exercise.md   The lab brief — read this first
starter/      Starting code with TODOs for you to fill in
solution/     A complete reference implementation — don't look until you've had a go
```

Each topic folder may also contain extra materials where relevant (for example `standalone/` in the Concurrency lab).

## How the labs work

Most exercises ship with their test **already written** in `starter/` — your job is to make it pass. The pattern is deliberate:

1. Run `go test ./... -v` from `starter/` **before touching anything**. Watch it fail.
2. Read the failure — it tells you exactly what the exercise is asking for.
3. Implement the TODO, following the exercise brief.
4. Re-run `go test ./... -v` and confirm the test is green.

A few exercises (the ones in Topic 12 and part of the capstone) deliberately ask you to **write the test yourself** — that's the point of the Testing topic, so don't skip ahead.

Command-line verification steps (`go run ...`) are still worth doing where the brief includes them — the tests check correctness, but running the real program is how you see it work.

## Tooling

To work through the course you'll need:

- **Go 1.26** — from [go.dev/dl](https://go.dev/dl). Verify with `go version`.
- **An editor** — the course uses VS Code, but any editor with Go support works.
- **Docker** (Topic 13 only) — Docker Desktop, including guidance for Docker Desktop on WSL.

## Course outline

### Topic 1 — Go Essentials

- Getting Started with Go
- Writing and Running Go Programs
- Using an IDE

### Topic 2 — Core Language Features

- Getting Started
- Variables and Types
- Functions
- Formatted I/O — a Primer on the `fmt` Package
- Organizing Packages

### Topic 3 — Flow Control & Data Structures

- Looping
- Decision Making
- Pointers
- Structs
- Arrays
- Slices
- Ranges
- Maps

### Topic 4 — Object-Oriented Programming

- Introduction to Object-Oriented Programming
- Defining and Using Classes in Go
- Initialization

### Topic 5 — Inheritance and Interfaces

- Getting Started with Inheritance
- Inheritance in Go
- Interfaces
- Polymorphism

### Topic 6 — Functional Programming

- Functional Programming
- Higher Order Functions
- Additional Techniques

### Topic 7 — Concurrency

- Deferred Function Execution
- Getting Started with Concurrency
- Channels
- Channel Techniques
- Buffered Channels

### Topic 8 — Introduction to Design Patterns

- Introduction to Design Patterns and Evaluating Trade-offs
- Example Design Pattern: Singleton
- Implementing a Thread-Safe Singleton
- Leveraging `sync.Once`
- The Handler/Service/Repository Pattern
- Dependency Injection Fundamentals

### Topic 9 — Additional Design Patterns

- Strategy Pattern
- Abstract Factory Pattern
- Decorator Pattern
- Builder Pattern
- Dependency Injection in Practice

### Topic 10 — REST Services

- Overview of REST Services
- Creating a REST Service in Go
- Structuring Services with Handler/Service/Repository
- Additional Techniques
- Implementing a Full CRUD REST Service
- Running Services Locally with Docker Compose

### Topic 11 — gRPC and Protocol Buffers

- Overview of gRPC and Protocol Buffers
- The Benefits of Protocol Buffers
- Getting gRPC and Protocol Buffers
- Simple Example of gRPC Services
- A Closer Look at Service Methods

### Topic 12 — Testing

- Setting the Scene
- Getting Started with Go Tests
- Testing for Panics
- Parameterized Tests
- Test Coverage
- Benchmarking

### Topic 13 — Docker & Deployment

- Containerizing Go Applications
- Multi-Stage Docker Builds
- `scratch` vs. `distroless` Base Images
- Docker Compose for Local Development
- Working with Docker Desktop on WSL
- An Introduction to Deployment Options

### Capstone — Lab 14: a concurrent URL-status spider

The course ends with a capstone that combines four things into one small, realistic tool: real file I/O, a real HTTP client, CLI flags, and the worker-pool concurrency pattern from Topic 7. It reads a list of URLs from a file, checks every one concurrently, and writes a CSV report of what happened to each.
