# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gen` is a Go library (`github.com/bobTheBuilder7/gen`) for programmatic Go source code generation. It provides a fluent builder API where code elements (imports, structs, functions, etc.) are composed as `fmt.Stringer` blocks and written to any `io.Writer`. Zero external dependencies.

## Commands

```bash
go build ./...           # Build library and examples
go vet ./...             # Static analysis
go run ./examples/       # Run the example generator (outputs generated.go in examples/)
```

No tests exist yet. When adding them, use standard `go test ./...`.

## Architecture

All code lives in two files in the root package:

- **gen.go** — Core API. `File` struct holds a package name and ordered list of `fmt.Stringer` blocks. Blocks are added via `AddBlock()` and flushed via `WriteTo(io.Writer)`. Each Go construct (import, var, const, struct, interface, function) is a private type implementing `fmt.Stringer`. Public constructors (`Import()`, `Struct()`, `Func()`, `MethodFunc()`, etc.) create these types.
- **util.go** — `join()` helper for joining `fmt.Stringer` slices with a separator (used by `Call()`).

Key design decisions:
- `File` is thread-safe via `sync.RWMutex` — blocks can be added concurrently.
- `WriteTo` does **not** auto-format output (removed in cfdf011). Callers handle formatting (e.g., pipe through `go/format`).
- `Func()` creates standalone functions; `MethodFunc()` creates methods with a receiver.
- `Call()` generates function call expressions with optional assignment (`assigns := name(args...)`).
- `ErrCheck()` generates `if err != nil { return ... }` boilerplate; empty `Arg("")` produces single-return error check.
- Value wrappers (`String()`, `Int()`, `Bool()`, `Float()`, `Rune()`) produce properly formatted Go literals as `fmt.Stringer`.
