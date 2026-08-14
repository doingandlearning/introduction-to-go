# --- build stage ---
FROM golang:1.22 AS builder
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 keeps this a fully static binary — the same "one
# dependency-free binary" property from Topic 1. That property is
# exactly what lets the final image below be tiny: there's no runtime,
# no interpreter, and no shared library dependencies to drag along.
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# --- final stage ---
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/server /usr/local/bin/server

EXPOSE 8080
ENTRYPOINT ["server"]

# Note for teaching: this final image is small because it's copying in
# one static binary and nothing else — no Go toolchain, no source, no
# package manager cache from the build stage. Swapping alpine for
# gcr.io/distroless/static or FROM scratch shrinks it further; Topic 13
# goes deeper on that tradeoff. This topic only needs "small enough to
# run locally without thinking about it."
