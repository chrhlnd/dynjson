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

Stats off my meger vm

```
go test --bench=".*"
PASS
BenchmarkResolveChildStr         1000000              1385 ns/op
BenchmarkResolveChildU64         1000000              1298 ns/op
BenchmarkResolveChildBool        1000000              1261 ns/op
BenchmarkResolve2ChildF64         200000              9952 ns/op
BenchmarkResolve2ChildString      200000              9743 ns/op
BenchmarkResolve3ChildArrayStr    500000              4375 ns/op
BenchmarkResolve2ChildCached     2000000               899 ns/op
BenchmarkResolve4ChildStr         100000             16996 ns/op
BenchmarkResolve4ChildCachedStr  2000000               926 ns/op
ok      github.com/chrhlnd/dynjson      17.753s
```
