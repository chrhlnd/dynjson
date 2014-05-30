package dynjson

import (
	"encoding/json"
	"strconv"
)

func New()(DynNode) {
	return &dynnode{}
}

func NewFromBytes(bytes []byte)(DynNode) {
	return &dynnode{ data: bytes }
}

type dynnode struct {
	parent *dynnode
	path string
	data json.RawMessage
	ary []json.RawMessage
	obj map[string]json.RawMessage
}

func (n *dynnode)Path()(string) {
	return n.path
}

func (n *dynnode)Parent()(DynNode) {
	return n.parent
}

func (n *dynnode)Root()(DynNode) {
	if n.parent != nil {
		return n.parent.Root()
	}
	return n
}

var null_node dynnode

func (n *dynnode)Data()([]byte) {
	return n.data
}

func (n *dynnode)SetVal( v interface{} )(err error) {
	n.data,err = json.Marshal(v)
	n.ary = nil
	n.obj = nil
	return
}

func (n *dynnode)Resolve(path string)(DynNode,error) {
	if n.IsNull() {
		return &null_node,nil
	}

	makeArrayNode := func ( key string, idx int )(*dynnode,error) {
		var ary []json.RawMessage
		var err error

		if ary, err = n.Ary(); err != nil {
			return &null_node, err
		}

		if idx < 0 || ary == nil || idx >= len(ary) {
			return &null_node,nil
		}

		return &dynnode{
			parent : n,
			path : key,
			data : ary[idx],
		},nil
	}

	makeObjNode := func ( key string )(*dynnode,error) {
		var obj map[string]json.RawMessage
		var err error
		if obj, err = n.Obj(); err != nil {
			return &null_node, err
		}

		if obj == nil || obj[key] == nil {
			return &null_node,nil
		}

		return &dynnode{
			parent : n,
			path : key,
			data : n.obj[key],
		},nil
	}

	makeNode := func ( key string )(*dynnode,error) {
		idx, numerr := strconv.ParseInt(key,10,32)
		if numerr == nil {
			return makeArrayNode( key, int(idx) )
		} else {
			return makeObjNode( key )
		}
	}

	tok := path[0]

	for i := 1; i < len(path); i++ {
		if path[i] == tok {
			var child *dynnode
			var err error
			if child, err = makeNode(path[1:i]); err != nil {
				return child, err
			}

			if child.IsNull() {
				return child,nil
			}
			return child.Resolve(path[i:])
		}
	}
	return makeNode(path[1:])
}

func (n *dynnode)Obj()(map[string]json.RawMessage,error) {
	if n.obj != nil { return n.obj, nil }

	if err := json.Unmarshal( n.data, &n.obj ); err != nil {
		return nil, err
	}
	return n.obj, nil
}

func (n *dynnode)Ary()( ret []json.RawMessage, err error) {
	if n.ary != nil { ret = n.ary; return }

	if err = json.Unmarshal( n.data, &n.ary ); err != nil {
		return
	}

	ret = n.ary
	return
}

func (n *dynnode)IsNull()(bool) {
	return n == &null_node
}

func (n *dynnode)U64()(ret uint64,err error) {
	err = json.Unmarshal( n.data, &ret )
	return
}

func (n *dynnode)I64()(ret int64,err error) {
	err = json.Unmarshal( n.data, &ret )
	return
}

func (n *dynnode)F64()(ret float64, err error) {
	err = json.Unmarshal( n.data, &ret )
	return
}

func (n *dynnode)Str()(ret string, err error) {
	err = json.Unmarshal( n.data, &ret )
	return
}

func (n *dynnode)Bool()(ret bool,err error) {
	err = json.Unmarshal( n.data, &ret )
	return
}

