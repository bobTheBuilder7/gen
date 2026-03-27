package gen

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

type File struct {
	packageName string
	comment     string
	blocks      []fmt.Stringer
	mu          sync.RWMutex
}

type imprt struct {
	alias string
	path  string
}

func (i imprt) String() string {
	return fmt.Sprintf("import %s\"%s\"", i.alias, i.path)
}

type vr struct {
	name  string
	value fmt.Stringer
}

func (v vr) String() string {
	return fmt.Sprintf("var %s = %s", v.name, v.value.String())
}

type cnst struct {
	name  string
	value fmt.Stringer
}

func (c cnst) String() string {
	return fmt.Sprintf("const %s = %s", c.name, c.value.String())
}

func Import(alias, path string) imprt {
	return imprt{alias, path}
}

func Var(name string, value fmt.Stringer) vr {
	return vr{name, value}
}

func Const(name string, value fmt.Stringer) cnst {
	return cnst{name, value}
}

func NewFile(packageName string, comments ...string) *File {
	file := new(File)

	file.packageName = packageName
	file.comment = strings.Join(comments, "\n")

	return file
}

type raw string

func (r raw) String() string {
	return string(r)
}

func String(v string) fmt.Stringer {
	return raw(fmt.Sprintf("\"%s\"", v))
}

func Int(v int) fmt.Stringer {
	return raw(fmt.Sprintf("%d", v))
}

func Bool(v bool) fmt.Stringer {
	return raw(fmt.Sprintf("%t", v))
}

func Float(v float64) fmt.Stringer {
	return raw(fmt.Sprintf("%f", v))
}

func Line(s string) fmt.Stringer {
	return raw(s)
}

func Rune(v rune) fmt.Stringer {
	return raw(fmt.Sprintf("'%c'", v))
}

func Arg(v string) fmt.Stringer {
	return raw(v)
}

func Call(assigns, name string, args ...fmt.Stringer) fmt.Stringer {
	call := fmt.Sprintf("%s(%s)", name, join(args, ", "))
	if assigns != "" {
		call = fmt.Sprintf("%s := %s", assigns, call)
	}
	return raw(call)
}

func ErrCheck(value fmt.Stringer) fmt.Stringer {
	if value.String() != "" {
		return raw(fmt.Sprintf("if err != nil {\nreturn %s, err\n}", value.String()))
	}
	return raw("if err != nil {\nreturn err\n}")
}

type typeAlias struct {
	name       string
	underlying string
}

func (t typeAlias) String() string {
	return fmt.Sprintf("type %s %s", t.name, t.underlying)
}

func TypeAlias(name, underlying string) typeAlias {
	return typeAlias{name, underlying}
}

type Field struct {
	Name string
	Type string
	Tag  string
}

func (f Field) String() string {
	if f.Tag != "" {
		return fmt.Sprintf("%s %s `%s`", f.Name, f.Type, f.Tag)
	}
	return fmt.Sprintf("%s %s", f.Name, f.Type)
}

type strct struct {
	name   string
	fields []Field
}

func (s strct) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("type %s struct {\n", s.name))
	for _, f := range s.fields {
		out.WriteString(f.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

func Struct(name string, fields ...Field) strct {
	return strct{name, fields}
}

type Method struct {
	Name    string
	Params  string
	Returns string
}

func (m Method) String() string {
	return fmt.Sprintf("%s(%s) %s", m.Name, m.Params, m.Returns)
}

type iface struct {
	name    string
	methods []Method
}

func (i iface) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("type %s interface {\n", i.name))
	for _, m := range i.methods {
		out.WriteString(m.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

func Interface(name string, methods ...Method) iface {
	return iface{name, methods}
}

type fn struct {
	name     string
	receiver string
	params   string
	returns  string
	body     []fmt.Stringer
}

func (f fn) String() string {
	var out strings.Builder
	out.WriteString("func ")
	if f.receiver != "" {
		out.WriteString(fmt.Sprintf("(%s) ", f.receiver))
	}
	out.WriteString(fmt.Sprintf("%s(%s)", f.name, f.params))
	if f.returns != "" {
		out.WriteString(" " + f.returns)
	}
	out.WriteString(" {\n")
	for _, stmt := range f.body {
		out.WriteString(stmt.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

func Func(name, params, returns string, body ...fmt.Stringer) fn {
	return fn{name: name, params: params, returns: returns, body: body}
}

func MethodFunc(receiver, name, params, returns string, body ...fmt.Stringer) fn {
	return fn{name: name, receiver: receiver, params: params, returns: returns, body: body}
}

func (f *File) AddBlock(block fmt.Stringer) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blocks = append(f.blocks, block)
}

func (f *File) WriteTo(w io.Writer) (int64, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var total int64

	if f.comment != "" {
		for line := range strings.SplitSeq(f.comment, "\n") {
			n, err := fmt.Fprintf(w, "// %s\n", line)
			total += int64(n)
			if err != nil {
				return total, err
			}
		}
	}

	n, err := fmt.Fprintf(w, "package %s\n\n", f.packageName)
	total += int64(n)
	if err != nil {
		return total, err
	}

	for _, b := range f.blocks {
		n, err = fmt.Fprintf(w, "%s\n", b.String())
		total += int64(n)
		if err != nil {
			return total, err
		}
	}

	return total, nil
}
