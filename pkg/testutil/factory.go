package testutil

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrWithManyNoValues                   = errors.New("no values passed to WithMany")
	ErrWithManyHasDifferentType           = errors.New("values passed to WithMany have different types")
	ErrInsertAssWithoutAss                = errors.New("inserting associations without any associations")
	ErrInsertWithAssOnlyWithBuild         = errors.New("InsertWithAss can only be used with Build")
	ErrInsertListWithAssOnlyWithBuildList = errors.New("InsertListWithAss can only be used with BuildList")
	ErrDestValueNotStruct                 = errors.New("destination value is not a struct")
	ErrSourceValueNotStruct               = errors.New("source value is not a struct")
	ErrDestAndSourceIsDiff                = errors.New("destination and source type is different")
)

// bluePrintFunc is a client-defined function to create a new value
type bluePrintFunc[T any] func(i int, last T) T

// inserter is a client-defined function to insert a value into the database
type inserter[T any] func(db *sql.DB, v T) (T, error)

type Factory[T any] struct {
	db        *sql.DB
	dataType  reflect.Type
	index     int
	last      *T
	empty     T
	list      []*T
	bluePrint bluePrintFunc[T]
	inserter  inserter[T]
	errors    []error

	// map from name to list of associations
	// e.g. "User" -> []*User
	associations map[string][]interface{}
	// map from tag to field name
	// e.g. "User" -> "UserID"
	tagToField map[string]string
}

// NewFactory creates a new factory instance
func NewFactory[T any](db *sql.DB, v T, bluePrint bluePrintFunc[T], inserter inserter[T]) *Factory[T] {
	dataType := reflect.TypeOf(v)

	tagToField := map[string]string{}
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		tag := field.Tag.Get("factory")
		if tag == "" {
			continue
		}

		tagToField[tag] = field.Name
	}

	last := reflect.New(dataType).Elem().Interface().(T)

	return &Factory[T]{
		db:           db,
		dataType:     dataType,
		index:        0,
		last:         &last,
		empty:        last,
		bluePrint:    bluePrint,
		inserter:     inserter,
		associations: map[string][]interface{}{},
		tagToField:   tagToField,
	}
}

// Build creates a new value using the bluePrint function
func (f *Factory[T]) Build() *Factory[T] {
	v := f.bluePrint(f.index+1, *f.last)

	f.index++
	f.last = &v

	return f
}

// BuildList creates a list of n values using the bluePrint function
func (f *Factory[T]) BuildList(n int) *Factory[T] {
	f.list = make([]*T, n)
	for i := 0; i < n; i++ {
		v := f.bluePrint(i+1, *f.last)
		f.list[i] = &v
		f.index++
		f.last = f.list[i]
	}

	return f
}

// Overwrite overwrites the last value with the given value
func (f *Factory[T]) Overwrite(ow T) *Factory[T] {
	if err := insertValues(f.last, ow); err != nil {
		f.errors = append(f.errors, err)
	}

	return f
}

// Overwrites overwrites the last n values with the given values
func (f *Factory[T]) Overwrites(ows []T) *Factory[T] {
	len := len(ows)
	if len == 0 || f.list == nil {
		return f
	}

	var ow T
	for i, v := range f.list {
		if i < len {
			ow = ows[i]
		}

		if err := insertValues(v, ow); err != nil {
			f.errors = append(f.errors, err)
		}
		f.list[i] = v
	}

	return f
}

// Insert inserts the value using the inserter function
func (f *Factory[T]) Insert() (T, error) {
	return f.inserter(f.db, *f.last)
}

