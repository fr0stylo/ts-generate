package parser

import (
	"fmt"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TypeDefinition struct {
	Name   string
	Values any
	Type   *string
}

func (t TypeDefinition) String() string {
	if t.Type == nil {
		return t.Name
	}

	return fmt.Sprintf(*t.Type, t.Name)
}

func (t TypeDefinition) Extend() string {
	if t.Type == nil {
		return "%s"
	}

	return *t.Type
}

func refString(ref string) *string {
	return &ref
}

func getType(value any, key string) *TypeDefinition {
	switch x := value.(type) {
	case time.Time:
		return &TypeDefinition{Name: "Date", Values: nil}
	case int:
		return &TypeDefinition{Name: "number", Values: nil}
	case string:
		return &TypeDefinition{Name: "string", Values: nil}
	case bool:
		return &TypeDefinition{Name: "boolean", Values: nil}
	case float64:
		return &TypeDefinition{Name: "number", Values: nil}
	case []any:
		if len(x) > 0 {
			typeDef := getType(x[0], key)

			return &TypeDefinition{Name: typeDef.Name, Values: typeDef.Values, Type: refString(fmt.Sprintf("Array<%s>", typeDef.Extend()))}
		}
	case any:
		return &TypeDefinition{Name: cases.Title(language.English, cases.Compact).String(key), Values: value, Type: refString("Partial<%s>")}
	}

	return &TypeDefinition{Name: "unknown"}
}
