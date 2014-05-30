package dynjson

import "encoding/json"

type DynNode interface {
	Path()(string)
	Parent()(DynNode)
	Root()(DynNode)

	//Path format is first char is delimter, between characters are keys
	//  Arrays use Numeric keys
	//         #0#customer#name
	//         /books/0/name
	//         
	Resolve( path string )(DynNode,error)

	IsNull()(bool)

	U64()(uint64,error)
	I64()(int64,error)
	F64()(float64,error)
	Str()(string,error)
	Bool()(bool,error)
	Obj()(map[string]json.RawMessage,error)
	Ary()([]json.RawMessage,error)

	SetVal(v interface{})(error)
	Data()([]byte)
}
