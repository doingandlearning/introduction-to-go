# Lab 13: Docker & Deployment — containerize a Go service, break it on purpose

Starter code is in `starter/` (an incomplete `Dockerfile`, a buggy
`docker-compose.yml`, TODOs implied by the comments in each). A complete
reference is in `solution/` — don't look until you've had a go.

Both directories share the same Go program (`cmd/server/main.go`): a
small HTTP API with `/healthz`, `/items`, and `/outbound` endpoints,
standard library only. By the end of this lab you'll have written a real
multi-stage Dockerfile, watched a container fail on a missing
dependency it never had, wired up docker-compose networking correctly
after getting it wrong on purpose, and fixed a line-ending bug that
looks nothing like a line-ending bug when it happens.

---

## Exercise 1: Write a real multi-stage Dockerfile

**Objective:** Turn the naive, single-stage starter Dockerfile into the
multi-stage pattern from lecture.

**Context:** `starter/Dockerfile` builds and runs, but ships the entire
`golang:1.22` toolchain in the final image just to run one binary.

**Tasks:**

1. Rewrite `starter/Dockerfile` as two stages: a `builder` stage (based
   on `golang:1.22`) that compiles the binary with
   `CGO_ENABLED=0 GOOS=linux go build`, and a final stage that starts
   `FROM gcr.io/distroless/static-debian12` and copies in only the
   compiled binary with `COPY --from=builder`.
2. Build it: `docker build -t docker-lab .`
3. Run it: `docker run -p 8080:8080 docker-lab`
4. Confirm `curl localhost:8080/healthz` and `curl localhost:8080/items`
   both respond correctly.
5. Compare `docker images` output for your multi-stage build against
   the original single-stage version. Note the size difference.

**Key Learning:** The builder stage's size doesn't matter — it never
ships. Only what the final `FROM` starts with, plus whatever you `COPY`
into it, ends up in the image you actually deploy.

---

## Exercise 2: `scratch` vs `distroless` — the missing CA certs

**Objective:** Hit a real, commonly-encountered Docker gotcha instead of
just being told about it.

**Tasks:**

1. Starting from your working Exercise 1 Dockerfile, change the final
   stage's base image from `gcr.io/distroless/static-debian12` to
   `scratch`. Rebuild and rerun.
2. Hit `curl localhost:8080/outbound`. Read the error carefully.
3. Reason about why it fails on `scratch` but didn't on `distroless` —
   what does `distroless/static-debian12` include that plain `scratch`
   doesn't?
4. Fix it: add `COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt`
   above your existing `COPY --from=builder /server /server` line.
   Rebuild, rerun, confirm `/outbound` now succeeds.

**Key Learning:** `scratch` is genuinely empty — no CA trust store, no
shell, nothing. Any outbound HTTPS call from a `scratch`-based container
fails until you explicitly copy in a CA bundle. `distroless` images
include one by default, which is most of why they're a gentler starting
point than `scratch` day one.

---

## Exercise 3: Wire up docker-compose correctly

**Objective:** Fix two real bugs in the starter `docker-compose.yml`,
one at a time, and understand why each one is wrong before fixing it.

**Context:** `starter/docker-compose.yml` brings up `api` and a
`db` (Postgres) service, but `api` is configured to reach Postgres at
`localhost`, and there's no healthcheck gating the startup order.

**Tasks:**

1. Run `docker-compose up --build` as-is. Watch the `api` service's logs
   for a connection failure (or reason through why `DATABASE_URL`
   pointing at `localhost` can never reach the `db` container from
   inside the `api` container).
2. Fix the hostname: change `localhost` to `db` — the service name
   Compose gives the `db` container on the shared network it creates
   automatically.
3. Add a `healthcheck` block to the `db` service
   (`test: ["CMD-SHELL", "pg_isready -U app"]`, with `interval`,
   `timeout`, and `retries`), then change `api`'s `depends_on` from a
   plain list (`- db`) to the long form gated on
   `condition: service_healthy`.
4. Tear everything down (`docker-compose down -v`, which also drops the
   named volume) and bring it up fresh. Confirm `api` now waits for
   Postgres to be genuinely ready, not just "container started."

