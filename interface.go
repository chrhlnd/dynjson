package dynjson

import "encoding/json"

type DynNode interface {
	LocalPath() string
	FullPath() string
	Parent() DynNode
	Root() DynNode

	//Path format is: first rune is the delimter, between delimiters are key names
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
	Obj() (map[string]json.RawMessage, error)
	Ary() ([]json.RawMessage, error)

	AsU64() uint64  // panic
	AsI64() int64   // panic
	AsF64() float64 // panic
	AsStr() string  // panic
	AsBool() bool   // panic

	SetVal(v interface{}) error
	Data() []byte
}
