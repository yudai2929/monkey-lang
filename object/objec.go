package object

import (
	"bytes"
	"fmt"
	"gtihub.com/yudai2929/monkey-lang/ast"
	"hash/fnv"
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
	// STRING_OBJ is the string object type
	STRING_OBJ = "STRING"
	// BUILTIN_OBJ is the built-in function object type
	BUILTIN_OBJ = "BUILTIN"
	// ARRAY_OBJ is the array object type
	ARRAY_OBJ = "ARRAY"
	// HASH_OBJ is the hash object type
	HASH_OBJ = "HASH"
)

// HashKey is the hash key object
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Hashable is the interface that all hashable objects in the interpreter implement
type Hashable interface {
	HashKey() HashKey
}

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

// HashKey returns the hash key of the object
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Boolean is the boolean object
type Boolean struct {
	Value bool
}

// Type returns the type of the object
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Inspect returns the string representation of the object
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// HashKey returns the hash key of the object
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

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

// String is the string object
type String struct {
	Value string
}

// Type returns the type of the object
func (s *String) Type() ObjectType { return STRING_OBJ }

// Inspect returns the string representation of the object
func (s *String) Inspect() string { return s.Value }

// HashKey returns the hash key of the object
func (s *String) HashKey() HashKey {
	var h = fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// BuiltinFunctionType is the type of the built-in function
type BuiltinFunctionType func(args ...Object) Object

// Builtin is the built-in function object
type Builtin struct {
	Fn BuiltinFunctionType
}

// Type returns the type of the object
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

// Inspect returns the string representation of the object
func (b *Builtin) Inspect() string { return "builtin function" }

// Array is the array object
type Array struct {
	Elements []Object
}

// Type returns the type of the object
func (ao *Array) Type() ObjectType { return ARRAY_OBJ }

// Inspect returns the string representation of the object
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, el := range ao.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(fmt.Sprintf("%s", elements))
	out.WriteString("]")

	return out.String()
}

// HashPair is the key-value pair of the hash object
type HashPair struct {
	Key   Object
	Value Object
}

// Hash is the hash object
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns the type of the object
func (h *Hash) Type() ObjectType { return HASH_OBJ }

// Inspect returns the string representation of the object
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(fmt.Sprintf("%s", pairs))
	out.WriteString("}")

	return out.String()
}
