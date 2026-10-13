package objectrepo

import (
	"reflect"
	"sync"
)

type object struct {
	objectType reflect.Type
	object     any
	objectTag  string
	fields     []field
}

type field struct {
	parent reflect.Value
	field  reflect.StructField
	tag    string
}

type ObjectRepository struct {
	m       sync.Mutex
	objects map[string]object
}

func New() *ObjectRepository {
	return &ObjectRepository{
		objects: make(map[string]object, 8),
	}
}

func (or *ObjectRepository) getByTagAndType(t reflect.Type, tag string) any {
	o, ok := or.objects[tag]
	if ok {
		if isType(t, o.objectType) {
			if tag == "" || tag == o.objectTag {
				return o.object
			}
		}
	}

	return nil
}

func (or *ObjectRepository) getByType(t reflect.Type) any {
	for _, o := range or.objects {
		if isType(t, o.objectType) {
			return o.object
		}
	}

	return nil
}

func isType(field, object reflect.Type) bool {
	switch field.Kind() {
	case reflect.Pointer:
		{
			return field == object
		}
	case reflect.Interface:
		{
			return object.Implements(field)
		}
	}

	return false
}

// Return first instance for the type
func (or *ObjectRepository) Find(t reflect.Type, tag string) any {
	if t == nil {
		return nil
	}

	if tag != "" {
		return or.getByTagAndType(t, tag)
	}

	return or.getByType(t)

}

// Return all instances for the type
func (or *ObjectRepository) Lookup(t reflect.Type, tag string) []any {
	res := make([]any, 0, len(or.objects))
	if t != nil {
		switch t.Kind() {
		case reflect.Pointer:
			{
				for _, o := range or.objects {
					if o.objectType == t {
						if tag == "" || tag == o.objectTag {
							res = append(res, o.object)
						}
					}
				}

			}
		case reflect.Interface:
			{
				for _, o := range or.objects {
					if o.objectType.Implements(t) {
						if tag == "" || tag == o.objectTag {
							res = append(res, o.object)
						}
					}
				}
			}
		}
	}
	return res
}

// !!! Биндятся только поинтеры для структур !!!
func (or *ObjectRepository) Bind(value any) {
	if value == nil {
		panic("value must be defined")
	}

	if reflect.ValueOf(value).Type().Kind() != reflect.Pointer {
		panic("value must be pointer")
	}

	v := reflect.ValueOf(value).Elem()

	if v.Type().Kind() != reflect.Struct {
		panic("value must be pointer to struct")
	}

	tag := v.Type().Name()

	if _, ok := or.objects[tag]; ok {
		panic("an instance with that tag already exists")
	}

	fl := or.collectFields(v)

	or.m.Lock()
	or.objects[tag] = object{
		objectType: reflect.TypeOf(value),
		object:     value,
		objectTag:  tag,
		fields:     fl,
	}
	or.m.Unlock()
}

// !!! Собираются только атрибуты первого уровня !!!
func (or *ObjectRepository) collectFields(v reflect.Value) []field {
	t := v.Type()
	fl := make([]field, 0, t.NumField())
	for _, f := range reflect.VisibleFields(t) {
		tv, ok := f.Tag.Lookup("component")
		if ok {

			ftk := or.getFieldType(f)

			// !!! Поле структуры для заполнения может быть либо интерфейсом, либо указателем на структуру !!!
			if ftk != reflect.Invalid {
				fi := field{
					parent: v,
					field:  f,
					tag:    tv,
				}
				or.setValueForField(fi)
				fl = append(fl, fi)
			}
		}
	}

	return fl
}

func (or *ObjectRepository) getFieldType(f reflect.StructField) reflect.Kind {
	switch k := f.Type.Kind(); k {
	case reflect.Interface:
		{
			return k
		}
	case reflect.Pointer:
		{
			et := f.Type.Elem().Kind()
			if et == reflect.Struct {
				return et
			}
		}
	case reflect.Slice:
		{
			et := f.Type.Elem().Kind()
			if et == reflect.Interface {
				return et
			}
			if et == reflect.Pointer {
				ept := f.Type.Elem().Elem().Kind()
				if ept == reflect.Struct {
					return et
				}
			}
		}
	}

	return reflect.Invalid
}

func (or *ObjectRepository) Apply() {
	or.m.Lock()
	defer or.m.Unlock()
	for _, o := range or.objects {
		for _, f := range o.fields {
			or.setValueForField(f)
		}
	}
}

func (or *ObjectRepository) setValueForField(f field) {
	switch f.field.Type.Kind() {
	case reflect.Slice:
		{
			valueForField := or.Lookup(f.field.Type.Elem(), f.tag)
			f := f.parent.FieldByName(f.field.Name)
			sliceVal := reflect.MakeSlice(f.Type(), len(valueForField), len(valueForField))
			for i, s := range valueForField {
				v := reflect.ValueOf(s)
				sliceVal.Index(i).Set(v)
			}

			f.Set(sliceVal)
		}
	default:
		{
			valueForField := or.Find(f.field.Type, f.tag)
			if valueForField != nil {
				v := reflect.ValueOf(valueForField)
				f := f.parent.FieldByName(f.field.Name)
				f.Set(v)
			}
		}
	}

}
