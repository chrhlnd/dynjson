dynjson
=======

Small library to allow for dynamic JSON access in GO


API is make a node object with the encoded data. Call Resolve to get things out. Convert to type with Convert methods.
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
  
  if node, err = dyn.Resolve("/some/path/that/is/cool"); err != nil {
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
See the dyn_test.go file for more usage. You can SetVal too and it will serialize into bytes for you.

When you call resolve it, dynamiclly decodes the path. So if you have some big object buried in the hierarchy you should save out
the parent node then query off that.

congress_books, err = dyn.Resolve("/library/congress/collection")

The array object is cached in the node, so it doesn't have to reparse that.
```
book_name err = congress_books.Resolve("/3000/name")
book_name, err = congress_books.Resolve("/3001/name") 
```
Better then digging through all objects every time.

Objects are pathed by their key name, arrays are pathed by their index number.
```
node.Resolve("/name/of/keys/in/json/object")
node.Resolve("/json/array/500")

```

