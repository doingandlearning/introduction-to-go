---
title: "**Docker & Deployment**"
sub_title: Go Programming — Topic 13
author: Kevin Cunningham
---

## Opening scenario

A teammate on a Python service tells you their Docker image is **900MB**
and takes **3 minutes** to build. You mention your Go service's image is
**8MB** and builds in **20 seconds**.

They don't believe you. Same job — an HTTP API with a handful of
endpoints and a database dependency. Wildly different numbers.

**Type in chat: what do you think is actually different here? Guess before we get to the explanation.**

<!--
speaker_note: |
  Let a few guesses land - people often reach for "Go is just faster"
  or "Python is bloated," both true-ish but not the actual mechanism.
  Don't correct yet, just bank the guesses. The real answer is coming
  via the multi-stage build and static binary discussion.
-->

<!-- end_slide -->

## Why this topic exists

The client scoping call was explicit: **you use WSL and docker-compose
for local development, and you want real Docker coverage** — not a
five-minute aside.

Topic 10 never touched Docker at all, so there's no assumed prior
exposure to build on here — this topic starts from zero and covers all
of it:

<!-- incremental_lists: true -->

- Docker fundamentals, from a Go developer's angle specifically
- Multi-stage builds, properly explained
- WSL specifics, in real depth
- Deploying and containerizing what you've built across this course

<!-- incremental_lists: false -->

<!--
speaker_note: |
  Say this plainly rather than assuming it: nobody in the room has
  built or run a Docker image in this course before today, not even
  the "just enough to get a REST service up" version. Today is real
  first contact - pace the next section accordingly.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Docker, from zero
===

<!--
speaker_note: |
  This section exists because Topic 10 never touched Docker at all -
  don't assume ANY prior exposure walking in here, not even "I've seen
  a docker-compose.yml before." Treat this as a genuine first contact,
  not a recap.
-->

<!-- end_slide -->

## An image is a recipe, a container is what you get when you follow it

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Image**

A read-only bundle: a filesystem plus metadata about what to run.
Built once, with `docker build`. Doesn't change once built.

<!-- column: 1 -->

**Container**

A running (or stopped) instance of an image — the same relationship a
class has to an object. `docker run` starts one.

<!-- reset_layout -->

**One image, run many times, is many independent containers** — same
starting filesystem each time, no shared state between them unless you
explicitly wire it up (a volume, a network — both coming up shortly).

<!-- end_slide -->

## The Dockerfile: instructions for building the image

A `Dockerfile` is a plain text file, one instruction per line, read top
to bottom:

<!-- incremental_lists: true -->

- `FROM` — start from an existing image (an OS, a language runtime, or something minimal)
- `WORKDIR` — set the working directory inside the image for everything that follows
- `COPY` — copy files from your machine into the image
- `RUN` — execute a command *while building* the image (installing something, compiling something)
- `ENTRYPOINT` — the command that runs when a container starts from this image

<!-- incremental_lists: false -->

Each instruction adds a new layer on top of the last. `docker build`
runs them in order and produces the final image.

<!-- end_slide -->

## A first Dockerfile for a Go service

