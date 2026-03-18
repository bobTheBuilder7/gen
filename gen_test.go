package gen

import (
	"bytes"
	"testing"
)

func TestImport(t *testing.T) {
	got := Import("", "fmt").String()
	want := `import "fmt"`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestImportWithAlias(t *testing.T) {
	got := Import("f ", "fmt").String()
	want := `import f "fmt"`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestVar(t *testing.T) {
	got := Var("x", Int(42)).String()
	want := "var x = 42"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestConst(t *testing.T) {
	got := Const("name", String("hello")).String()
	want := `const name = "hello"`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestValueWrappers(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"String", String("foo").String(), `"foo"`},
		{"Int", Int(7).String(), "7"},
		{"Bool", Bool(true).String(), "true"},
		{"Float", Float(3.14).String(), "3.140000"},
		{"Rune", Rune('A').String(), "'A'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestStruct(t *testing.T) {
	got := Struct("User",
		Field{Name: "Name", Type: "string"},
		Field{Name: "Age", Type: "int", Tag: `json:"age"`},
	).String()
	want := "type User struct {\nName string\nAge int `json:\"age\"`\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestInterface(t *testing.T) {
	got := Interface("Reader",
		Method{Name: "Read", Params: "p []byte", Returns: "(int, error)"},
	).String()
	want := "type Reader interface {\nRead(p []byte) (int, error)\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFunc(t *testing.T) {
	got := Func("Add", "a, b int", "int",
		Line("return a + b"),
	).String()
	want := "func Add(a, b int) int {\nreturn a + b\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestMethodFunc(t *testing.T) {
	got := MethodFunc("s *Server", "Start", "", "error",
		Line("return nil"),
	).String()
	want := "func (s *Server) Start() error {\nreturn nil\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestCallAndErrCheck(t *testing.T) {
	got := Call("result, err", "doSomething", Arg("ctx"), String("test")).String()
	want := `result, err := doSomething(ctx, "test")`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	got = ErrCheck(Arg("nil")).String()
	want = "if err != nil {\nreturn nil, err\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	got = ErrCheck(Arg("")).String()
	want = "if err != nil {\nreturn err\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestTypeAlias(t *testing.T) {
	got := TypeAlias("UserID", "int64").String()
	want := "type UserID int64"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFieldWithoutTag(t *testing.T) {
	got := (Field{Name: "ID", Type: "int"}).String()
	want := "ID int"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFieldWithTag(t *testing.T) {
	got := (Field{Name: "Name", Type: "string", Tag: `json:"name,omitempty"`}).String()
	want := "Name string `json:\"name,omitempty\"`"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStructEmpty(t *testing.T) {
	got := Struct("Empty").String()
	want := "type Empty struct {\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestInterfaceMultipleMethods(t *testing.T) {
	got := Interface("ReadWriter",
		Method{Name: "Read", Params: "p []byte", Returns: "(int, error)"},
		Method{Name: "Write", Params: "p []byte", Returns: "(int, error)"},
	).String()
	want := "type ReadWriter interface {\nRead(p []byte) (int, error)\nWrite(p []byte) (int, error)\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFuncNoReturn(t *testing.T) {
	got := Func("doWork", "", "",
		Line("fmt.Println()"),
	).String()
	want := "func doWork() {\nfmt.Println()\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFuncMultipleBodyLines(t *testing.T) {
	got := Func("greet", "name string", "",
		Var("msg", String("hello")),
		Call("", "fmt.Println", Arg("msg")),
	).String()
	want := "func greet(name string) {\nvar msg = \"hello\"\nfmt.Println(msg)\n}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestCallNoAssign(t *testing.T) {
	got := Call("", "fmt.Println", String("hi")).String()
	want := `fmt.Println("hi")`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestCallNoArgs(t *testing.T) {
	got := Call("", "doWork").String()
	want := "doWork()"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLine(t *testing.T) {
	got := Line("x := 5").String()
	want := "x := 5"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFileWriteToNoComment(t *testing.T) {
	f := NewFile("util")
	f.AddBlock(Var("x", Int(1)))

	var buf bytes.Buffer
	if err := f.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}

	want := "package util\n\nvar x = 1\n"
	if buf.String() != want {
		t.Fatalf("got %q, want %q", buf.String(), want)
	}
}

func TestFileWriteToMultiLineComment(t *testing.T) {
	f := NewFile("main", "line one", "line two")

	var buf bytes.Buffer
	if err := f.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}

	want := "// line one\n// line two\npackage main\n\n"
	if buf.String() != want {
		t.Fatalf("got %q, want %q", buf.String(), want)
	}
}

func TestFileWriteTo(t *testing.T) {
	f := NewFile("main", "generated code")
	f.AddBlock(Import("", "fmt"))
	f.AddBlock(Func("main", "", "",
		Call("", "fmt.Println", String("hello")),
	))

	var buf bytes.Buffer
	if err := f.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}

	want := "// generated code\npackage main\n\nimport \"fmt\"\nfunc main() {\nfmt.Println(\"hello\")\n}\n"
	if buf.String() != want {
		t.Fatalf("got %q, want %q", buf.String(), want)
	}
}
