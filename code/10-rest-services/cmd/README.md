# Topic 10 — standalone live-coding demos

Small, self-contained `main.go`s for the moments in the slides marked
`Demo:` or `speaker_note`, separate from the full CRUD app in
`cmd/server`. Each one is a single file, no `internal/` imports, safe to
run, break, and restart without touching the layered app.

- **minimal-ping** — the opening "what a REST service is" slide. Type it
  live from an empty file if you want the full effect.
- **routing** — Go 1.22 method + path-param routing, three routes, one
  path shape.
- **struct-tag-footgun** — flip one `json` tag to demo the silent-failure
  footgun. Instructions are in the file's header comment.
- **recover-middleware** — panic in one handler, `/ping` still answers
  right after. Header comment explains the optional "without recover"
  contrast run.
- **nested-json** — one struct showing three JSON-nesting shapes
  together: a named nested field (real nesting), an embedded/anonymous
  struct (fields flatten/promote instead of nesting — the gotcha), and a
  slice of structs (JSON array of objects).

Each is run individually from the topic's module root:

```
go run ./cmd/minimal-ping
go run ./cmd/routing
go run ./cmd/struct-tag-footgun
go run ./cmd/recover-middleware
go run ./cmd/nested-json
```

Only one at a time — they all bind `:8080`.