```dockerfile
FROM golang:1.22
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

Start from an image with the Go toolchain already installed, copy the
source in, compile it, run the result.

<!-- pause -->

**Demo:**

```
docker build -t docker-demo .
docker run -p 8080:8080 docker-demo
```

```
curl localhost:8080/healthz
curl localhost:8080/items
```

<!--
speaker_note: |
  Run this live end to end - build, run, curl - before saying anything
  about image size. Let it just work first. The size problem lands
  better once they've seen the happy path succeed, and it's exactly
  the Dockerfile this topic keeps coming back to over the next few
  slides ("the naive, single-stage Dockerfile") - it's not a throwaway
  example, name it as the same one again when you get there.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Why Go and Docker are an unusually good fit
===

<!-- end_slide -->

## Remember Topic 1?

You were sold on this: **Go compiles to one dependency-free binary. No
runtime needed on the target machine.**

<!-- pause -->

A Docker image is, at its core, a filesystem plus instructions for what
to run. If your program needs a runtime and a pile of installed packages
just to exist, all of that has to be *inside* the image.

**If your program needs nothing but itself, the image can be almost
nothing but the binary.**

<!--
speaker_note: |
  This is the direct callback to Topic 1's single-binary pitch - make
  the connection explicit rather than assuming they'll draw the line
  themselves. This slide is the thesis statement for the whole topic.
-->

<!-- end_slide -->

## What a Python or Node image actually contains

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**A typical Python/Node service image:**

- The full language runtime (CPython, or Node + V8)
- Every installed package (`pip install`, `npm install`) baked into a
  layer
- Often a package manager left in the final image, because removing it
  cleanly is its own project

<!-- column: 1 -->

**A typical Go service image:**

- The compiled binary
- Maybe CA certificates, if it makes HTTPS calls
- That's genuinely close to the whole list

<!-- reset_layout -->

**This is why 900MB vs 8MB is a real, typical gap — not a cherry-picked example.**

<!-- end_slide -->

## This has to be earned, not assumed

Go doesn't hand you a tiny image for free. The Dockerfile you just
built and ran — a naive, single-stage one — ships the entire Go
toolchain in the final image too: compiler, module cache, everything.

**Demo:** run `docker images` against the image you just built and look
at its size.

<!-- pause -->

**The mechanism that actually gets you the 8MB image is the multi-stage
build.** That's next.

<!-- end_slide -->

<!-- jump_to_middle -->

Multi-stage builds
===

<!-- end_slide -->

## The pattern

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

```dockerfile
# --- build stage ---
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /server ./cmd/server

# --- final stage ---
FROM gcr.io/distroless/static-debian12
COPY --from=builder /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

<!-- column: 1 -->

Two `FROM` lines, two stages.

The **builder** stage has the full Go toolchain — large, but it never
ships.

The **final** stage starts fresh from a tiny base and copies in only the
compiled binary via `COPY --from=builder`.

Only the final stage becomes your image.

<!-- reset_layout -->

**Demo:** build this, then run `docker images` and compare it against a single-stage version of the same Dockerfile.

<!--
speaker_note: |
  If you have Docker available to demo live, do it - the size number
  landing on screen is more persuasive than any slide. If not, walk
  through the Dockerfile line by line and have delegates predict the
  final image's contents before you reveal it.
-->

<!-- end_slide -->

## `CGO_ENABLED=0` — the gotcha worth naming

`CGO_ENABLED=0` disables cgo, which forces a **truly statically linked**
binary with no dependency on the system's C library (`libc`).

<!-- pause -->

**This is directly Topic 1's "no runtime needed" claim, made literal.**
If `CGO_ENABLED` is left at its default (`1` in some setups, especially
if any dependency pulls in cgo), your binary can silently pick up a
dynamic dependency on `libc` — and then it *won't* run standalone on a
minimal or `scratch`-based final image. It'll fail to start, with an
error that doesn't obviously point at cgo.

**Set it explicitly. Don't rely on the default.**

<!--
speaker_note: |
  This genuinely bites people - a binary that builds fine and runs
  fine on the build machine, then fails mysteriously in a minimal
  final-stage image because it silently linked against glibc. Worth
  a real pause here, not a rushed mention.
-->

<!-- end_slide -->

## `scratch` vs `distroless`

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**`FROM scratch`**

The absolute minimal base image — literally empty. Works, and is as
small as it gets.

But: no CA certificates, no shell, no `/etc/passwd`, nothing.

Any outbound HTTPS call fails until you also
`COPY --from=builder /etc/ssl/certs/ca-certificates.crt ...`

<!-- column: 1 -->

**`FROM gcr.io/distroless/static-debian12`**

Small, but not empty — includes a CA certificate bundle and basic
runtime files, still no shell and no package manager.

Slightly larger than `scratch`, meaningfully easier to work with on day
one.

