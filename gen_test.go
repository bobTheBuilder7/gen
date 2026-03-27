package gen

import (
	"bytes"
	"testing"

	"github.com/bobTheBuilder7/gen/assert"
)

func TestImport(t *testing.T) {
	assert.Equal(t, Import("", "fmt").String(), `import "fmt"`)
}

func TestImportWithAlias(t *testing.T) {
	assert.Equal(t, Import("f ", "fmt").String(), `import f "fmt"`)
}

func TestVar(t *testing.T) {
	assert.Equal(t, Var("x", Int(42)).String(), "var x = 42")
}

func TestConst(t *testing.T) {
	assert.Equal(t, Const("name", String("hello")).String(), `const name = "hello"`)
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
			assert.Equal(t, tt.got, tt.want)
		})
	}
}

func TestStruct(t *testing.T) {
	got := Struct("User",
		Field{Name: "Name", Type: "string"},
		Field{Name: "Age", Type: "int", Tag: `json:"age"`},
	).String()
	assert.Equal(t, got, "type User struct {\nName string\nAge int `json:\"age\"`\n}")
}

func TestInterface(t *testing.T) {
	got := Interface("Reader",
		Method{Name: "Read", Params: "p []byte", Returns: "(int, error)"},
	).String()
	assert.Equal(t, got, "type Reader interface {\nRead(p []byte) (int, error)\n}")
}

func TestFunc(t *testing.T) {
	got := Func("Add", "a, b int", "int",
		Line("return a + b"),
	).String()
	assert.Equal(t, got, "func Add(a, b int) int {\nreturn a + b\n}")
}

func TestMethodFunc(t *testing.T) {
	got := MethodFunc("s *Server", "Start", "", "error",
		Line("return nil"),
	).String()
	assert.Equal(t, got, "func (s *Server) Start() error {\nreturn nil\n}")
}

func TestCallAndErrCheck(t *testing.T) {
	got := Call("result, err", "doSomething", Arg("ctx"), String("test")).String()
	assert.Equal(t, got, `result, err := doSomething(ctx, "test")`)

	got = ErrCheck(Arg("nil")).String()
	assert.Equal(t, got, "if err != nil {\nreturn nil, err\n}")

	got = ErrCheck(Arg("")).String()
	assert.Equal(t, got, "if err != nil {\nreturn err\n}")
}

func TestTypeAlias(t *testing.T) {
	assert.Equal(t, TypeAlias("UserID", "int64").String(), "type UserID int64")
}

func TestFieldWithoutTag(t *testing.T) {
	assert.Equal(t, (Field{Name: "ID", Type: "int"}).String(), "ID int")
}

func TestFieldWithTag(t *testing.T) {
	assert.Equal(t, (Field{Name: "Name", Type: "string", Tag: `json:"name,omitempty"`}).String(), "Name string `json:\"name,omitempty\"`")
}

func TestStructEmpty(t *testing.T) {
	assert.Equal(t, Struct("Empty").String(), "type Empty struct {\n}")
}

func TestInterfaceMultipleMethods(t *testing.T) {
	got := Interface("ReadWriter",
		Method{Name: "Read", Params: "p []byte", Returns: "(int, error)"},
		Method{Name: "Write", Params: "p []byte", Returns: "(int, error)"},
	).String()
	assert.Equal(t, got, "type ReadWriter interface {\nRead(p []byte) (int, error)\nWrite(p []byte) (int, error)\n}")
}

func TestFuncNoReturn(t *testing.T) {
	got := Func("doWork", "", "",
		Line("fmt.Println()"),
	).String()
	assert.Equal(t, got, "func doWork() {\nfmt.Println()\n}")
}

func TestFuncMultipleBodyLines(t *testing.T) {
	got := Func("greet", "name string", "",
		Var("msg", String("hello")),
		Call("", "fmt.Println", Arg("msg")),
	).String()
	assert.Equal(t, got, "func greet(name string) {\nvar msg = \"hello\"\nfmt.Println(msg)\n}")
}

func TestCallNoAssign(t *testing.T) {
	assert.Equal(t, Call("", "fmt.Println", String("hi")).String(), `fmt.Println("hi")`)
}

func TestCallNoArgs(t *testing.T) {
	assert.Equal(t, Call("", "doWork").String(), "doWork()")
}

func TestLine(t *testing.T) {
	assert.Equal(t, Line("x := 5").String(), "x := 5")
}

func TestFileWriteToNoComment(t *testing.T) {
	f := NewFile("util")
	f.AddBlock(Var("x", Int(1)))

	var buf bytes.Buffer
	assert.Nil(t, f.WriteTo(&buf))
	assert.Equal(t, buf.String(), "package util\n\nvar x = 1\n")
}

func TestFileWriteToMultiLineComment(t *testing.T) {
	f := NewFile("main", "line one", "line two")

	var buf bytes.Buffer
	assert.Nil(t, f.WriteTo(&buf))
	assert.Equal(t, buf.String(), "// line one\n// line two\npackage main\n\n")
}

func TestFileWriteTo(t *testing.T) {
	f := NewFile("main", "generated code")
	f.AddBlock(Import("", "fmt"))
	f.AddBlock(Func("main", "", "",
		Call("", "fmt.Println", String("hello")),
	))

	var buf bytes.Buffer
	assert.Nil(t, f.WriteTo(&buf))
	assert.Equal(t, buf.String(), "// generated code\npackage main\n\nimport \"fmt\"\nfunc main() {\nfmt.Println(\"hello\")\n}\n")
}
