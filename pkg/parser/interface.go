package parser

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"
)

type InterfaceDefinition struct {
	Name       string
	Properties map[string]interface{}
}

func (d *InterfaceDefinition) Equals(other *InterfaceDefinition) bool {
	dJ, _ := json.Marshal(d)
	otherJ, _ := json.Marshal(other)

	return strings.EqualFold(string(dJ), string(otherJ))
}

func (d *InterfaceDefinition) Hash() []byte {
	dJ, _ := json.Marshal(d)

	return sha1.New().Sum(dJ)
}

func parse(input map[string]any, name string) ([]*InterfaceDefinition, error) {
	result := make([]*InterfaceDefinition, 0)

	props := map[string]interface{}{}
	for key, value := range input {
		t := getType(value, key)

		props[key] = t

		if t.Values != nil {
			r, ok := t.Values.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("cannot convert %v to map[string]any", t.Values)
			}
			id, err := parse(r, t.Name)
			if err != nil {
				return nil, fmt.Errorf("cannot parse %v: %w", key, err)
			}

			result = append(result, id...)
		}
	}

	result = append(result, &InterfaceDefinition{Name: name, Properties: props})

	return result, nil
}