<!-- reset_layout -->

**A real, commonly-hit gotcha:** ship `scratch` without the CA cert copy, then wonder why every outbound HTTPS call from that container fails with a certificate error.

<!--
speaker_note: |
  Worth demoing if possible - build the same service FROM scratch
  without the cert copy, hit an endpoint that makes an HTTPS call, and
  read the actual error together. It's not an obvious error at first
  glance.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

docker-compose, from zero
===

<!-- end_slide -->

## Why docker-compose exists

A real service rarely runs alone — it needs a database, maybe a cache,
maybe another internal service. Running each with its own `docker run`
command, remembering every flag, every time, doesn't scale past "one
container."

<!-- pause -->

`docker-compose` describes a whole stack — every service, its image or
build context, its ports, its environment — in one YAML file, and
brings the whole thing up or down with one command.

<!-- end_slide -->

## A first docker-compose.yml

```yaml
services:
  api:
    build: .
    ports:
      - "8080:8080"
  db:
    image: postgres:16
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
      POSTGRES_DB: appdb
```

Two services: `api` (built from the `Dockerfile` in this directory) and
`db` (pulled ready-made from Docker Hub — no Dockerfile needed for
something that already ships an official image).

<!-- pause -->

**Demo:**

```
docker-compose up --build
```

```
curl localhost:8080/healthz
```

```
docker-compose down
```

<!--
speaker_note: |
  Keep this first pass deliberately simple - no healthcheck, no named
  volume yet. The point is just "one file, one command, two containers
  come up talking to each other." The production-shaped version with
  healthchecks and a persistent volume is the very next section, framed
  explicitly as an upgrade from this one - not as new, unrelated
  content.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

docker-compose for local Go development
===

<!-- end_slide -->

## The shape of it

Same two services as the simple version you just saw — now with the
pieces a real dev setup actually needs: an explicit `DATABASE_URL`, a
named volume so data survives a restart, and a healthcheck gating when
`api` is allowed to start.

```yaml
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://app:app@db:5432/appdb?sslmode=disable
    depends_on:
      db:
        condition: service_healthy
  db:
    image: postgres:16
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
      POSTGRES_DB: appdb
    volumes:
      - dbdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app"]
      interval: 5s
      timeout: 5s
      retries: 5
volumes:
  dbdata:
```

<!-- end_slide -->

## Service names are hostnames, not `localhost`

Compose puts every service in the file on one shared network
automatically. From inside the `api` container, Postgres isn't at
`localhost` — it's at `db`, the service name.

<!-- pause -->

**This is a genuinely common first mistake.** `localhost` inside a
container means *that container*, not the host machine and not a
sibling container. Point your connection string at the service name.

**Demo:** point `DATABASE_URL` at `localhost` instead of `db`, bring the stack up, and read the connection failure together.

<!--
speaker_note: |
  Everyone who's containerized a multi-service app for the first time
  has hit this exactly once. Let them watch the failure rather than
  just describing it - it's a five-second demo with a memorable payoff.
-->

<!-- end_slide -->

## `depends_on` doesn't mean "ready"

A plain `depends_on: - db` only waits for the `db` **container process
to start**. Postgres inside it can still be mid-initialization —
accepting no connections yet — when your `api` container starts and
tries to connect.

<!-- pause -->

**The fix is the healthcheck-gated form:**

```yaml
depends_on:
  db:
    condition: service_healthy
```

Combined with a `healthcheck` block on `db` (`pg_isready`), `api` now
waits for Postgres to be genuinely ready to accept connections — a real
race condition, closed properly.

<!-- end_slide -->

## Volumes: surviving `down` and `up`

```yaml
volumes:
  - dbdata:/var/lib/postgresql/data
```

A **named volume** persists data outside the container's own
filesystem. `docker-compose down` removes containers; the data in
`dbdata` survives. `docker-compose up` again, and your data is still
there.

<!-- pause -->

`docker-compose down -v` removes the volume too — genuinely wipes the
database. Know which one you're running.

