package efactory

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory/db"
)

// Config is the configuration for the factory
type Config[T any] struct {
	// BluePrint is a client-defined function to create a new value
	// if not provided, non-zero default values is set
	// BluePrint must follow the below signature
	// type bluePrintFunc[T any] func(i int, last T) T
	BluePrint bluePrintFunc[T]

	// DB is the interface for the Database
	// must be provided if want to do database operations
	// use esql if using raw sql
	// use egorm if using gorm
	// use emongo if using mongo
	DB db.Database

	// StorageName is the name specified where the value is stored
	// it will be table name for sql, and collection name for mongodb
	// it is optional
	// if not provided, camel case of the type name will be used
	StorageName string

	// OmitID is to check if the ID field should be omitted
	// it is optional
	// if not provided, it will be true
	OmitID bool

	// isSetZeroValue is to determine if the zero value should be set
	// it is optional
	// if not provided, it will be true
	IsSetZeroValue *bool
}

type Factory[T any] struct {
	db             db.Database
	bluePrint      bluePrintFunc[T]
	storageName    string
	dataType       reflect.Type
	empty          T
	index          int
	omitID         bool
	isSetZeroValue bool

	// map from name to trait function
	traits map[string]setTraiter[T]

	// map from name to list of associations
	// e.g. "User" -> []*User
	associations map[string][]interface{}

	// map from tag to metadata
	// e.g. "User" -> {tableName: "users", fieldName: "UserID"}
	tagToInfo map[string]tagInfo
}

// bluePrintFunc is a client-defined function to create a new value
type bluePrintFunc[T any] func(i int, last T) T

// SetTrait is a client-defined function to add a trait to mutate the value
type setTraiter[T any] func(v *T)

// tagInfo is the metadata for the tag
type tagInfo struct {
	tableName string
	fieldName string
}

// builder is for building a single value
type builder[T any] struct {
	v      *T
	errors []error
	f      *Factory[T]
}

// builderList is for building a list of values
type builderList[T any] struct {
	list   []*T
	errors []error
	f      *Factory[T]
}

func New[T any](v T) *Factory[T] {
	dataType := reflect.TypeOf(v)

	return &Factory[T]{
		dataType:       dataType,
		empty:          reflect.New(dataType).Elem().Interface().(T),
		associations:   map[string][]interface{}{},
		tagToInfo:      map[string]tagInfo{},
		index:          1,
		isSetZeroValue: true,
	}
}

// SetConfig sets the configuration for the factory
func (f *Factory[T]) SetConfig(c Config[T]) *Factory[T] {
	f.bluePrint = c.BluePrint
	f.db = c.DB
	f.omitID = c.OmitID

	if c.StorageName == "" {
		f.storageName = fmt.Sprintf("%ss", CamelToSnake(f.dataType.Name()))
	} else {
		f.storageName = c.StorageName
	}

	if c.IsSetZeroValue != nil {
		f.isSetZeroValue = *c.IsSetZeroValue
	}

	return f
}

// SetTrait adds a trait to the factory
func (f *Factory[T]) SetTrait(name string, tr setTraiter[T]) *Factory[T] {
	if f.traits == nil {
		f.traits = map[string]setTraiter[T]{}
	}

	f.traits[name] = tr
	return f
}

// Reset resets the factory to its initial state
func (f *Factory[T]) Reset() {
	f.index = 1
}

// Build builds a value
func (f *Factory[T]) Build() *builder[T] {
	var v T
	if f.bluePrint != nil {
		v = f.bluePrint(f.index, v)
	}

	if f.isSetZeroValue {
		setNonZeroValues(f.index, &v)
	}

	if !f.omitID {
		if err := f.db.SetIDField(&v, f.index); err != nil {
			return &builder[T]{errors: []error{err}, f: f}
		}
	}

	f.index++

	return &builder[T]{
		v:      &v,
		errors: []error{},
		f:      f,
	}
}

// BuildList creates a list of n values
func (f *Factory[T]) BuildList(n int) *builderList[T] {
	list := make([]*T, n)
	errs := []error{}
	if n < 1 {
		errs = append(errs, errors.New("BuildList: n must be greater than 0"))
		return &builderList[T]{errors: errs}
	}

	for i := 0; i < n; i++ {
		var v T
		if f.bluePrint != nil {
			v = f.bluePrint(f.index, v)
		}

		if f.isSetZeroValue {
			setNonZeroValues(f.index, &v)
		}

		if !f.omitID {
			if err := f.db.SetIDField(&v, f.index); err != nil {
				errs = append(errs, err)
				break
			}
		}

		list[i] = &v
		f.index++
	}

	return &builderList[T]{
		list:   list,
		errors: errs,
		f:      f,
	}
}

