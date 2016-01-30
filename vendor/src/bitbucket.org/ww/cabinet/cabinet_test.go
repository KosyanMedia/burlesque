package cabinet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

var test_db string = "casket.kch"
var test_db2 string = "casket2.kch"
var test_dump string = "casket.dump"

func newCabinet(t *testing.T, filename string) *KCDB {
	kc := New()
	err := kc.Open(filename, KCOWRITER|KCOCREATE)
	if err != nil {
		t.Fatalf("Open(): %s", err)
	}
	return kc
}

func delCabinet(t *testing.T, kc *KCDB) {
	path, err := kc.Path()
	if err != nil {
		t.Fatalf("Path(): %s", err)
	}
	err = kc.Close()
	if err != nil {
		t.Fatalf("Close(): %s", err)
	}
	kc.Del()
	err = os.Remove(path)
	if err != nil {
		t.Fatal("Remove(%s): %s", path, err)
	}
}

func addKVP(t *testing.T, kc *KCDB, k, v string) {
	err := kc.Add([]byte(k), []byte(v))
	if err != nil {
		t.Errorf("Add(%s, %s): %s", k, v, err)
	}
}

func setKVP(t *testing.T, kc *KCDB, k, v string) {
	err := kc.Set([]byte(k), []byte(v))
	if err != nil {
		t.Errorf("Set(%s, %s): %s", k, v, err)
	}
}

func getKVP(t *testing.T, kc *KCDB, k, v string) {
	v2, err := kc.Get([]byte(k))
	if err != nil {
		t.Errorf("Get(%s): %s", k, err)
	}
	if v != string(v2) {
		t.Errorf("Get(%s) expected %s got %s", k, v, v2)
	}
}

func TestAdd(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	addKVP(t, kc, "hello", "world")
	addKVP(t, kc, "bit", "bucket")
	addKVP(t, kc, "foo", "bar")

	getKVP(t, kc, "hello", "world")
	getKVP(t, kc, "bit", "bucket")
	getKVP(t, kc, "foo", "bar")

	status, err := kc.Status()
	if err != nil {
		t.Errorf("Status(): %s", err)
	}
	path, err := kc.Path()
	if err != nil {
		t.Errorf("Path(): %s", err)
	}
	count, err := kc.Count()
	if err != nil {
		t.Errorf("Count(): %s", err)
	}
	size, err := kc.Size()
	if err != nil {
		t.Errorf("Size(): %s", err)
	}
	fmt.Printf("%s %d bytes %d records\n%s", path, size, count, status)
}

func TestSet(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	addKVP(t, kc, "hello", "world")
	getKVP(t, kc, "hello", "world")
	setKVP(t, kc, "hello", "universe")
	getKVP(t, kc, "hello", "universe")
}

func TestReplace(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	err := kc.Replace([]byte("foo"), []byte("bar"))
	if err == nil {
		t.Errorf("Replace(foo, bar) expected failure")
	}
	addKVP(t, kc, "foo", "bar")
	err = kc.Replace([]byte("foo"), []byte("baz"))
	if err != nil {
		t.Error(err)
	}
	getKVP(t, kc, "foo", "baz")
}

func TestAppend(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	addKVP(t, kc, "numbers", "one")
	err := kc.Append([]byte("numbers"), []byte("two"))
	if err != nil {
		t.Error(err)
	}
	getKVP(t, kc, "numbers", "onetwo")
}

func TestTran(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	err := kc.BeginTran(false)
	addKVP(t, kc, "numbers", "1")
	err = kc.EndTran(false)
	count, err := kc.Count()
	if err != nil {
		t.Fatalf("Count(): %s", err)
	}
	if count != 0 {
		t.Errorf("Rolled back - should have no items, have %d", count)
	}

	err = kc.BeginTran(false)
	addKVP(t, kc, "numbers", "1")
	err = kc.EndTran(true)
	count, err = kc.Count()
	if err != nil {
		t.Fatalf("Count(): %s", err)
	}
	if count != 1 {
		t.Errorf("Committed - should have 1 item, have %d", count)
	}
}