<!-- end_slide -->

## Two honest dev-loop options

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Full compose**

```
docker-compose up --build
```

Rebuilds the Go image on every code
change. Closer to what actually ships
to production. Slower edit/rebuild
loop.

<!-- column: 1 -->

**Host Go, containerized dependency only**

```
docker-compose up db
go run ./cmd/server
```

Fast iteration — no image rebuild at
all. Drifts further from the
production environment.

<!-- reset_layout -->

**Type in chat: which would you reach for day-to-day, and which for the last check before opening a PR?**

<!--
speaker_note: |
  There's no declared right answer here and the slide shouldn't
  pretend there is - let the room actually disagree a little. Most
  working teams land on "host Go day-to-day, full compose before
  pushing," but that's a discussion beat, not a rule to hand down.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

WSL 
===

<!--
speaker_note: |
  Signal explicitly that this section exists because the client asked
  for it by name on the scoping call. This isn't generic Docker
  Desktop boilerplate - every point here is chosen because it's a real
  thing this specific client's team will hit.
-->

<!-- end_slide -->

## Run Docker commands from inside WSL, not from Windows

Docker Desktop's WSL2 backend runs the actual Docker engine inside a
lightweight Linux VM that WSL2 manages.

<!-- pause -->

**Run `docker` and `docker-compose` from inside a WSL distro's shell**
(Ubuntu, etc.) — not PowerShell or CMD. Running from the Windows side
works in a limited sense, but you lose performance and correct
file-permission behavior that running natively inside the Linux VM
gives you for free.

<!-- end_slide -->

## Filesystem location matters — a lot

WSL2 keeps two separate filesystems in play:

<!-- incremental_lists: true -->

- **Linux-native**, on a real ext4-backed virtual disk — e.g. `~/projects/myapp`
- **Windows-mounted**, e.g. `/mnt/c/Users/you/projects/myapp` — a bind mount across the Windows/Linux boundary

<!-- incremental_lists: false -->

<!-- pause -->

**Keep your project on the Linux side.** Bind-mounting a `/mnt/c/...`
path into a container crosses that boundary on *every single file
operation*.

<!-- end_slide -->

## What crossing the boundary actually costs you

<!-- incremental_lists: true -->

- `go build` measurably slower, when source or the module cache lives cross-filesystem
- File-watching tools (hot-reload, `air`, etc.) can miss changes entirely — Windows filesystem-change notifications don't always propagate cleanly across the boundary
- Every `COPY . .` in a Docker build context pays the same cost, on every build

<!-- incremental_lists: false -->

**This isn't a minor tax. For a team building and rebuilding containers all day, it's the difference between a fast loop and a frustrating one.**

<!--
speaker_note: |
  This is worth landing concretely - "slow" is abstract, "your hot
  reload tool silently stops noticing your changes" is the kind of
  thing that actually gets reported as a bug against the wrong tool.
-->

<!-- end_slide -->

## Line endings: a real, cryptic failure mode

A repo checked out or edited on the Windows side can pick up **CRLF**
line endings. Copy a shell script with CRLF endings — `entrypoint.sh`,
say — into a Linux container, and it fails.

<!-- pause -->

```
exec /entrypoint.sh: no such file or directory
```

**The script is right there. The error lies to you.** CRLF appends an
invisible `\r` to the shebang line, so the kernel looks for an
interpreter literally named `/bin/sh\r` — which doesn't exist.

<!-- pause -->

**Fix:** a `.gitattributes` file with `* text=auto eol=lf`, plus an
editor configured to save this project's files with LF.

<!--
speaker_note: |
  This is the lab's Exercise 5 - if time allows, cause this failure
  live rather than only describing it. The gap between "script clearly
  exists" and "no such file or directory" is what makes this stick.
-->

<!-- end_slide -->

## Two smaller WSL notes

<!-- column_layout: [1, 1] -->

<!-- column: 0 -->

**Resource limits**