// Get returns the value
func (b *builder[T]) Get() (T, error) {
	if len(b.errors) > 0 {
		return b.f.empty, genFinalError(b.errors)
	}

	return *b.v, nil
}

// Get returns the list of values
func (b *builderList[T]) Get() ([]T, error) {
	if len(b.errors) > 0 {
		return nil, genFinalError(b.errors)
	}

	output := make([]T, len(b.list))
	for i, v := range b.list {
		output[i] = *v
	}

	return output, nil
}

// Insert inserts the value into the database
func (b *builder[T]) Insert() (T, error) {
	if b.f.db == nil {
		b.errors = append(b.errors, errors.New("Insert: DB is not provided"))
	}

	if len(b.errors) > 0 {
		return b.f.empty, genFinalError(b.errors)
	}

	val, err := b.f.db.Insert(db.InserParams{StorageName: b.f.storageName, Value: b.v})
	if err != nil {
		b.errors = append(b.errors, err)
		return b.f.empty, genFinalError(b.errors)
	}

	v, ok := val.(*T)
	if !ok {
		b.errors = append(b.errors, fmt.Errorf("Insert: can't convert to pointer"))
		return b.f.empty, genFinalError(b.errors)
	}

	return *v, nil
}

// Insert inserts the list of values into the database
func (b *builderList[T]) Insert() ([]T, error) {
	if b.f.db == nil {
		b.errors = append(b.errors, errors.New("Insert: DB is not provided"))
	}

	if len(b.errors) > 0 {
		return nil, genFinalError(b.errors)
	}

	// convert to any type
	input := make([]interface{}, len(b.list))
	for i, v := range b.list {
		input[i] = v
	}
	vals, err := b.f.db.InsertList(db.InserListParams{StorageName: b.f.storageName, Values: input})
	if err != nil {
		b.errors = append(b.errors, err)
		return nil, genFinalError(b.errors)
	}

	// convert to []T
	output := make([]T, len(vals))
	for i, val := range vals {
		v, ok := val.(*T)
		if !ok {
			b.errors = append(b.errors, fmt.Errorf("Insert: can't convert to pointer"))
			return nil, genFinalError(b.errors)
		}

		output[i] = *v
	}

	return output, nil
}

// Overwrite overwrites the value with the given value
func (b *builder[T]) Overwrite(ow T) *builder[T] {
	if len(b.errors) > 0 {
		return b
	}

	if err := copyValues(b.v, ow); err != nil {
		b.errors = append(b.errors, err)
	}

	return b
}

// Overwrites overwrites the values with the given values
func (b *builderList[T]) Overwrites(ows ...T) *builderList[T] {
	if len(b.errors) > 0 {
		return b
	}

	for i := 0; i < len(ows) && i < len(b.list); i++ {
		if err := copyValues(b.list[i], ows[i]); err != nil {
			b.errors = append(b.errors, err)
			return b
		}
	}

	return b
}

// Overwrite overwrites the values with the given one value
func (b *builderList[T]) Overwrite(ow T) *builderList[T] {
	if len(b.errors) > 0 {
		return b
	}

	for i := 0; i < len(b.list); i++ {
		if err := copyValues(b.list[i], ow); err != nil {
			b.errors = append(b.errors, err)
			return b
		}
	}

	return b
}

// WithTrait invokes the traiter based on given name
func (b *builder[T]) WithTrait(name string) *builder[T] {
	if len(b.errors) > 0 {
		return b
	}

	tr, ok := b.f.traits[name]
	if !ok {
		b.errors = append(b.errors, fmt.Errorf("WithTrait: %s is not defiend at SetTrait", name))
		return b
	}

	tr(b.v)

	return b
}

// WithTraits invokes the traiter based on given names
func (b *builderList[T]) WithTraits(names ...string) *builderList[T] {
	if len(b.errors) > 0 {
		return b
	}

	for i := 0; i < len(names) && i < len(b.list); i++ {
		tr, ok := b.f.traits[names[i]]
		if !ok {
			b.errors = append(b.errors, fmt.Errorf("WithTrait: %s is not defiend at SetTrait", names[i]))
			return b
		}

		tr(b.list[i])
	}

	return b
}

