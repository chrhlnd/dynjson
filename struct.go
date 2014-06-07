package dynjson

import (
	"encoding/json"
	"strconv"
	"strings"
	"bytes"
//	"log"
)

func New()(DynNode) {
	return &dynnode{}
}

func NewFromBytes(bytes []byte)(DynNode) {
	return &dynnode{ data: bytes }
}

type dynnode struct {
	version int
	parent *dynnode
	path string
	data json.RawMessage
	ary []json.RawMessage
	obj map[string]json.RawMessage
}

func (n *dynnode)LocalPath()(string) {
	return n.path
}

func (n *dynnode)buildPath( paths []string)(string) {
	paths = append(paths, n.path)
	if n.parent != nil {
		return n.parent.buildPath(paths)
	}
	return strings.Join( paths, "" )
}

func (n *dynnode)checkVersion( v int )(bool) {
	if n.version == v {
		if n.parent != nil {
			return n.parent.checkVersion(v)
		}
		return true
	}
	return false
}

func (n *dynnode)inVersion()(bool) {
	return n.checkVersion( n.version )
}

func (n *dynnode)FullPath()(string) {
	n.syncVersion()
	return n.buildPath(make([]string,0,3))
}

func (n *dynnode)Parent()(DynNode) {
	n.syncVersion()
	return n.parent
}

func (n *dynnode)Root()(DynNode) {
	n.syncVersion()
	if n.parent != nil {
		return n.parent.Root()
	}
	return n
}

var null_node dynnode

func (n *dynnode)Data()([]byte) {
	n.syncVersion()
	return n.data
}

func (n *dynnode)incVersion() {
	n.version++
	n.syncVersion()
}

func (n *dynnode)syncChild( child string, child_data []byte )(int) {
	idx, numerr := strconv.ParseInt(child[1:],10,32)
	if numerr == nil {
		var err error
		var ary []json.RawMessage

		if ary, err = n.Ary(); err != nil {
			panic(err) // invalid state, if child is a number its parent better be an array
		}

		buf := new (bytes.Buffer)
		buf.Grow(len(n.data)-len(ary[idx])+len(child_data))

		ary[idx] = child_data

		buf.WriteString("[")
		for _, v := range ary {
			buf.Write(v)
			buf.WriteString(",")
		}
		buf.Truncate(buf.Len()-len(","))
		buf.WriteString("]")

		n.data = buf.Bytes()

		if n.parent != nil {
			n.parent.syncChild(n.path, n.data)
			n.syncVersion()
		} else {
			n.incVersion()
		}

		return n.version
	} else {
		var err error
		var obj map[string]json.RawMessage

		if obj, err = n.Obj(); err != nil {
			panic(err)
		}

		buf := new(bytes.Buffer)
		buf.Grow(len(n.data)-len(obj[child])+len(child_data))

		obj[child[1:]] = child_data

		buf.WriteString("{")
		for k, v := range obj {
			buf.WriteString("\"")
			buf.WriteString(k)
			buf.WriteString("\":")
			buf.Write(v)
			buf.WriteString(",")
		}
		buf.Truncate(buf.Len()-len(","))
		buf.WriteString("}")

		n.data = buf.Bytes()

		if n.parent != nil {
			n.parent.syncChild(n.path, n.data)
			n.syncVersion()
		} else {
			n.incVersion()
		}

		return n.version
	}
}

func (n *dynnode)syncVersion() {
	if !n.inVersion() {
		node, err := n.parent.Node(n.path)
		if err != nil {
			panic(err) // mal formed data
		}

		if node.IsNull() {
			// orphaned
			n.parent = nil
		} else {
			// mutated
			n.copyFrom(node)
		}
	}
}

func (n *dynnode)copyFrom( other DynNode ) {
	switch o := other.(type) {
	case *dynnode:
		*n = *o
	default:
		panic("maybe the copy interface should be public")
	}
}

func (n *dynnode)SetVal( v interface{} )(err error) {
	n.data,err = json.Marshal(v)
	n.ary = nil
	n.obj = nil
	if n.parent != nil {
		//log.Printf("Syncing parent with new data @ path %v", n.path)
		n.parent.syncChild(n.path, n.data)
		n.syncVersion()
	} else {
		//log.Printf("Syncing no parent incing version")
		n.incVersion()
	}
	return
}

func (n *dynnode)AsNode(path string)(r DynNode) {
	r, err := n.Node(path)
	if err != nil { panic(err)	}
	return
}

func makeArrayNode(n *dynnode, key string, idx int)(*dynnode,error) {
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

func makeObjNode(n *dynnode,key string)(*dynnode,error) {
	var obj map[string]json.RawMessage
	var err error
	if obj, err = n.Obj(); err != nil {
		return &null_node, err
	}

	if obj == nil || obj[key[1:]] == nil {
		return &null_node,nil
	}

	ret := &dynnode{
		parent : n,
		path : key,
		data : n.obj[key[1:]],
	}
	return ret, nil
}

func makeNode(n *dynnode,key string)(*dynnode,error) {
	idx, numerr := strconv.ParseInt(key[1:],10,32)
	if numerr == nil {
		return makeArrayNode( n, key, int(idx) )
	} else {
		return makeObjNode( n, key )
	}
}

func (n *dynnode)Node(path string)(DynNode,error) {
	n.syncVersion()
	if n.IsNull() {
		return &null_node,nil
	}

	tok := path[0]

	for i := 1; i < len(path); i++ {
		if path[i] == tok {
			var child *dynnode
			var err error
			if child, err = makeNode(n,path[0:i]); err != nil {
				return child, err
			}

			if child.IsNull() {
				return child,nil
			}
			return child.Node(path[i:])
		}
	}
	return makeNode(n,path)
}

func (n *dynnode)Obj()(map[string]json.RawMessage,error) {
	n.syncVersion()
	if n.obj != nil { return n.obj, nil }

	if err := json.Unmarshal( n.data, &n.obj ); err != nil {
		return nil, err
	}
	return n.obj, nil
}

func (n *dynnode)Ary()( ret []json.RawMessage, err error) {
	n.syncVersion()
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

func (n *dynnode)AsU64()(r uint64) {
	if n.IsNull() { panic("path NULL"); }
	r, err := n.U64()
	if err != nil { panic(err) }
	return
}

func (n *dynnode)AsI64()(r int64) {
	if n.IsNull() { panic("path NULL"); }
	r, err := n.I64()
	if err != nil { panic(err) }
	return
}

func (n *dynnode)AsF64()(r float64) {
	if n.IsNull() { panic("path NULL"); }
	r, err := n.F64()
	if err != nil { panic(err) }
	return
}

func (n *dynnode)AsStr()(r string) {
	if n.IsNull() { panic("path NULL"); }
	r, err := n.Str()
	if err != nil { panic(err) }
	return
}

func (n *dynnode)AsBool()(r bool) {
	if n.IsNull() { panic("path NULL"); }
	r, err := n.Bool()
	if err != nil { panic(err) }
	return
}