WSL2's memory/CPU allocation to its Linux VM is configurable via a
`.wslconfig` file in the Windows user profile. If containers start
getting OOM-killed unexpectedly during a `docker-compose up` with
several services, this is the first place to look.

<!-- column: 1 -->

**Port forwarding**

Modern WSL2 + Docker Desktop generally just works —
`localhost:8080` on the Windows host reaches a container's exposed
port automatically. If you're on an older or manual
Docker-Engine-in-WSL setup and it *doesn't* just work, that's the
likely reason: it may need manual port-proxy configuration.

<!-- reset_layout -->

<!-- end_slide -->

<!-- jump_to_middle -->

A brief, honest word on production
===

<!-- end_slide -->

## This isn't a Kubernetes course

The same small, static, multi-stage image that runs locally via
`docker-compose up` is exactly what you'd push to a registry and deploy.

<!-- pause -->

**Where it actually runs varies widely:** a single VM running `docker
run`, a managed platform like ECS or Cloud Run, or a Kubernetes cluster.
That decision — and the CI/CD pipeline that gets an image there — is
genuinely out of scope for this course.

**What isn't out of scope: the image-building discipline.** Small,
static, dependency-free. That discipline is identical regardless of
where the image lands.

<!--
speaker_note: |
  Be honest here rather than hand-waving - name Kubernetes and CI/CD
  explicitly as things this course didn't cover, rather than letting
  delegates assume they got a partial version of that content. Point
  at it as a real next step, don't pretend to have taught it.
-->

<!-- end_slide -->

## Summary

<!-- incremental_lists: true -->

1. **Go's single-binary story makes small images the natural outcome**: no runtime, no package pile, just the binary
2. **Multi-stage builds are how you actually get there**: a large builder stage that never ships, a minimal final stage that does
3. **`CGO_ENABLED=0` and `scratch` vs `distroless` are real, specific gotchas**: silent libc dependencies and missing CA certs, not hypotheticals
4. **docker-compose networking uses service names, not `localhost`**, and `depends_on` needs a healthcheck to mean "actually ready"
5. **WSL rewards Linux-side project location and LF line endings**: the filesystem boundary and CRLF failures are real, specific, and avoidable
6. **The same image-building discipline carries to production**, wherever that ends up being

<!-- end_slide -->

## Back to the opening scenario

900MB versus 8MB. Three minutes versus twenty seconds. Now you know
exactly why: a Python or Node image ships a full runtime plus every
installed package, in every layer. A Go image, built with a proper
multi-stage Dockerfile, ships a statically linked binary and little
else.

<!-- pause -->

**It's not that Go is "faster" in some vague sense. It's that Go had
almost nothing to ship in the first place — and Docker just reflects
that back at you, honestly, in megabytes.**

<!--
speaker_note: |
  This resolves the opening chat poll - go back to a couple of the
  original guesses if you noted them, and show which ones were closest
  to the real mechanism (static linking + multi-stage builds) versus
  the vaguer "Go is just faster" instinct.
-->

<!-- end_slide -->

## That's the course

Day one, Topic 1, the pitch was: **`go build`, hand over the file,
done.** No interpreter, no runtime, no "works on my machine."

<!-- pause -->

Four days later, that exact same property is why a Go container image
is a handful of megabytes instead of hundreds — the single-binary story
from Topic 1 and the tiny Docker image from this topic are **the same
fact**, looked at from two different angles.

<!-- pause -->

**Everything in between — types, interfaces, concurrency, patterns, REST
and gRPC services, tests — is what you now know how to build. This topic
is how you ship it.**

<!--
speaker_note: |
  This is the close of the whole course, not just this topic - there
  is no Topic 14, say that out loud. Give the throughline a beat before
  moving to Questions - the teammate from Topic 1 who couldn't run a
  Python script without the right interpreter and packages installed
  is the same problem this topic solved for container images.
-->

<!-- end_slide -->

<!-- jump_to_middle -->

Questions?
===

<!-- end_slide -->
