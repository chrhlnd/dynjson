dynjson
=======

Small library to allow for dynamic JSON access in GO


Make a node object with the encoded []byte data.

Call Node to get things out.

Convert to type with Convert methods.

```
import (
  "github.com/chrhlnd/dynjson"
  "os"
  "log"
)

func main() {
  dyn := dynjson.NewFromBytes([]byte(data))
  
  var node dynjson.DynNode
  var err error
  
  if node, err = dyn.Node("/some/path/that/is/cool"); err != nil {
    log.Errorf("My path didn't parse error %v", err)
    os.Exit(1)
  }
  
  var val string
  if val, err = node.Str(); err != nil {
    log.Errorf("My value wasn't a string? %v", err)
    os.Exit(1)
  }
}
```
See the dyn_test.go file for more usage.

n.SetVal will serialize into bytes for you, recomposing data up the Parent chain.

When you call Node (AsNode), it dynamiclly decodes the path. So if you have some big object buried in the hierarchy you should save out the parent node then query off that.

```
congress_books, err = dyn.Node("/library/congress/collection")
```

The array object is cached in the node, so it doesn't have to reparse that.
```
book_name, err = congress_books.Node("/3000/name")
book_name, err = congress_books.Node("/3001/name") 
```
Better then digging through all objects every time.

Objects are addressed by their key name.
Arrays are addressed by their index number.

```
node.Node("/name/of/keys/in/json/object")
node.Node("/json/array/500")

```

The As'method's are for convience and panicing.


```
node.AsNode("/name/of/val").AsI64()
```

Will panic if path doesn't exist and if the value isn't parsable as I64

Mutation is handeled by SetVal, this back propagates in the doc tree, possibly nullifying other nodes if the data set doesn't contain them anymore. This is done with a version number that is incremented to the root. If an orphaned node is accessed it will find its not the right version then possibly orphan itself upon access.

```
{
	"name" : "parent"
	"children" : [ "zero", "one", "two" ]
}
```

```
child_zero := root.AsNode("/children/0")
child_two := root.AsNode("/children/2")

root.AsNode("/children").SetVal( []string{ "zero", "one" } )

if child_two.Root() == child_two {
	fmt.Printf("I'm an orphan")
}

if child_zero.Root() == root {
	fmt.Printf("I'm still in the graph whee")
}
```

=======