// InsertList inserts the list of values using the inserter function
func (f *Factory[T]) InsertList() ([]T, error) {
	if f.list == nil {
		return nil, nil
	}

	list := make([]T, len(f.list))
	for i, v := range f.list {
		v, err := f.inserter(f.db, *v)
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
	last := reflect.New(f.dataType).Elem().Interface().(T)
	f.last = &last
	f.list = nil
	f.index = 0
}

// Value returns the last value
func (f *Factory[T]) Value() T {
	return *f.last
}

// WihtOne set one association to the factory value
func (f *Factory[T]) WithOne(value interface{}) *Factory[T] {
	typeOfValue := reflect.TypeOf(value)
	valueOfValue := reflect.ValueOf(value)

	if valueOfValue.Kind() != reflect.Ptr {
		f.errors = append(f.errors, fmt.Errorf("type: %v, value: %v passed to WithOne is not a pointer", typeOfValue.Name(), value))
		return f
	}

	name := typeOfValue.Elem().Name()
	v := reflect.New(typeOfValue.Elem()).Interface()
	if valueOfValue.Elem().Kind() != reflect.Struct {
		f.errors = append(f.errors, fmt.Errorf("type %v, value: %v passed to WithOne is not a struct", name, v))
		return f
	}

	if _, ok := f.tagToField[name]; !ok {
		f.errors = append(f.errors, fmt.Errorf("type %v, value: %v passed to WithOne is not found at tag", name, v))
		return f
	}

	setNonZeroValues(value, f.index)
	f.associations[name] = []interface{}{value}
	f.index++
	return f
}

// WithMany sets multiple associations to the factory value
func (f *Factory[T]) WithMany(i int, values ...interface{}) *Factory[T] {
	lenValues := len(values)
	if lenValues == 0 {
		f.errors = append(f.errors, ErrWithManyNoValues)
		return f
	}

	firstVal := values[0]
	if reflect.ValueOf(firstVal).Kind() != reflect.Ptr {
		f.errors = append(f.errors, fmt.Errorf("type: %v, value: %v passed to WithMany is not a pointer", reflect.TypeOf(firstVal).Name(), firstVal))
		return f
	}

	name := reflect.TypeOf(firstVal).Elem().Name()
	emptyVal := reflect.New(reflect.TypeOf(firstVal).Elem()).Interface()
	vs := make([]interface{}, 0, i)
	listLen := 1
	if f.list != nil {
		listLen = len(f.list)
	}
	for k := 0; k < i && k < listLen; k++ {
		curVal := emptyVal
		if k < lenValues {
			curVal = values[k]
		}

		typeOfCurVal := reflect.TypeOf(curVal)
		valueOfCurVal := reflect.ValueOf(curVal)

		curName := typeOfCurVal.Elem().Name()
		if curName != name {
			f.errors = append(f.errors, fmt.Errorf("%s: one is %s, the other is %s", ErrWithManyHasDifferentType, name, curName))
			return f
		}

		if valueOfCurVal.Kind() != reflect.Ptr {
			f.errors = append(f.errors, fmt.Errorf("type: %v, value: %v passed to WithMany is not a pointer", typeOfCurVal.Name(), curVal))
			return f
		}

		if valueOfCurVal.Elem().Kind() != reflect.Struct {
			f.errors = append(f.errors, fmt.Errorf("type %v, value: %v passed to WithMany is not a struct", typeOfCurVal.Elem().Name(), valueOfCurVal.Elem().Interface()))
			return f
		}

		if _, ok := f.tagToField[curName]; !ok {
			f.errors = append(f.errors, fmt.Errorf("type %v, value: %v passed to WithMany is not found at tag", typeOfCurVal.Elem().Name(), valueOfCurVal.Elem().Interface()))
			return f
		}

		setNonZeroValues(curVal, f.index)
		vs = append(vs, curVal)
		f.index++
	}

	f.associations[name] = vs
	return f
}

// InsertWithAss inserts the value with associations
func (f *Factory[T]) InsertWithAss() (T, []interface{}, error) {
	if f.list != nil {
		return f.empty, nil, ErrInsertWithAssOnlyWithBuild
	}

	ass, err := f.genAndInsertAss()
	if err != nil {
		return f.empty, nil, err
	}

	if err := f.setAss(); err != nil {
		return f.empty, nil, err
	}

	val, err := f.Insert()
	if err != nil {
		return f.empty, nil, err
	}

	return val, ass, nil
}

// InsertListWithAss inserts the list of value with associations
func (f *Factory[T]) InsertListWithAss() ([]T, []interface{}, error) {
	if f.list == nil {
		return nil, nil, ErrInsertListWithAssOnlyWithBuildList
	}

	ass, err := f.genAndInsertAss()
	if err != nil {
		return nil, nil, err
	}

	if err := f.setAssWithList(); err != nil {
		return nil, nil, err
	}

	vals, err := f.InsertList()
	if err != nil {
		return nil, nil, err
	}

	return vals, ass, nil
}

// genAndInsertAss insert associations value into db,
// and returns value with incremental id
func (f *Factory[T]) genAndInsertAss() ([]interface{}, error) {
	if f.associations == nil {
		return nil, ErrInsertAssWithoutAss
	}

	if f.errors != nil {
		return nil, genFinalError(f.errors)
	}

	result := []interface{}{}
	for name, vs := range f.associations {
		for _, v := range vs {
			if err := f.insertAss(name, v); err != nil {
				return nil, err
			}

			// deference the value from pointer v
			result = append(result, reflect.ValueOf(v).Elem().Interface())
		}
	}

	return result, nil
}

// setAss sets the assoication connection to one factory value.
// setAss can only be used with Build
func (f *Factory[T]) setAss() error {
	for name, vs := range f.associations {
		// use vs[0] because we can make sure setAss(InsertWithAss) only invoke with Build function
		// which means there's only one factory value
		// so that each associations only allow one value
		if err := f.setField(f.last, f.tagToField[name], vs[0]); err != nil {
			return err
		}
	}

	return nil
}

// setAssWithList sets the assoication connection to list of factory values
// setAssWithList can only be used with BuildList
func (f *Factory[T]) setAssWithList() error {
	cachePrev := map[string]interface{}{}

	for i, l := range f.list {
		for name, vs := range f.associations {
			var v interface{}
			if i >= len(vs) {
				v = cachePrev[name]
			} else {
				v = vs[i]
				cachePrev[name] = vs[i]
			}

			if err := f.setField(l, f.tagToField[name], v); err != nil {
				return err
			}
		}
	}

	return nil
}

// setField sets the association value to the target field with the given name and source.
// It set the name field of the target with the ID field of the source
func (f *Factory[T]) setField(target *T, name string, source interface{}) error {
	targetField := reflect.ValueOf(target).Elem().FieldByName(name)
	if !targetField.IsValid() {
		return fmt.Errorf("target field not found: %s at setField", name)
	}

	if !targetField.CanSet() {
		return fmt.Errorf("target field cannot be set: %s at setField", name)
	}

	sourceIDField := reflect.ValueOf(source).Elem().FieldByName("ID")
	if !sourceIDField.IsValid() {
		return fmt.Errorf("source field not found: ID at setField")
	}

	if sourceIDField.Kind() != reflect.Int ||
		sourceIDField.Kind() != reflect.Int64 ||
		sourceIDField.Kind() != reflect.Int32 ||
		sourceIDField.Kind() != reflect.Int16 ||
		sourceIDField.Kind() != reflect.Int8 {
		return fmt.Errorf("ID field is not an integer at setField")
	}

	targetField.SetInt(sourceIDField.Int())

	return nil
}

// insertAss inserts the association value into the database
func (f *Factory[T]) insertAss(name string, v interface{}) error {
	tableName := strings.ToLower(name) + "s"
	fieldNames := []string{}
	fieldValues := []interface{}{}

	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		n := val.Type().Field(i).Name
		if n == "ID" {
			continue
		}

		fieldNames = append(fieldNames, n)
		fieldValues = append(fieldValues, val.Field(i).Interface())
	}

	stmt := `INSERT INTO ` + tableName + ` (`
	for i, v := range fieldNames {
		stmt += v
		if i < len(fieldNames)-1 {
			stmt += ", "
		}
	}
	stmt += `) VALUES (`
	for i := range fieldValues {
		stmt += "?"
		if i < len(fieldValues)-1 {
			stmt += ", "
		}
	}
	stmt += `)`

	res, err := f.db.Exec(stmt, fieldValues...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	val.FieldByName("ID").SetInt(id)
	return nil
}

// setNonZeroValues sets non-zero values to the given struct
func setNonZeroValues(v interface{}, i int) {
	val := reflect.ValueOf(v).Elem()
	typeOfVal := val.Type()

	for k := 0; k < val.NumField(); k++ {
		field := val.Field(k)
		zeroValue := reflect.Zero(typeOfVal.Field(k).Type)

		if reflect.DeepEqual(field.Interface(), zeroValue.Interface()) &&
			typeOfVal.Field(k).Name != "ID" {
			v := genNonZeroValue(typeOfVal.Field(k).Type, i)
			field.Set(reflect.ValueOf(v))
		}
	}
}

// genNonZeroValue generates a non-zero value for the given type
func genNonZeroValue(fieldType reflect.Type, i int) interface{} {
	switch fieldType.Kind() {
	case reflect.Int:
		return i
	case reflect.Int8:
		return int8(i)
	case reflect.Int16:
		return int16(i)
	case reflect.Int32:
		return int32(i)
	case reflect.Int64:
		return int64(i)
	case reflect.Uint:
		return uint(i)
	case reflect.Uint8:
		return uint8(i)
	case reflect.Uint16:
		return uint16(i)
	case reflect.Uint32:
		return uint32(i)
	case reflect.Uint64:
		return uint64(i)
	case reflect.Float32:
		return float32(i)
	case reflect.Float64:
		return float64(i)
	case reflect.Bool:
		return true
	case reflect.String:
		return "test" + fmt.Sprint(i)
	case reflect.Struct:
		return reflect.New(fieldType).Elem().Interface()
	case reflect.Ptr:
		return reflect.New(fieldType.Elem()).Elem().Interface()
	case reflect.Slice:
		return reflect.MakeSlice(fieldType, 0, 0).Interface()
	default:
		return reflect.New(fieldType).Elem().Interface()
	}
}

// insertValues inserts non-zero values from src to dest
func insertValues[T any](dest *T, src T) error {
	destValue := reflect.ValueOf(dest).Elem()
	srcValue := reflect.ValueOf(src)

	if destValue.Kind() != reflect.Struct {
		return ErrDestValueNotStruct
	}

	if srcValue.Kind() != reflect.Struct {
		return ErrSourceValueNotStruct
	}

	if destValue.Type() != srcValue.Type() {
		return ErrDestAndSourceIsDiff
	}

	for i := 0; i < destValue.NumField(); i++ {
		destField := destValue.Field(i)
		srcField := srcValue.FieldByName(destValue.Type().Field(i).Name)

		if srcField.IsValid() && destField.Type() == srcField.Type() && !srcField.IsZero() {
			destField.Set(srcField)
		}
	}

	return nil
}

// genFinalError generates a final error message from the given errors
func genFinalError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	errorMessages := make([]string, len(errs))
	for i, err := range errs {
		errorMessages[i] = err.Error()
	}

	return fmt.Errorf("encountered the following errors:\n%s", strings.Join(errorMessages, "\n"))
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
