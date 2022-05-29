package qdoc

import (
	"github.com/ThilinaTLM/quick-doc/schema"
	"github.com/getkin/kin-openapi/openapi3"
)

// Schema is document data scheme configuration
func Schema(value interface{}) SchemaConfig {
	return SchemaConfig{
		Object: value,
		builder: schema.NewBuilder(&schema.Options{
			ExploreNilStruct: false,
			PreferJsonTag:    true,
		}),
	}
}

type SchemaConfig struct {
	Object  interface{}
	builder schema.Builder
}

func (sb *SchemaConfig) toOpenAPI() *openapi3.Schema {
	prop, err := sb.builder.GetSchema(sb.Object)
	if err != nil {
		return nil
	}

	return propToOpenAPI(prop)
}

func propToOpenAPI(prop *schema.Property) *openapi3.Schema {
	if prop == nil {
		return openapi3.NewSchema()
	}

	switch prop.Type {
	case schema.PropType_STRING:
		return &openapi3.Schema{
			Type:        "string",
			Example:     prop.Value,
			Title:       prop.Name,
			Description: prop.Description,
		}
	case schema.PropType_INTEGER:
		return &openapi3.Schema{
			Type:        "integer",
			Example:     prop.Value,
			Title:       prop.Name,
			Description: prop.Description,
		}
	case schema.PropType_NUMBER:
		return &openapi3.Schema{
			Type:        "number",
			Example:     prop.Value,
			Title:       prop.Name,
			Description: prop.Description,
		}
	case schema.PropType_BOOLEAN:
		return &openapi3.Schema{
			Type:        "boolean",
			Example:     prop.Value,
			Title:       prop.Name,
			Description: prop.Description,
		}
	case schema.PropType_ARRAY:
		return &openapi3.Schema{
			Type: "array",
			Items: &openapi3.SchemaRef{
				Value: propToOpenAPI(&prop.Properties[0]),
			},
			Title:       prop.Name,
			Description: prop.Description,
		}
	case schema.PropType_OBJECT:
		properties := make(map[string]*openapi3.SchemaRef)

		for _, p := range prop.Properties {
			properties[p.Name] = &openapi3.SchemaRef{
				Value: propToOpenAPI(&p),
			}
		}
		return &openapi3.Schema{
			Type:        "object",
			Properties:  properties,
			Title:       prop.Name,
			Description: prop.Description,
		}
	default:
		panic("unknown type")
	}
}