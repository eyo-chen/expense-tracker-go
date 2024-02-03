package testutil

import (
	"database/sql"
	"reflect"
)

// bluePrintFunc is a client-defined function to create a new value
type bluePrintFunc[T any] func(i int, last T) T

// inserter is a client-defined function to insert a value into the database
type inserter[T any] func(db *sql.DB, v T) (T, error)

type Factory[T any] struct {
	db        *sql.DB
	dataType  reflect.Type
	index     int
	last      T
	list      []T
	bluePrint bluePrintFunc[T]
	inserter  inserter[T]
}

// NewFactory creates a new factory instance
func NewFactory[T any](db *sql.DB, v T, bluePrint bluePrintFunc[T], inserter inserter[T]) *Factory[T] {
	return &Factory[T]{
		db:        db,
		dataType:  reflect.PtrTo(reflect.TypeOf(v)),
		index:     0,
		last:      v,
		bluePrint: bluePrint,
		inserter:  inserter,
	}
}

// Overwrite overwrites the last value with the given value
func (f *Factory[T]) Overwrite(ow *T) *Factory[T] {
	insertValues(&f.last, *ow)
	return f
}

// Overwrites overwrites the last n values with the given values
func (f *Factory[T]) Overwrites(ows []*T) *Factory[T] {
	len := len(ows)
	if len == 0 || f.list == nil {
		return f
	}

	var ow *T
	for i, v := range f.list {
		if i < len {
			ow = ows[i]
		}

		insertValues(&v, *ow)
		f.list[i] = v
	}

	return f
}

// Build creates a new value using the bluePrint function
func (f *Factory[T]) Build() *Factory[T] {
	v := f.bluePrint(f.index+1, f.last)

	f.index++
	f.last = v

	return f
}

// Insert inserts the last value using the inserter function
func (f *Factory[T]) Insert() (T, error) {
	return f.inserter(f.db, f.last)
}

// BuildList creates a list of n values using the bluePrint function
func (f *Factory[T]) BuildList(n int) *Factory[T] {
	f.list = make([]T, n)
	for i := 0; i < n; i++ {
		f.list[i] = f.bluePrint(i+1, f.last)
		f.index++
		f.last = f.list[i]
	}

	return f
}

// InsertList inserts the list of values using the inserter function
func (f *Factory[T]) InsertList() ([]T, error) {
	if f.list == nil {
		return nil, nil
	}

	list := make([]T, len(f.list))
	for i, v := range f.list {
		v, err := f.inserter(f.db, v)
		if err != nil {
			return nil, err
		}
		list[i] = v
	}

	f.list = nil
	return list, nil
}

// Reset resets the factory to its initial state
func (f *Factory[T]) Reset() {
	f.last = reflect.New(f.dataType.Elem()).Elem().Interface().(T)
	f.list = nil
	f.index = 0
}

// Value returns the last value
func (f *Factory[T]) Value() T {
	return f.last
}

// insertValues inserts non-zero values from src to dest
func insertValues[T any](dest *T, src T) {
	destValue := reflect.ValueOf(dest).Elem()
	srcValue := reflect.ValueOf(src)

	if destValue.Kind() != reflect.Struct || srcValue.Kind() != reflect.Struct {
		return
	}

	if destValue.Type() != srcValue.Type() {
		return
	}

	for i := 0; i < destValue.NumField(); i++ {
		destField := destValue.Field(i)
		srcField := srcValue.FieldByName(destValue.Type().Field(i).Name)

		if srcField.IsValid() && destField.Type() == srcField.Type() && !srcField.IsZero() {
			destField.Set(srcField)
		}
	}
}

// type Builder[T any] struct {
// 	f         *Factory[T]
// 	index     int
// 	last      T
// 	overwrite *T
// 	bluePrint bluePrintFunc[T]
// 	inserter  inserter[T]
// }

// func (f *Factory[T]) NewBuilder(bluePrint bluePrintFunc[T], inserter inserter[T], last T, overwrite *T) *Builder[T] {
// 	return &Builder[T]{
// 		f:         f,
// 		index:     0,
// 		last:      last,
// 		overwrite: overwrite,
// 		bluePrint: bluePrint,
// 		inserter:  inserter,
// 	}
// }
