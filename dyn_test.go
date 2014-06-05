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
func BenchmarkNodeChildStr(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Node("/name"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Str(); err != nil || val != "fred" {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkNodeChildU64(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Node("/age"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.U64(); err != nil || val != 33 {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkNodeChildBool(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Node("/fun"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Bool(); err != nil || val != false {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkNode2ChildF64(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Node("/sister/rating"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.F64(); err != nil || val != 10.34 {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkNode2ChildString(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Node("/sister/name"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Str(); err != nil || val != "joey" {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkNode3ChildArrayStr(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Node("/friends/2"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Str(); err != nil || val != "bob" {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkNode2ChildCached(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	var node DynNode
	var err error
	if node, err = root.Node("/friends/0"); err != nil {
		b.Fatal("Failed")
	}
	for i := 0; i < b.N; i++ {
		if val, err := node.Str(); err != nil || val != "joe" {
			b.Fatal("Failed")
		}
	}
}

func BenchmarkNode4ChildStr(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if node, err := root.Node("/sister/friends/2/name"); err != nil {
			b.Fatal("Failed")
		} else {
			if val, err := node.Str(); err != nil || val != "mum" {
				b.Fatal("Failed")
			}
		}
	}
}

func BenchmarkNode4ChildCachedStr(b * testing.B) {
	root := NewFromBytes([]byte(customer))

	b.ResetTimer()

	var node DynNode
	var err error

	if node, err = root.Node("/sister/friends/2/name"); err != nil {
		b.Fatal("Failed")
	}

	for i := 0; i < b.N; i++ {
		if val, err := node.Str(); err != nil || val != "mum" {
			b.Fatal("Failed")
		}
	}
}

func BenchmarkMutate1stChild(b *testing.B) {
	root := NewFromBytes([]byte(customer))
	for i := 0; i < b.N; i++ {
		n := root.AsNode("/name");
		n.SetVal(68.33)
	}
}

func BenchmarkMutate2rdChild(b *testing.B) {
	root := NewFromBytes([]byte(customer))
	for i := 0; i < b.N; i++ {
		n := root.AsNode("/sister/name");
		n.SetVal(300)
	}
}

func BenchmarkMutate3rdChild(b *testing.B) {
	root := NewFromBytes([]byte(customer))


	for i := 0; i < b.N; i++ {
		n := root.AsNode("/sister/friends/0");
		n.SetVal(true)
	}
}

func TestMutate(t *testing.T) {
	root := NewFromBytes([]byte(customer))

	var node DynNode
	var node1 DynNode
	var err error

	errIf := func (err error) {
		if err != nil {
			t.Error(err)
			panic("blah")
		}
	}

	node, err = root.Node("/name");
	errIf(err);

	node.SetVal("Sup")


	node1, err = root.Node("/name");
	errIf(err);

	if node1.AsStr() != node.AsStr() {
		t.Errorf("Mutation failed expected %v, got %v\n", node.AsStr(), node1.AsStr())
		return
	}

	node, err = root.Node("/friends/0")
	errIf(err)

	node.SetVal("MONKEYMAN")

	t.Logf("Data was %v\n", string(root.Data()))

	node1, err = root.Node("/friends/0")
	errIf(err)

	if node1.AsStr() != node.AsStr() {
		t.Errorf("Mutation failed expec %v got %v\n", node.AsStr(), node1.AsStr())
		return
	}

	node, err = root.Node("/sister/friends/1")
	errIf(err)

	node.SetVal("SUPERMONKEYMAN")

	if root.AsNode("/sister/friends/1").AsStr() != "SUPERMONKEYMAN" {
		t.Error("Mutate failed")
		return
	}

	root.AsNode("/sister/name").SetVal(33.65)
	t.Logf("Data was %v\n", string(root.Data()))

	if root.AsNode("/sister/name").AsF64() != 33.65 {
		t.Error("F64 fail")
		return
	}


	t.Logf("Data was %v\n", string(root.Data()))

}

func TestCover(t *testing.T) {
	root := NewFromBytes([]byte(customer))

	if root.AsNode("/age").AsU64() != 33 {
		t.Error("failed")
	}

	if root.AsNode("/age").AsI64() != 33 {
		t.Error("failed")
	}

	if root.AsNode("/sister/rating").AsF64() != 10.34 {
		t.Error("failed")
	}

	if root.AsNode("/name").AsStr() != "fred" {
		t.Error("failed")
	}

	if root.AsNode("/fun").AsBool() != false {
		t.Error("failed")
	}

	if root.AsNode("/notthere").IsNull() != true {
		t.Error("failed")
	}

	ageN := root.AsNode("/age")
	if ageN.IsNull() {
		t.Error("failed")
	}

	if err := ageN.SetVal(34); err != nil {
		t.Error("failed")
	}

	if root.AsNode("/age").AsU64() != 34 {
		t.Error("failed")
	}

	sisF := root.AsNode("/sister/friends")
	if sisF.IsNull() {
		t.Error("failed")
	}
	friend2 := sisF.AsNode("/2")
	if friend2.IsNull() {
		t.Error("failed")
	}

	if err := sisF.SetVal( []string { "jim", "bob" } ); err != nil {
		t.Error(err)
	}

	if friend2.Root() != friend2 {
		t.Error("failed orphan")
	}

	friend1 := sisF.AsNode("/1")
	if friend1.IsNull() {
		t.Error("failed")
	}

	if friend1.AsStr() != "bob" {
		t.Error("failed mutation")
	}
}

func TestNode(t *testing.T) {
	b := []byte(customer)

	root := NewFromBytes(b)

	var err error
	var node DynNode

	if node, err = root.Node("/name"); err == nil && !node.IsNull() {
		if str, err := node.Str(); str != "fred" {
			t.Errorf("Expected 'fred' for /name got [%v] err", str, err)
		} else {
			t.Logf("/name = %v", str)
		}

	} else {
		t.Errorf("Node failed to find /name expected DynNode value = 'fred' err was %v",err);
	}

	if node, err = root.Node("/friends"); err == nil && !node.IsNull() {
		if ary, err := node.Ary(); err != nil {
			t.Errorf("Expected Array of friends got error %v", err)
		} else {
			t.Logf("Friends = %v", ary)
			if len(ary) != 3 {
				t.Errorf("Expected 3 friends got %v", len(ary))
			}
		}

	} else {
		t.Errorf("Node failed to find /friends expected DynNode value = [ '', '', '' ] err was %v",err);
	}

	if node, err = root.Node("/friends/1"); err == nil && !node.IsNull() {
		if str, err := node.Str(); err != nil || str != "douge" {
			t.Errorf("Expected Middle friend of 'douge' got error %v", err)
		} else {
			t.Logf("Found 2nd friend named %v", str)
		}

	} else {
		t.Errorf("Node failed to find /friends/1 expected DynNode value = 'douge' err was %v",err);
	}

	if node, err = root.Node("/sister/rating"); err == nil && !node.IsNull() {
		if f64, err := node.F64(); err != nil || f64 != 10.34 {
			t.Errorf("Expected float got error %v",err)
		} else {
			t.Logf("Found rating %v", f64)
		}
	} else {
		t.Errorf("Failed err %v", err)
	}

	if node, err = root.Node("/sister/fun"); err == nil && !node.IsNull() {
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

	if node, err = root.Node("/Friends/0/Rank"); err != nil {
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

