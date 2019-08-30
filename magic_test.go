package magic

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func val(i interface{}) reflect.Value {
	return reflect.ValueOf(i)
}

func assert(t *testing.T, i1, i2 interface{}) {
	t.Helper()
	if i1 != i2 {
		t.Fatalf("%v != %v", i1, i2)
	}
}

func assertTrue(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatal(false)
	}
}

type testType1 struct {
	ID       int
	Name     string
	Password string
	Tags     []string
}

type testType2 struct {
	ID      int
	Name    string
	Created time.Time
	Tags    []string
}

type testType3 struct {
	ID   *int
	Name *string
}

type testType4 struct {
	ID   int
	Tags []*string
}

type testType5 struct {
	ID      int
	Created time.Time
}

type testType6 struct {
	ID      int
	Created int64
}

func timeToUnix(from, to reflect.Value) (bool, error) {
	if to.Type() != reflect.TypeOf(int64(0)) {
		return false, nil
	}
	t, ok := from.Interface().(time.Time)
	if ok {
		to.SetInt(t.Unix())
		return true, nil
	}

	return false, nil
}

func TestMapStruct(t *testing.T) {
	t1 := testType1{34, "John", "111", []string{"1"}}
	t2 := testType2{}

	err := Map(t1, &t2)
	assert(t, err, nil)
	assert(t, t2.ID, t1.ID)
	assert(t, t2.Name, t1.Name)

	t2 = testType2{}
	err = Map(&t1, &t2)
	assert(t, err, nil)
	assert(t, t2.ID, t1.ID)
	assert(t, t2.Name, t1.Name)
}

func TestMapStructWithPointers(t *testing.T) {
	t1 := testType1{34, "John", "111", []string{"1"}}
	t3 := testType3{}

	err := Map(t1, &t3)
	assert(t, err, nil)
	assert(t, *t3.ID, t1.ID)
	assert(t, *t3.Name, t1.Name)

	id := 44
	name := "test"
	t3 = testType3{&id, &name}
	t1 = testType1{}

	err = Map(t1, &t3)
	assert(t, err, nil)
	assert(t, *t3.ID, t1.ID)
	assert(t, *t3.Name, t1.Name)
}

func TestMapStructWithPointersSlice(t *testing.T) {
	t1 := testType1{34, "John", "111", []string{"1"}}

	t4 := testType4{}

	err := Map(t1, &t4)
	assert(t, err, nil)
	assert(t, t4.ID, t1.ID)
	assert(t, *t4.Tags[0], t1.Tags[0])
}

func TestMapSlice(t *testing.T) {
	t1 := []testType1{{34, "John", "111", []string{"1"}}}

	t2 := []testType2{}

	err := Map(t1, &t2)
	assert(t, err, nil)
	assert(t, len(t2), 1)
	assert(t, t2[0].ID, t1[0].ID)
	assert(t, t2[0].Name, t1[0].Name)
	assert(t, t2[0].Tags[0], t1[0].Tags[0])
}

func TestMapPointersSlice(t *testing.T) {
	t1 := []testType1{{34, "John", "111", []string{"1"}}}
	t2 := []*testType2{}

	err := Map(t1, &t2)
	assert(t, err, nil)
	assert(t, len(t2), 1)
	assert(t, t2[0].ID, t1[0].ID)
	assert(t, t2[0].Name, t1[0].Name)
	assert(t, t2[0].Tags[0], t1[0].Tags[0])

	t2 = []*testType2{
		{43, "John", time.Now(), []string{"2"}},
	}
	t1 = []testType1{}
	err = Map(t2, &t1)
	assert(t, err, nil)
	assert(t, len(t1), 1)
	assert(t, t1[0].ID, t2[0].ID)
	assert(t, t1[0].Name, t2[0].Name)
	assert(t, len(t1[0].Tags), 1)
	assert(t, t1[0].Tags[0], t2[0].Tags[0])
}

func TestInvalidType(t *testing.T) {
	s1 := struct {
		ID int
	}{}
	s2 := struct {
		ID string
	}{}
	err := Map(s1, &s2)
	if err == nil || err.Error() != "ID: cannot convert int to string" {
		t.Fatal(err)
	}
}

func TestInvalidSlice(t *testing.T) {
	s1 := struct {
		Tags []int
	}{[]int{1}}
	s2 := struct {
		Tags []string
	}{}
	err := Map(s1, &s2)
	if err == nil || err.Error() != "Tags: cannot convert int to string" {
		t.Fatal(err)
	}
}

func TestPtrToType(t *testing.T) {
	i := 4385
	type Foo1 struct {
		Bar int
	}
	type Foo2 struct {
		Bar int
		S   float64
	}

	f := Foo1{56}
	s1 := struct {
		ID  *int
		Foo *Foo1
	}{&i, &f}
	s2 := struct {
		ID  int
		Foo Foo2
	}{}

	err := Map(s1, &s2)
	assert(t, err, nil)
	assert(t, *s1.ID, s2.ID)
	assert(t, s1.Foo.Bar, s2.Foo.Bar)

	s1.ID = nil
	s1.Foo = nil
	err = Map(s1, &s2)
	assert(t, err, nil)
	assert(t, s2.ID, i)
	if s1.Foo != nil {
		t.Fatal(s1.Foo)
	}
}

