# gen

A lightweight Go library for programmatic Go source code generation.

Build Go files with a clean, composable API — no templates, no string concatenation, no external dependencies.

## Install

```bash
go get github.com/bobTheBuilder7/gen
```

## Quick Start

```go
package main

import (
    "os"
    "github.com/bobTheBuilder7/gen"
)

func main() {
    f := gen.NewFile("models")

    f.AddBlock(gen.Import("", "time"))

    f.AddBlock(gen.Struct("User",
        gen.Field{Name: "ID", Type: "int", Tag: `json:"id" db:"id"`},
        gen.Field{Name: "Name", Type: "string", Tag: `json:"name"`},
        gen.Field{Name: "CreatedAt", Type: "time.Time", Tag: `json:"created_at"`},
    ))

    f.AddBlock(gen.Interface("UserStore",
        gen.Method{Name: "FindByID", Params: "id int", Returns: "(*User, error)"},
        gen.Method{Name: "Save", Params: "u *User", Returns: "error"},
    ))

    out, _ := os.Create("models.go")
    defer out.Close()
    f.WriteTo(out)
}
```

This generates:

```go
package models

import "time"

type User struct {
    ID        int       `json:"id" db:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

type UserStore interface {
    FindByID(id int) (*User, error)
    Save(u *User) error
}
```

## API

### File

```go
f := gen.NewFile("pkg")              // create a new file with package name
f := gen.NewFile("pkg", "comment")   // with a file-level comment
f.AddBlock(block)                    // add any code block (thread-safe)
f.WriteTo(w)                         // write generated code to an io.Writer
```

### Imports

```go
gen.Import("", "fmt")               // import "fmt"
gen.Import("pb ", "google/pb")      // import pb "google/pb"
```

### Constants & Variables

```go
gen.Const("Version", gen.String("1.0.0"))   // const Version = "1.0.0"
gen.Const("MaxRetries", gen.Int(3))         // const MaxRetries = 3
gen.Var("Debug", gen.Bool(false))           // var Debug = false
```

### Type Aliases

```go
gen.TypeAlias("ID", "int64")                // type ID int64
gen.TypeAlias("Handler", "func()")          // type Handler func()
```

### Structs

```go
gen.Struct("Config",
    gen.Field{Name: "Host", Type: "string", Tag: `json:"host"`},
    gen.Field{Name: "Port", Type: "int"},
)
```

### Interfaces

```go
gen.Interface("Reader",
    gen.Method{Name: "Read", Params: "p []byte", Returns: "(int, error)"},
)
```

### Functions & Methods

```go
// Standalone function
gen.Func("NewUser", "name string", "*User",
    gen.Line("return &User{Name: name}"),
)

// Method with receiver
gen.MethodFunc("u *User", "Validate", "", "error",
    gen.Call("err", "validate", gen.Arg("u.Name")),
    gen.ErrCheck(gen.Arg("")),
    gen.Line("return nil"),
)
```

### Value Helpers

| Helper | Example | Output |
|--------|---------|--------|
| `gen.String("hi")` | `gen.Const("X", gen.String("hi"))` | `const X = "hi"` |
| `gen.Int(42)` | `gen.Const("X", gen.Int(42))` | `const X = 42` |
| `gen.Bool(true)` | `gen.Var("X", gen.Bool(true))` | `var X = true` |
| `gen.Float(3.14)` | `gen.Const("X", gen.Float(3.14))` | `const X = 3.140000` |
| `gen.Rune(':')` | `gen.Const("X", gen.Rune(':'))` | `const X = ':'` |

### Code Helpers

```go
gen.Call("result", "doWork", gen.Arg("ctx"))   // result := doWork(ctx)
gen.Call("", "fmt.Println", gen.Arg("msg"))    // fmt.Println(msg)
gen.ErrCheck(gen.Arg(""))                      // if err != nil { return err }
gen.ErrCheck(gen.Arg("nil"))                   // if err != nil { return nil, err }
gen.Line("return nil")                         // raw code line
```

## Features

- **Zero dependencies** — only Go standard library
- **Thread-safe** — add blocks concurrently with `AddBlock()`
- **Composable** — all blocks implement `fmt.Stringer`, mix and match freely
- **Flexible output** — write to files, buffers, stdout, or pipe through `go/format`

## License

MIT