// WihtOne set one association to the factory value
func (b *builder[T]) WithOne(v interface{}) *builder[T] {
	if len(b.errors) > 0 {
		return b
	}

	// set tagToInfo if it's not set
	if len(b.f.tagToInfo) == 0 {
		t, err := genTagToInfo(b.f.dataType)
		if err != nil {
			b.errors = append(b.errors, err)
			return b
		}
		b.f.tagToInfo = t
	}

	if err := setAssValue(v, b.f.tagToInfo, b.f.index, "WithOne"); err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	name := reflect.TypeOf(v).Elem().Name()
	b.f.associations[name] = []interface{}{v}
	b.f.index++
	return b
}

// WihtOne set one association to the factory value
func (b *builderList[T]) WithOne(v interface{}) *builderList[T] {
	if len(b.errors) > 0 {
		return b
	}

	// set tagToInfo if it's not set
	if len(b.f.tagToInfo) == 0 {
		t, err := genTagToInfo(b.f.dataType)
		if err != nil {
			b.errors = append(b.errors, err)
			return b
		}
		b.f.tagToInfo = t
	}

	if err := setAssValue(v, b.f.tagToInfo, b.f.index, "WithOne"); err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	name := reflect.TypeOf(v).Elem().Name()
	b.f.associations[name] = []interface{}{v}
	b.f.index++
	return b
}

// WithMany set many associations to the factory value
func (b *builderList[T]) WithMany(values ...interface{}) *builderList[T] {
	if len(b.errors) > 0 {
		return b
	}

	// set tagToInfo if it's not set
	if len(b.f.tagToInfo) == 0 {
		t, err := genTagToInfo(b.f.dataType)
		if err != nil {
			b.errors = append(b.errors, err)
			return b
		}
		b.f.tagToInfo = t
	}

	var curValName string
	for _, v := range values {
		if err := setAssValue(v, b.f.tagToInfo, b.f.index, "WithMany"); err != nil {
			b.errors = append(b.errors, err)
			return b
		}

		// check if the provided values are of the same type
		// because we have to make sure all the value is pointer (setAssValue does that for us)
		// before we can use Elem()
		if curValName != "" && curValName != reflect.TypeOf(v).Elem().Name() {
			b.errors = append(b.errors, fmt.Errorf("WithMany: provided values are not the same type"))
			return b
		}

		name := reflect.TypeOf(v).Elem().Name()
		b.f.associations[name] = append(b.f.associations[name], v)
		b.f.index++
		curValName = name
	}

	return b
}

// InsertWithAss inserts the value with the associations into the database
func (b *builder[T]) InsertWithAss() (T, []interface{}, error) {
	if len(b.errors) > 0 {
		return b.f.empty, nil, genFinalError(b.errors)
	}

	// generate and insert the associations
	assVals, err := genAndInsertAss(b.f.db, b.f.associations, b.f.tagToInfo)
	if err != nil {
		b.errors = append(b.errors, err)
		return b.f.empty, nil, genFinalError(b.errors)
	}

	// set the connection between the factory value and the associations
	for name, vals := range b.f.associations {
		// use vs[0] because we can make sure InsertWithAss only invoke with Build function
		// which means there's only one factory value
		// so that each associations only allow one value
		fieldName := b.f.tagToInfo[name].fieldName
		if err := setField(b.v, fieldName, vals[0], "InsertWithAss"); err != nil {
			b.errors = append(b.errors, err)
			return b.f.empty, nil, genFinalError(b.errors)
		}
	}

	// insert the factory value
	v, err := b.Insert()
	if err != nil {
		b.errors = append(b.errors, err)
		return b.f.empty, nil, err
	}

	return v, assVals, nil
}

// InsertListWithAss inserts the list of values with the associations into the database
func (b *builderList[T]) InsertWithAss() ([]T, []interface{}, error) {
	if len(b.errors) > 0 {
		return nil, nil, genFinalError(b.errors)
	}

	// generate and insert
	assVals, err := genAndInsertAss(b.f.db, b.f.associations, b.f.tagToInfo)
	if err != nil {
		b.errors = append(b.errors, err)
		return nil, nil, genFinalError(b.errors)
	}

	// set the connection between the factory value and the associations
	// use cachePrev because multiple values can have one association value
	// e.g. multiple transaction belongs to one user
	cachePrev := map[string]interface{}{}
	for i, l := range b.list {
		for name, vs := range b.f.associations {
			var v interface{}
			if i >= len(vs) {
				v = cachePrev[name]
			} else {
				v = vs[i]
				cachePrev[name] = vs[i]
			}

			fieldName := b.f.tagToInfo[name].fieldName
			if err := setField(l, fieldName, v, "InsertWithAss"); err != nil {
				b.errors = append(b.errors, err)
				return nil, nil, genFinalError(b.errors)
			}
		}
	}

	// insert the factory value
	v, err := b.Insert()
	if err != nil {
		b.errors = append(b.errors, err)
		return nil, nil, genFinalError(b.errors)
	}

	return v, assVals, nil
}

