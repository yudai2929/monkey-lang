package object

import (
	"bytes"
	"fmt"
	"gtihub.com/yudai2929/monkey-lang/ast"
)

// ObjectType is the type of the object
type ObjectType string

const (
	// INTEGER_OBJ is the integer object type
	INTEGER_OBJ = "INTEGER"
	// BOOLEAN_OBJ is the boolean object type
	BOOLEAN_OBJ = "BOOLEAN"
	// NULL_OBJ is the null object type
	NULL_OBJ = "NULL"
	// RETURN_VALUE_OBJ is the return value object type
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	// ERROR_OBJ is the error object type
	ERROR_OBJ = "ERROR"
	// FUNCTION_OBJ is the function object type
	FUNCTION_OBJ = "FUNCTION"
)

// Object is the interface that all objects in the interpreter implement
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer is the integer object
type Integer struct {
	Value int64
}

// Type returns the type of the object
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Inspect returns the string representation of the object
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean is the boolean object
type Boolean struct {
	Value bool
}

// Type returns the type of the object
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Inspect returns the string representation of the object
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null is the null object
type Null struct{}

// Type returns the type of the object
func (n *Null) Type() ObjectType { return NULL_OBJ }

// Inspect returns the string representation of the object
func (n *Null) Inspect() string { return "null" }

// ReturnValue is the return value object
type ReturnValue struct {
	Value Object
}

// Type returns the type of the object
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Inspect returns the string representation of the object
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error is the error object
type Error struct {
	Message string
}

// Type returns the type of the object
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Inspect returns the string representation of the object
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function is the function object
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type returns the type of the object
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

// Inspect returns the string representation of the object
func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(fmt.Sprintf("%s", params))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
