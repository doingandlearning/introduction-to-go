# Topic 13 sample code — Docker & Deployment

One small HTTP API (`cmd/server`), one `Dockerfile`, one `docker-compose.yml`.
Run everything from this directory (`code/13-docker-and-deployment/`, or
wherever this folder lands in the course repo).

## `cmd/server`

Three endpoints, standard library only:

- `GET /healthz` — plain `200 ok`, for compose healthchecks and manual pokes
- `GET /items` — a small hardcoded JSON list, standing in for Topic 10's
  REST service
- `GET /outbound` — makes a real outbound HTTPS call, used later to make
  the "`scratch` has no CA certificates" gotcha observable rather than
  theoretical

Reads `PORT` (defaults to `8080`) and `DATABASE_URL` (logged, not
connected to — this stand-in doesn't need a real database to be worth
containerizing).

## Build and run the image directly

```
docker build -t docker-demo .
docker run -p 8080:8080 docker-demo

curl localhost:8080/healthz
curl localhost:8080/items
```

Check the image size once it's built — this is the whole pitch from the
lecture:

```
docker images docker-demo
# a handful of MB: the binary, plus whatever distroless/static-debian12
# needs for TLS and timezones. No Go toolchain, no shell, no package
# manager shipped in this image at all.
```

## Bring up the service plus a real Postgres

```
docker-compose up --build
```

This starts `api` and `db`. Compose won't start `api` until `db`'s
healthcheck (`pg_isready`) reports healthy — not just "container
started" — so there's no race between the API booting and Postgres
actually being ready to accept connections.

```
curl localhost:8080/healthz
docker-compose down        # stops both containers
docker-compose up          # data in the dbdata volume survives — no --build needed
                            # unless cmd/server changed
```

Tear down and wipe the database entirely with `docker-compose down -v`.

## Faster local dev loop

Rebuilding the Go image on every code change is a slower loop than
editing Go normally. An alternative for local iteration: only start the
dependency in a container, and run the Go service on the host:

```
docker-compose up db      # just Postgres, healthcheck and all
DATABASE_URL="postgres://app:app@localhost:5432/appdb?sslmode=disable" go run ./cmd/server
```

Note the hostname flips to `localhost` here — from the host machine,
Postgres's exposed port is reachable at `localhost`; it's only
*other containers* that need to say `db`. This is faster to iterate on,
at the cost of drifting further from what actually ships to production.
Neither approach is "correct" — it's a genuine tradeoff, covered in the
lecture.

---

### WSL callout

If you're on Windows, run every command above from **inside a WSL
distro's shell** (Ubuntu, etc.) — not PowerShell or CMD — and keep this
project on the Linux filesystem (`~/projects/...`), not a Windows-drive
mount (`/mnt/c/...`). Both matter for build speed and file-watching
correctness. See the lecture slides for the full reasoning.
