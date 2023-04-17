package schema

import (
	"fmt"
	"reflect"
	"strings"
)

func NewBuilderDefault() Builder {
	return Builder{
		Options: &Options{
			ExploreNilStruct: true,
			PreferJsonTag:    true,
		},
	}
}

func NewBuilder(opts *Options) Builder {
	return Builder{
		Options: opts,
	}
}

type Options struct {
	ExploreNilStruct bool
	PreferJsonTag    bool
}

type Builder struct {
	Options *Options
}

func (b *Builder) GetSchema(obj interface{}) (*Property, error) {
	if obj == nil {
		return nil, nil
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	visited := make(map[string][][]int)
	currentPath := []int{1}
	return b.inspect(t, v, visited, currentPath)
}

func (b *Builder) inspect(t reflect.Type, v reflect.Value, visited map[string][][]int, currentPath []int) (*Property, error) {
	switch t.Kind() {
	case reflect.Interface:
		if v.Interface() == nil {
			return nil, nil
		}
		return b.inspect(v.Elem().Type(), v.Elem(), visited, currentPath)
	case reflect.Ptr:
		if !v.IsValid() {
			return b.inspect(t.Elem(), reflect.Value{}, visited, currentPath)
		}
		return b.inspect(t.Elem(), v.Elem(), visited, currentPath)
	case reflect.String:
		return &Property{
			Type:  PropType_STRING,
			Value: b.valueString(v),
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Property{
			Type:  PropType_INTEGER,
			Value: b.valueString(v),
		}, nil
	case reflect.Float32, reflect.Float64:
		return &Property{
			Type:  PropType_NUMBER,
			Value: b.valueString(v),
		}, nil
	case reflect.Bool:
		return &Property{
			Type:  PropType_BOOLEAN,
			Value: b.valueString(v),
		}, nil
	case reflect.Struct:
		if !v.IsValid() && !b.Options.ExploreNilStruct {
			return nil, nil
		}
		structType := t.String()
		props := make([]Property, 0)
		var err error
		visitedPath := (visited)[structType]
		if visitedPath == nil || len(visitedPath) == 0 {
			(visited)[structType] = [][]int{currentPath}
			props, err = b.inspectStruct(t, v, visited, currentPath)
			if err != nil {
				return nil, err
			}
			return &Property{
				Type:       PropType_OBJECT,
				Properties: props,
			}, nil
		} else {
			for _, v := range visitedPath {
				prefixCandidate := currentPath[0:len(v)]
				if reflect.DeepEqual(prefixCandidate, v) {
					//stopping the further exploration when a cyclic dependency is found in the tree branch
					return &Property{
						Type:       PropType_OBJECT,
						Properties: props,
					}, nil
				}
			}
		}
		visitedPath = append(visitedPath, currentPath)
		(visited)[structType] = visitedPath
		props, err = b.inspectStruct(t, v, visited, currentPath)
		if err != nil {
			return nil, err
		}
		return &Property{
			Type:       PropType_OBJECT,
			Properties: props,
		}, nil
	case reflect.Slice:
		props := make([]Property, 0)
		if !v.IsValid() || v.Len() == 0 {
			prop, err := b.inspect(t.Elem(), reflect.Value{}, visited, currentPath)
			if err != nil {
				return nil, err
			}
			if prop != nil {
				props = append(props, *prop)
			}
		} else {
			for i := 0; i < v.Len(); i++ {
				prop, err := b.inspect(t.Elem(), v.Index(i), visited, currentPath)
				if err != nil {
					return nil, err
				}
				props = append(props, *prop)
			}
		}
		return &Property{
			Type:       PropType_ARRAY,
			Properties: props,
		}, nil
	case reflect.Map:
		panic("not implemented")
	default:
		panic("unknown type")
	}
}

func (b *Builder) inspectStruct(t reflect.Type, v reflect.Value, visited map[string][][]int, currentPath []int) ([]Property, error) {
	props := make([]Property, 0)
	for i := 0; i < t.NumField(); i++ {
		nextPath := make([]int, len(currentPath))
		copy(nextPath, currentPath)
		nextPath = append(nextPath, i)

		_field := t.Field(i)
		var _value reflect.Value
		if v.IsValid() {
			_value = v.Field(i)
		} else {
			_value = reflect.Value{}
		}
		prop, err := b.inspect(_field.Type, _value, visited, nextPath)
		if err != nil {
			return nil, err
		}
		if prop == nil {
			continue
		}
		prop = prop.
			WithName(b.structFieldName(_field))
		props = append(props, *prop)
	}
	return props, nil
}

func (b *Builder) structFieldName(sf reflect.StructField) string {
	if b.Options.PreferJsonTag {
		jsonTag := strings.Split(sf.Tag.Get("json"), ",")[0]
		if jsonTag != "" {
			return jsonTag
		}
	}
	return sf.Name
}

func (b *Builder) valueString(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
