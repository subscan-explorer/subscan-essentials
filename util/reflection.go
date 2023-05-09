package util

import (
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/exp/slog"
)

func MapInterfaceAsStruct[T any](iface interface{}) (T, error) {
	var t T
	typ := reflect.TypeOf(t)
	strukt, err := mapInterfaceAsStruct(typ, iface)
	if err != nil {
		return t, err
	}
	return strukt.Interface().(T), nil
}

func checkAssignable(a, b ValAndType) error {
	if !a.Type.AssignableTo(b.Type) {
		return fmt.Errorf("type mismatch. type of %+v is %+v, but %+v is %+v", a.Val, a.Type, b.Val, b.Type)
	}
	return nil
}

type ValAndType struct {
	Val  interface{}
	Type reflect.Type
}

func valAndType(val interface{}, typ reflect.Type) ValAndType {
	return ValAndType{Val: val, Type: typ}
}

func valType(val interface{}) ValAndType {
	return valAndType(val, reflect.TypeOf(val))
}

func mapInterfaceAsStruct(typ reflect.Type, iface interface{}) (reflect.Value, error) {
	if typ.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("invalid type")
	}

	ret := reflect.New(typ).Elem()

	recurse := func(field reflect.StructField, v interface{}) error {
		if field.Type.Kind() == reflect.Slice && reflect.TypeOf(v).AssignableTo(reflect.TypeOf([]interface{}{})) {
			elem := field.Type.Elem()

			slice := reflect.MakeSlice(field.Type, 0, 0)
			for _, val := range v.([]interface{}) {
				if elem.Kind() == reflect.Struct && reflect.TypeOf(val).AssignableTo(reflect.TypeOf(map[string]interface{}{})) {
					v, err := mapInterfaceAsStruct(elem, val)
					if err != nil {
						return err
					}
					slice = reflect.Append(slice, v)
				} else {
					if err := checkAssignable(valAndType(field.Name, elem), valType(v)); err != nil {
						return err
					}
					slice = reflect.Append(slice, reflect.ValueOf(val))
				}
			}
			if !slice.Type().AssignableTo(field.Type) {
				return fmt.Errorf("type mismatch. type of %+v is %+v, but %+v is %+v", field.Name, field.Type, slice, slice.Type())
			}
			ret.FieldByIndex(field.Index).Set(slice)
			return nil
		} else if field.Type.Kind() == reflect.Struct && reflect.TypeOf(v).AssignableTo(reflect.TypeOf(map[string]interface{}{})) {
			v, err := mapInterfaceAsStruct(field.Type, v)
			if err != nil {
				return err
			}
			ret.FieldByIndex(field.Index).Set(v)
			return nil
		} else {
			if !reflect.TypeOf(v).AssignableTo(field.Type) {
				return fmt.Errorf("invalid type. type of %+v is %+v, but %+v is %+v", field.Name, field.Type, v, reflect.TypeOf(v))
			}
			ret.FieldByIndex(field.Index).Set(reflect.ValueOf(v))
			return nil
		}
	}

	if asMap, ok := iface.(map[string]interface{}); ok {
		for i := 0; i < typ.NumField(); i++ {
			fieldTyp := typ.Field(i)
			if val, ok := asMap[fieldTyp.Name]; ok {
				if err := recurse(fieldTyp, val); err != nil {
					return reflect.Value{}, err
				}
			} else {
				if rename, ok := fieldTyp.Tag.Lookup("json"); ok {
					if val, ok := asMap[rename]; ok {
						if err := recurse(fieldTyp, val); err != nil {
							return reflect.Value{}, err
						}
						continue
					} else {
						slog.Warn("missing field", "name", fieldTyp.Name, "fromMap", asMap)
					}
				} else {
					slog.Warn("missing field", "name", fieldTyp.Name, "fromMap", asMap)
					continue
				}
			}
		}
	} else {
		return reflect.Value{}, errors.New("invalid interface")
	}

	return ret, nil
}