**Key Learning:** Container-to-container networking inside a Compose
project uses service names as hostnames, never `localhost` — and a plain
`depends_on` only waits for a container process to start, not for
whatever's running inside it to be ready to do useful work. Those are
two different, both-common mistakes.

---

## Exercise 4: The WSL filesystem boundary

**Objective:** Confirm (or reason through) why project location on disk
matters under WSL2.

**Tasks — if you're working from a real WSL environment:**

1. From inside your WSL distro's shell, run `pwd` in this project's
   directory. Confirm it does **not** start with `/mnt/c/` (or any other
   `/mnt/<drive>/`) — if it does, you're on the Windows-filesystem side
   of the boundary.
2. If it does start with `/mnt/c/...`, copy the whole project to
   somewhere under your Linux home directory instead (e.g.
   `~/projects/docker-lab`), and re-run from there.
3. Time `docker-compose build` (or just `docker build .`) from both
   locations if you can. Note the difference.

**Tasks — if you're not on WSL, reason through it in writing instead:**

1. `COPY . .` in the Dockerfile, and any file-watching tool during local
   dev, both touch every file in the project repeatedly. Explain, in
   your own words, why crossing the Windows-filesystem/Linux-VM boundary
   on every one of those operations would be slower than staying
   entirely on one side of it.
2. Name one concrete symptom (beyond raw speed) that this boundary can
   cause for a developer relying on hot-reload tooling.

**Key Learning:** WSL2 keeps its own Linux-native virtual disk separate
from the Windows filesystem it mounts at `/mnt/c/...`. Every file
operation that crosses that boundary — which a `COPY . .` build context
or a file watcher does constantly — pays a real, measurable cost that a
project entirely on one side of the boundary doesn't.

---

## Exercise 5: Break it with a line ending, on purpose

**Objective:** Cause and then diagnose one of the most confusing
Docker/WSL failures there is: a script with the wrong line endings.

**Context:** `starter/Dockerfile.entrypoint-demo` runs the server via
`starter/entrypoint.sh` instead of invoking the binary directly. It
works as shipped.

**Tasks:**

1. Build and run it as-is first, to confirm it works:
   `docker build -f Dockerfile.entrypoint-demo -t docker-lab-entrypoint .`
   then `docker run -p 8080:8080 docker-lab-entrypoint`. Confirm
   `/healthz` responds.
2. Now corrupt it: convert `entrypoint.sh`'s line endings from LF to
   CRLF. (On a real Windows/WSL setup this often happens by accident —
   an editor set to Windows line endings, or a `git` checkout without
   `core.autocrlf`/`.gitattributes` configured. To do it deliberately:
   `unix2dos entrypoint.sh`, or open it in an editor and re-save with
   CRLF line endings.)
3. Rebuild and rerun. Read the error closely — it will **not** say
   anything about line endings.
4. Explain, in your own words, why a CRLF line ending on a shebang
   (`#!/bin/sh`) line produces the specific error you saw.
5. Fix it: convert back to LF (`dos2unix entrypoint.sh`, or re-save with
   LF), rebuild, confirm it works again. Then add a `.gitattributes`
   file to this project containing `* text=auto eol=lf`, so this can't
   silently happen again on a Windows checkout.

**Key Learning:** A CRLF line ending appends an invisible `\r` to the
end of the shebang line. The kernel then tries to find an interpreter
literally named `/bin/sh\r`, which doesn't exist — producing a cryptic
"no such file or directory" error that has nothing obviously to do with
line endings. This is a genuinely common failure mode for any repo that
crosses between Windows and Linux editing environments.

---

## Summary

By the end of this lab you should be able to:

- Write a multi-stage Dockerfile that compiles in one stage and ships
  only the binary in a small final stage
- Explain concretely why `scratch` needs its CA bundle copied in by hand
  and `distroless` doesn't
- Diagnose and fix both the `localhost`-vs-service-name mistake and the
  missing-healthcheck race condition in docker-compose networking
- Explain why project location matters under WSL2, in terms of what
  actually crosses the filesystem boundary
- Recognize a CRLF-corrupted shebang line from its symptom, not just its
  cause