func TestTypeToPtr(t *testing.T) {
	s1 := struct {
		ID int
	}{23984}
	s2 := struct {
		ID *int
	}{}

	err := Map(s1, &s2)
	assert(t, err, nil)
	assert(t, s1.ID, *s2.ID)
}

func TestConvertable(t *testing.T) {
	type Maps map[string]string
	type S1 struct {
		V map[string]string
	}
	type S2 struct {
		V Maps
	}

	s1 := S1{map[string]string{
		"foo": "bar",
	}}
	s2 := S2{}

	err := Map(s1, &s2)
	assert(t, err, nil)
	assert(t, len(s1.V), len(s2.V))
	assert(t, s1.V["foo"], s2.V["foo"])
}

func TestMapError(t *testing.T) {
	type S1 struct {
		ID int
	}
	s1 := S1{45}
	s2 := []string{}

	err := Map(s1, s2)
	if err == nil || err.Error() != "[]string is not addressable" {
		t.Fatal(err)
	}

	err = Map(s1, &s2)
	if err == nil || err.Error() != "Cannot map magic.S1 to *[]string" {
		t.Fatal(err)
	}
}

func TestConvertMap(t *testing.T) {
	type G1 struct {
		ID   int
		Name string
	}

	type G2 struct {
		ID int
	}

	type U1 struct {
		Groups map[string]G1
	}

	type U2 struct {
		Groups map[string]G2
	}

	u1 := U1{
		Groups: map[string]G1{
			"test": G1{43, "test"},
		},
	}
	u2 := U2{}

	err := Map(u1, &u2)
	assert(t, err, nil)
	assertTrue(t, len(u2.Groups) == 1)
	assertTrue(t, u2.Groups["test"].ID == 43)
}

func TestConvertMapToPtr(t *testing.T) {
	type G1 struct {
		ID   int
		Name string
	}

	type G2 struct {
		ID int
	}

	type U1 struct {
		Groups map[string]G1
	}

	type U2 struct {
		Groups *map[string]G2
	}

	u1 := U1{
		Groups: map[string]G1{
			"test": G1{43, "test"},
		},
	}
	u2 := U2{}

	err := Map(u1, &u2)
	assert(t, err, nil)
	assertTrue(t, u2.Groups != nil)
	assertTrue(t, len(*u2.Groups) == 1)
	assertTrue(t, (*(u2.Groups))["test"].ID == 43)
}

func TestConvertPtrToMap(t *testing.T) {
	type G1 struct {
		ID   int
		Name string
	}

	type G2 struct {
		ID int
	}

	type U1 struct {
		Groups *map[string]G1
	}

	type U2 struct {
		Groups map[string]G2
	}

	u1 := U1{
		Groups: &map[string]G1{
			"test": G1{43, "test"},
		},
	}
	u2 := U2{}

	err := Map(u1, &u2)
	assert(t, err, nil)
	assertTrue(t, len(u2.Groups) == 1)
	assertTrue(t, u2.Groups["test"].ID == 43)
}

func TestConvertInvalidMap(t *testing.T) {
	type U1 struct {
		Groups map[string]string
	}

	type U2 struct {
		Groups map[string]int
	}

	u1 := U1{
		Groups: map[string]string{
			"foo": "bar",
		},
	}
	u2 := U2{}

	err := Map(u1, &u2)
	assertTrue(t, err != nil)
	assert(t, err.Error(), "Groups: cannot convert string to int")
}

func TestConvertPtrToPtr(t *testing.T) {
	type G1 struct {
		ID int
	}

	type G2 struct {
		ID   int
		Name string
	}

	type U1 struct {
		ID    *int
		Group *G1
		G     *G1
	}

	type U2 struct {
		ID    *int
		Name  string
		Group *G2
		G     *G2
	}

	u1 := U1{Group: &G1{4}}
	u2 := U2{}

	err := Map(u1, &u2)
	assertTrue(t, err == nil)
	assertTrue(t, u2.Group != nil)
	assertTrue(t, u2.Group.ID == 4)
}

func TestConverter(t *testing.T) {
	now := time.Now()
	t5 := testType5{45, now}
	t6 := testType6{}

	err := Map(t5, &t6, WithConverters(timeToUnix))
	assert(t, err, nil)
	assert(t, t5.ID, t6.ID)
	assert(t, t6.Created, now.Unix())

	e := fmt.Errorf("E")
	err = Map(t5, &t6, WithConverters(func(v1, v2 reflect.Value) (bool, error) {
		return false, e
	}))
	assert(t, err.Error(), "Created: E")
}

func TestMapping(t *testing.T) {
	t1 := struct {
		ID string
	}{"111"}
	t2 := struct {
		UUID string
	}{}

	err := Map(t1, &t2, WithMapping(map[string]string{
		"ID": "UUID",
	}))
	assert(t, err, nil)
	assert(t, t1.ID, t2.UUID)
}
