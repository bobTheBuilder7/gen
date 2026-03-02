package main

import (
	"os"

	"github.com/bobTheBuilder7/gen"
)

func main() {
	f := gen.NewFile("gen")

	f.AddBlock(gen.Import("", "context"))
	f.AddBlock(gen.Import("", "net/http"))

	f.AddBlock(gen.Const("Version", gen.String("1.0.0")))
	f.AddBlock(gen.Const("MaxRetries", gen.Int(3)))
	f.AddBlock(gen.Const("Debug", gen.Bool(false)))
	f.AddBlock(gen.Const("Pi", gen.Float(3.14)))
	f.AddBlock(gen.Const("Separator", gen.Rune(':')))

	f.AddBlock(gen.Var("DefaultTimeout", gen.Int(30)))
	f.AddBlock(gen.Var("AppName", gen.String("myapp")))

	f.AddBlock(gen.TypeAlias("ID", "int"))
	f.AddBlock(gen.TypeAlias("Handler", "func(w http.ResponseWriter, r *http.Request)"))

	f.AddBlock(gen.Struct("User",
		gen.Field{Name: "ID", Type: "int", Tag: `json:"id" db:"id"`},
		gen.Field{Name: "Name", Type: "string", Tag: `json:"name"`},
		gen.Field{Name: "Email", Type: "string", Tag: `json:"email"`},
		gen.Field{Name: "Active", Type: "bool"},
	))

	f.AddBlock(gen.Struct("Config",
		gen.Field{Name: "Host", Type: "string", Tag: `json:"host" yaml:"host"`},
		gen.Field{Name: "Port", Type: "int", Tag: `json:"port"`},
		gen.Field{Name: "Debug", Type: "bool", Tag: `json:"debug"`},
	))

	f.AddBlock(gen.Interface("Repository",
		gen.Method{Name: "FindByID", Params: "id int", Returns: "(*User, error)"},
		gen.Method{Name: "Save", Params: "u *User", Returns: "error"},
		gen.Method{Name: "Delete", Params: "id int", Returns: "error"},
	))

	f.AddBlock(gen.Interface("Stringer",
		gen.Method{Name: "String", Returns: "string"},
	))

	f.AddBlock(gen.Const("GetUsersSQL", gen.String("select * from users")))

	f.AddBlock(gen.MethodFunc("q *Queries", "GetUsers", "ctx context.Context, id int", "error",
		gen.Call("err", "q.db.Exec", gen.Arg("ctx"), gen.Arg("GetUsersSQL")),
		gen.ErrCheck(gen.Arg("")),
		gen.Line("return nil"),
	))

	file, err := os.Create("./generated.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	f.WriteTo(file)
}
