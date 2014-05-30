package dynjson

import (
	"testing"
)

var customer = `
{
	"name" : "fred"
	,"age" : 33
	,"fun" : false
	,"friends" : [
		"joe"
		,"douge"
		,"bob"
	]
	,"sister" : {
		"name" : "joey"
		,"rating" : 10.34
		,"fun" : true
		,"friends" : [ "joe", "dag", { "name" : "mum" } ]
	}
}
`
func BenchmarkResolveChildStr(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Resolve("/name"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Str(); err != nil || val != "fred" {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkResolveChildU64(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Resolve("/age"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.U64(); err != nil || val != 33 {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkResolveChildBool(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Resolve("/fun"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Bool(); err != nil || val != false {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkResolve2ChildF64(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Resolve("/sister/rating"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.F64(); err != nil || val != 10.34 {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkResolve2ChildString(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Resolve("/sister/name"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Str(); err != nil || val != "joey" {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkResolve3ChildArrayStr(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Resolve("/friends/2"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Str(); err != nil || val != "bob" {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkResolve2ChildCached(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	var node DynNode
	var err error
	if node, err = root.Resolve("/friends/0"); err != nil {
		b.Fatal("Failed")
	}
	for i := 0; i < b.N; i++ {
		if val, err := node.Str(); err != nil || val != "joe" {
			b.Fatal("Failed")
		}
	}
}

func BenchmarkResolve4ChildStr(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Resolve("/sister/friends/2/name"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Str(); err != nil || val != "mum" {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkResolve4ChildCachedStr(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	var node DynNode
	var err error

	if node, err = root.Resolve("/sister/friends/2/name"); err != nil {
		b.Fatal("Failed")
	}

	for i := 0; i < b.N; i++ {
		if val, err := node.Str(); err != nil || val != "mum" {
			b.Fatal("Failed")
		}
	}
}

func TestResolve(t *testing.T) {
	b := []byte(customer)

	root := NewFromBytes(b)

	var err error
	var node DynNode

	if node, err = root.Resolve("/name"); err == nil && !node.IsNull() {
		if str, err := node.Str(); str != "fred" {
			t.Errorf("Expected 'fred' for /name got [%v] err", str, err)
		} else {
			t.Logf("/name = %v", str)
		}

	} else {
		t.Errorf("Resolve failed to find /name expected DynNode value = 'fred' err was %v",err);
	}

	if node, err = root.Resolve("/friends"); err == nil && !node.IsNull() {
		if ary, err := node.Ary(); err != nil {
			t.Errorf("Expected Array of friends got error %v", err)
		} else {
			t.Logf("Friends = %v", ary)
			if len(ary) != 3 {
				t.Errorf("Expected 3 friends got %v", len(ary))
			}
		}

	} else {
		t.Errorf("Resolve failed to find /friends expected DynNode value = [ '', '', '' ] err was %v",err);
	}

	if node, err = root.Resolve("/friends/1"); err == nil && !node.IsNull() {
		if str, err := node.Str(); err != nil || str != "douge" {
			t.Errorf("Expected Middle friend of 'douge' got error %v", err)
		} else {
			t.Logf("Found 2nd friend named %v", str)
		}

	} else {
		t.Errorf("Resolve failed to find /friends/1 expected DynNode value = 'douge' err was %v",err);
	}

	if node, err = root.Resolve("/sister/rating"); err == nil && !node.IsNull() {
		if f64, err := node.F64(); err != nil || f64 != 10.34 {
			t.Errorf("Expected float got error %v",err)
		} else {
			t.Logf("Found rating %v", f64)
		}
	} else {
		t.Errorf("Failed err %v", err)
	}

	if node, err = root.Resolve("/sister/fun"); err == nil && !node.IsNull() {
		if fun, err := node.Bool(); err != nil || !fun {
			t.Errorf("Expected true got %v error %v",fun,err)
		}
	} else {
		t.Errorf("Failed err %v", err)
	}

	type Friend struct {
		Name string
		Age int
		Rank float32
	}

	type MoreData struct {
		Name string
		Ranks []int
		Friends []Friend
	}

	d := MoreData{
		Name: "hello",
		Ranks: []int{ 0, 32, 6, 6, },
		Friends: []Friend{ {"dude",6,33.34},{"mecka",33,133.88},{"beeber",338,13.8}},
	}
	root = New()
	root.SetVal(d)

	if node, err = root.Resolve("/Friends/0/Rank"); err != nil {
		t.Fatal(err)
	}

	var f64 float64
	if f64, err = node.F64(); err != nil {
		t.Fatal(err)
	}

	if f64 != 33.34 {
		t.Fatalf("Wrong val got %v",f64)
	}



}