// setField sets the value to the name field of the target
func setField(target interface{}, name string, source interface{}, sourceFn string) error {
	targetField := reflect.ValueOf(target).Elem().FieldByName(name)
	if !targetField.IsValid() {
		return fmt.Errorf("%s: field %s is not found", sourceFn, name)
	}

	if !targetField.CanSet() {
		return fmt.Errorf("%s: field %s can not be set", sourceFn, name)
	}

	sourceIDField := reflect.ValueOf(source).Elem().FieldByName("ID")
	if !sourceIDField.IsValid() {
		return fmt.Errorf("%s: source field ID is not found", sourceFn)
	}

	sourceIDKind := sourceIDField.Kind()
	if sourceIDKind != reflect.Int &&
		sourceIDKind != reflect.Int64 &&
		sourceIDKind != reflect.Int32 &&
		sourceIDKind != reflect.Int16 &&
		sourceIDKind != reflect.Int8 &&
		sourceIDKind != reflect.Uint &&
		sourceIDKind != reflect.Uint64 &&
		sourceIDKind != reflect.Uint32 &&
		sourceIDKind != reflect.Uint16 &&
		sourceIDKind != reflect.Uint8 {
		return fmt.Errorf("%s: source field ID is not an integer", sourceFn)
	}

	// TODO: What if targetField is int, but sourceIDField is uint?
	switch sourceIDField.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		targetField.SetInt(sourceIDField.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		targetField.SetUint(sourceIDField.Uint())
	}

	return nil
}

// genAndInsertAss inserts the associations value into the database and returns with the inserted values
func genAndInsertAss(d db.Database, associations map[string][]interface{}, tagToInfo map[string]tagInfo) ([]interface{}, error) {
	if len(tagToInfo) == 0 {
		return nil, errors.New("tagToInfo is not set")
	}

	if len(associations) == 0 {
		return nil, errors.New("inserting associations without any associations")
	}

	result := []interface{}{}
	for name, vals := range associations {
		tableName := tagToInfo[name].tableName

		v, err := d.InsertList(db.InserListParams{StorageName: tableName, Values: vals})
		if err != nil {
			return nil, err
		}

		result = append(result, v...)
	}

	return result, nil
}

// setAssValue sets the value to the associations value
func setAssValue(v interface{}, tagToInfo map[string]tagInfo, index int, sourceFn string) error {
	typeOfV := reflect.TypeOf(v)

	// check if it's a pointer
	if typeOfV.Kind() != reflect.Ptr {
		name := typeOfV.Name()
		return fmt.Errorf("%s: type %s, value %v is not a pointer", sourceFn, name, v)
	}

	name := typeOfV.Elem().Name()
	// check if it's a pointer to a struct
	if typeOfV.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("%s: type %s, value %v is not a pointer to a struct", sourceFn, name, v)
	}

	// check if it's existed in tagToInfo
	if _, ok := tagToInfo[name]; !ok {
		return fmt.Errorf("%s: type %s, value %v is not found at tag", sourceFn, name, v)
	}

	setNonZeroValues(index, v)
	return nil
}

// genTagToInfo generates the map from tag to metadata
func genTagToInfo(dataType reflect.Type) (map[string]tagInfo, error) {
	tagToInfo := map[string]tagInfo{}
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		tag := field.Tag.Get("efactory")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		if len(parts) == 0 {
			return nil, errors.New("tag is in wrong format. It should be efactory:\"<struct_name>,<table_name>\"")
		}

		structName := parts[0]

		var tableName string
		if len(parts) == 2 {
			tableName = parts[1]
		} else {
			tableName = CamelToSnake(structName) + "s"
		}

		tagToInfo[structName] = tagInfo{tableName: tableName, fieldName: field.Name}
	}

	return tagToInfo, nil
}
