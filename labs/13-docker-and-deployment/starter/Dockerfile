# Starter Dockerfile — Exercise 1 & 2.
#
# This builds and runs, but it's a single stage: the final image ships
# the entire golang:1.22 toolchain (compiler, module cache, everything)
# just to run one binary. That's the opposite of the multi-stage pattern
# from lecture.
#
# Exercise 1: rewrite this as a proper multi-stage build — a "builder"
# stage that compiles the binary, and a small final stage that only
# contains the binary. Use gcr.io/distroless/static-debian12 as the
# final base to start with.
#
# Exercise 2: once Exercise 1 works, change the final stage's base image
# to `scratch` instead, rebuild, and hit /outbound. Reason about what
# breaks and why, then fix it.

FROM golang:1.22
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server
EXPOSE 8080
ENTRYPOINT ["/server"]
