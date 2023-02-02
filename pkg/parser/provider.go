package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/samber/lo"
)

const interfaceTemplate = `export interface {{ .Name }} {
	{{- range $k, $v :=.Properties }}
	{{ $k }}: {{ $v }};
	{{- end }}
}
`

type SchemaProvider struct {
	Schemas []*InterfaceDefinition
}

func (s *SchemaProvider) Provide(buf *bytes.Buffer) error {
	compiledTemplate := template.Must(template.New("templates").Parse(interfaceTemplate))

	for i, id := range s.Schemas {
		if i > 0 {
			buf.WriteString("\n")
		}

		if err := compiledTemplate.Execute(buf, id); err != nil {
			return err
		}
	}

	return nil
}

func (r *SchemaProvider) Merge(sp *SchemaProvider) {
	r.Schemas = append(r.Schemas, sp.Schemas...)
	r.Schemas = lo.UniqBy(r.Schemas, func(i *InterfaceDefinition) string {
		return string(i.Hash())
	})
}

func NewSchemaProvider() *SchemaProvider {
	return &SchemaProvider{
		Schemas: []*InterfaceDefinition{},
	}
}

func FromBuffer(buf io.Reader, name string) (*SchemaProvider, error) {
	var input any

	if err := json.NewDecoder(buf).Decode(&input); err != nil {
		return nil, err
	}

	switch x := input.(type) {
	case []any:
		defs, err := parse(x[0].(map[string]any), name)
		if err != nil {
			return nil, err
		}

		return &SchemaProvider{
			Schemas: defs,
		}, nil

	case map[string]any:
		defs, err := parse(x, "Main")
		if err != nil {
			return nil, err
		}

		return &SchemaProvider{
			Schemas: defs,
		}, nil
	}

	return nil, fmt.Errorf("failed to parse schema")
}