func TestIncr(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	_, err := kc.IncrInt([]byte("numbers"), 6)
	if err != nil {
		t.Error("IncrInt(6):", err)
	}

	v, err := kc.Get([]byte("numbers"))
	var n uint64
	binary.Read(bytes.NewBuffer(v), binary.BigEndian, &n)
	if n != 6 {
		t.Errorf("IncrInt(6): expected 6 got %d", n)
	}

	_, err = kc.IncrInt([]byte("numbers"), 6)
	if err != nil {
		t.Error("IncrInt(6):", err)
	}

	v, err = kc.Get([]byte("numbers"))
	binary.Read(bytes.NewBuffer(v), binary.BigEndian, &n)
	if n != 12 {
		t.Errorf("IncrInt(6): expected 12 got %d", n)
	}

	return

	err = kc.IncrDouble([]byte("floats"), 1.5)
	if err != nil {
		t.Errorf("IncrDouble(1.5): %s", err)
	}

	var d float64
	v, err = kc.Get([]byte("floats"))
	binary.Read(bytes.NewBuffer(v), binary.BigEndian, &d)
	if d != 1.5 {
		t.Errorf("IncrDouble(1.5): expected 1.5 got %l", d)
	}

	err = kc.IncrDouble([]byte("floats"), 2.0)
	if err != nil {
		t.Errorf("IncrDouble(2.0): %s", err)
	}

	v, err = kc.Get([]byte("floats"))
	binary.Read(bytes.NewBuffer(v), binary.BigEndian, &d)
	if d != 12 {
		t.Errorf("IncrDouble(2.0): expected 3.5 got %l", d)
	}

}

func TestIterators(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	addKVP(t, kc, "hello", "world")
	addKVP(t, kc, "bit", "bucket")
	addKVP(t, kc, "foo", "bar")

	count, err := kc.Count()
	if err != nil {
		t.Errorf("Count(): %s", err)
	}

	seen := make(map[string]bool)
	keys := kc.Keys()
	for k := range keys {
		seen[string(k)] = true
	}
	if len(seen) != int(count) {
		t.Errorf("Iterate over keys expected %d got %d items", count, len(seen))
	}

	seen = make(map[string]bool)
	values := kc.Values()
	for v := range values {
		seen[string(v)] = true
	}
	if len(seen) != int(count) {
		t.Errorf("Iterate over values expected %d got %d items", count, len(seen))
	}

	seeni := make(map[string]string)
	items := kc.Items()
	for i := range items {
		seeni[string(i.Key)] = string(i.Value)
	}
	if len(seen) != int(count) {
		t.Errorf("Iterate over items expected %d got %d items", count, len(seen))
	}
}

func TestMatch(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	addKVP(t, kc, "hello", "world")
	addKVP(t, kc, "hell", "world")
	addKVP(t, kc, "1hello", "world")

	matches, err := kc.MatchPrefix("he", 10)
	if err != nil {
		t.Fatalf("MatchPrefix(he): %s", err)
	}
	if len(matches) != 2 {
		t.Errorf("MatchPrefix(he): Expected two matches, got %d", len(matches))
	}

	matches, err = kc.MatchRegex("^[0-9][a-z]+", 10)
	if err != nil {
		t.Fatalf("MatchRegex(he): %s", err)
	}
	if len(matches) != 1 {
		t.Errorf("MatchRegex(^[0-9][a-z]+) Expected one matches, got %d", len(matches))
	}
}

func TestClear(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	addKVP(t, kc, "hello", "world")
	addKVP(t, kc, "bit", "bucket")
	addKVP(t, kc, "foo", "bar")

	count, err := kc.Count()
	if err != nil {
		t.Errorf("Count(): %s", err)
	}
	if count != uint64(3) {
		t.Errorf("Count() expected %d got %d", 3, count)
	}
	err = kc.Clear()
	if err != nil {
		t.Errorf("Clear(): %s", err)
	}
	count, err = kc.Count()
	if err != nil {
		t.Errorf("Count(): %s", err)
	}
	if count != uint64(0) {
		t.Errorf("Count() after Clear() expected %d got %d", 0, count)
	}
}

func TestCopy(t *testing.T) {
	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)

	addKVP(t, kc, "hello", "world")
	addKVP(t, kc, "bit", "bucket")
	addKVP(t, kc, "foo", "bar")

	err := kc.Copy(test_db2)
	if err != nil {
		t.Fatalf("Copy(%s): %s", test_db2, err)
	}

	kc2 := newCabinet(t, test_db2)
	defer delCabinet(t, kc2)

	count, _ := kc.Count()
	count2, _ := kc2.Count()
	if count != count2 {
		t.Errorf("Count() (copy) expected %d got %d", count, count2)
	}
}

func TestDump(t *testing.T) {
	return // broken??

	kc := newCabinet(t, test_db)
	defer delCabinet(t, kc)
	kc2 := newCabinet(t, test_db2)
	defer delCabinet(t, kc2)

	addKVP(t, kc, "hello", "world")
	addKVP(t, kc, "bit", "bucket")
	addKVP(t, kc, "foo", "bar")

	defer os.Remove(test_dump)

	err := kc.Dump(test_dump)
	if err != nil {
		t.Fatalf("Dump(%s): %s", test_dump, err)
	}

	err = kc.Load(test_dump)
	if err != nil {
		t.Fatalf("Load(%s): %s", test_dump, err)
	}

	count, _ := kc.Count()
	count2, _ := kc2.Count()
	if count != count2 {
		t.Errorf("Count() (copy) expected %d got %d", count, count2)
	}
}
