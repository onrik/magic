package magic

import (
	"reflect"
	"testing"
	"time"
)

func val(i interface{}) reflect.Value {
	return reflect.ValueOf(i)
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

type testType7 struct {
	T *testType1
}

type testType8 struct {
	T *testType2
}

func TestMapStruct(t *testing.T) {
	t1 := testType1{34, "John", "111", []string{"1"}}
	t2 := testType2{}

	err := Map(t1, &t2)
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

func TestInvalidType(t *testing.T) {
	s1 := struct {
		ID int
	}{}
	s2 := struct {
		ID string
	}{}
	err := Map(s1, &s2)
	if err == nil || err.Error() != "ID: Cannot convert int to string" {
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
	if err == nil || err.Error() != "Tags: Cannot convert int to string" {
		t.Fatal(err)
	}
}

func TestPtrToType(t *testing.T) {
	i := 4385
	s1 := struct {
		ID *int
	}{&i}
	s2 := struct {
		ID int
	}{}

	err := Map(s1, &s2)
	assert(t, err, nil)
	assert(t, *s1.ID, s2.ID)

	s1.ID = nil
	err = Map(s1, &s2)
	assert(t, err, nil)
	assert(t, s2.ID, i)
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

func TestConverter(t *testing.T) {
	now := time.Now()
	t5 := testType5{45, now}
	t6 := testType6{}

	err := Map(t5, &t6, WithConverters(timeToUnix))
	assert(t, err, nil)
	assert(t, t5.ID, t6.ID)
	assert(t, t6.Created, now.Unix())
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

func assert(t *testing.T, i1, i2 interface{}) {
	if i1 != i2 {
		t.Fatalf("%v != %v", i1, i2)
	}
}
