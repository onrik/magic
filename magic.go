package magic

import (
	"fmt"
	"reflect"
)

// Converter is a custom converter for differect types
type Converter func(from, to reflect.Value) (bool, error)

// isPtrOf checks that v1.Type is ptr of v2.Type
func isPtrOf(v1, v2 reflect.Value) bool {
	if v1.Type().Kind() != reflect.Ptr {
		return false
	}

	return v1.Type().Elem().Kind() == v2.Type().Kind()
}

func convert(from, to reflect.Value, opts *options) error {
	// fmt.Println("    convert", from.Type(), "(", from, ")", "->", to.Type(), "(", to, ")")
	// Same types
	if from.Type() == to.Type() {
		to.Set(from)
		return nil
	}

	// Convertible
	if from.Type().Kind() == to.Type().Kind() && from.Type().ConvertibleTo(to.Type()) {
		to.Set(from.Convert(to.Type()))
		return nil
	}

	// Ptr to Ptr
	if from.Type().Kind() == reflect.Ptr && to.Type().Kind() == reflect.Ptr {
		if from.Type().Elem().Kind() == to.Type().Elem().Kind() {
			if from.IsNil() {
				return nil
			}
			if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}

			return convert(from.Elem(), to.Elem(), opts)
		}
	}

	// Type to Ptr
	if isPtrOf(to, from) {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		return convert(from, to.Elem(), opts)
	}

	// Ptr to Type
	if isPtrOf(from, to) {
		if from.IsNil() {
			return nil
		}
		return convert(from.Elem(), to, opts)
	}

	// Different structs
	if from.Type().Kind() == reflect.Struct && to.Type().Kind() == reflect.Struct {
		return convertStruct(from, to, opts)
	}

	// Slices
	if from.Type().Kind() == reflect.Slice && to.Type().Kind() == reflect.Slice {
		return convertSlice(from, to, opts)
	}

	// Maps
	if from.Type().Kind() == reflect.Map && to.Type().Kind() == reflect.Map {
		if from.Type().Key() == to.Type().Key() {
			return convertMap(from, to, opts)
		}
	}

	// Use converters for different types
	for _, c := range opts.converters {
		ok, err := c(from, to)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	return fmt.Errorf("cannot convert %v to %v", from.Type(), to.Type())
}

func convertSlice(from, to reflect.Value, opts *options) error {
	for i := 0; i < from.Len(); i++ {
		elem := reflect.New(to.Type().Elem()).Elem()
		// fmt.Println("    convert slice", from.Index(i).Type(), elem.Type())
		err := convert(from.Index(i), elem, opts)
		if err != nil {
			return err
		}

		to.Set(reflect.Append(to, elem))
	}

	return nil
}

func convertMap(from, to reflect.Value, opts *options) error {
	if to.IsNil() {
		to.Set(reflect.MakeMap(to.Type()))
	}

	for _, k := range from.MapKeys() {
		elem := reflect.New(to.Type().Elem()).Elem()
		if err := convert(from.MapIndex(k), elem, opts); err != nil {
			return err
		}

		to.SetMapIndex(k, elem)
	}

	return nil
}

func convertStruct(from, to reflect.Value, opts *options) error {
	for i := 0; i < from.Type().NumField(); i++ {
		name := from.Type().Field(i).Name
		mappedName := opts.mapping[name]
		if mappedName == "" {
			mappedName = name
		}

		_, ok := to.Type().FieldByName(mappedName)
		if !ok {
			continue
		}

		v1 := from.FieldByName(name)
		v2 := to.FieldByName(mappedName)
		err := convert(v1, v2, opts)
		if err != nil {
			return fmt.Errorf("%s: %s", name, err.Error())
		}
	}
	return nil
}

// Map maps struct or slices values
func Map(from, to interface{}, opts ...func(*options)) error {
	valueFrom := reflect.ValueOf(from)
	valueTo := reflect.ValueOf(to)

	typeFrom := reflect.TypeOf(from)
	typeTo := reflect.TypeOf(to)
	// fmt.Println("Map", typeFrom, "(", valueFrom, ")", "->", typeTo, "(", valueTo, ")")

	if typeFrom.Kind() == reflect.Ptr {
		valueFrom = valueFrom.Elem()
		typeFrom = typeFrom.Elem()
	}

	if typeTo.Kind() == reflect.Ptr {
		valueTo = valueTo.Elem()
		typeTo = typeTo.Elem()
	}

	if !valueTo.CanAddr() {
		return fmt.Errorf("%T is not addressable", to)
	}

	o := options{}
	for _, option := range opts {
		option(&o)
	}

	if typeFrom.Kind() == reflect.Slice && typeTo.Kind() == reflect.Slice {
		return convertSlice(valueFrom, valueTo, &o)
	}
	if typeFrom.Kind() == reflect.Struct && typeTo.Kind() == reflect.Struct {
		return convertStruct(valueFrom, valueTo, &o)
	}

	return fmt.Errorf("Cannot map %T to %T", from, to)
}
