package dynjson

import "encoding/json"

type DynNode interface {
	LocalPath() string
	FullPath() string
	Parent() DynNode
	Root() DynNode

	//Node Path format is: first rune is the delimter, between delimiters are key names
	//  Arrays use Numeric keys
	//         #0#customer#name
	//         /books/0/name
	//
	Node(path string) (DynNode, error)
	AsNode(path string) DynNode // panic

	IsNull() bool

	U64() (uint64, error)
	I64() (int64, error)
	F64() (float64, error)
	Str() (string, error)
	Bool() (bool, error)
	// Struct Calls json.Unmarshal for you
	Struct(out interface{}) error

	//Obj this node is an object get the cached object
	Obj() (map[string]json.RawMessage, error)
	// Ary if this node is an array get the cached object
	Ary() ([]json.RawMessage, error)

	// Len will return the count of elements if this is a collection type node
	Len() int

	AsU64() uint64  // panic
	AsI64() int64   // panic
	AsF64() float64 // panic
	AsStr() string  // panic
	AsBool() bool   // panic

	SetVal(v interface{}) error
	Data() []byte
}
