package reflectutil

import "reflect"

type StructField struct {
	Index       int
	Tag         reflect.StructTag
	IsMandatory bool
	IsID        bool
	IsAltID     bool
}

var (
	structFieldsCache = make(map[reflect.Type]map[string]StructField)
	invalidType       = reflect.TypeOf(nil)
)

// GetStructFields returns the (filtered) fields of the given struct. It caches
// the result in order to speed up subsequent calls.
func GetStructFields(structVal reflect.Value) map[string]StructField {
	typ := structVal.Type()
	if cached, ok := structFieldsCache[typ]; ok {
		return cached
	}

	flds := make([]reflect.StructField, 0, structVal.NumField())
	for i := 0; i < structVal.NumField(); i++ {
		// the struct field `PkgPath` is empty for exported fields
		if structVal.Type().Field(i).PkgPath != "" {
			continue
		}
		flds = append(flds, structVal.Type().Field(i))
	}

	structFlds := make(map[string]StructField, len(flds))
	for _, fld := range flds {
		sf := StructField{
			Index:       fld.Index[0],
			Tag:         fld.Tag,
			IsMandatory: fld.Tag.Get("mandatory") == "true",
			IsID:        fld.Tag.Get("id") == "true",
			IsAltID:     fld.Tag.Get("altID") == "true",
		}
		structFlds[fld.Name] = sf
	}

	// cache fields info
	structFieldsCache[typ] = structFlds

	return structFlds
}

// Indirect dereferences pointers or interfaces until a non pointer/interface
// value is found and returns the result.
func Indirect(v reflect.Value) reflect.Value {
	knd := v.Kind()
	if knd == reflect.Pointer || knd == reflect.Interface {
		v = Indirect(v.Elem())
	}
	return v
}

// ZeroFromType returns the zero value of the given type.
func ZeroFromType(t reflect.Type) reflect.Value {
	if t == invalidType {
		return reflect.Value{}
	}
	vTop := reflect.New(t)
	v := vTop.Elem()
	for {
		switch t.Kind() {
		case reflect.Pointer:
			t = t.Elem()
			zv := reflect.New(t)
			if zv.IsValid() {
				v.Set(zv)
			}
			v = zv.Elem()
			continue
		}
		break
	}
	return vTop.Elem()
}

func ZeroFromValue(from reflect.Value) reflect.Value {
	if !from.IsValid() {
		return reflect.Value{}
	}
	t := from.Type()
	vTop := reflect.New(t)
	v := vTop.Elem()
	for {
		switch from.Kind() {
		case reflect.Pointer:
			from = from.Elem()
			t = t.Elem()
			zv := reflect.New(t)
			if zv.IsValid() {
				v.Set(zv)
			}
			v = zv.Elem()
			continue

		case reflect.Interface:
			from = from.Elem()
			t = from.Type()
			continue
		}
		break
	}
	return vTop.Elem()
}

// IsNil returns true if the given value is nil.
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return rv.IsNil()
	}
	return false
}
