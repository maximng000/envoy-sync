package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// SchemaField describes a single expected env key with optional constraints.
type SchemaField struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	Secret   bool   `json:"secret,omitempty"`
}

// Schema holds a collection of field definitions for an env file.
type Schema struct {
	Fields []SchemaField `json:"fields"`
}

// LoadSchema reads and parses a JSON schema file from disk.
func LoadSchema(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read schema: %w", err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse schema: %w", err)
	}
	return &s, nil
}

// SchemaViolation describes a single rule violation found during schema check.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// CheckSchema validates entries against the schema and returns any violations.
func CheckSchema(entries []Entry, s *Schema) []SchemaViolation {
	index := make(map[string]string, len(entries))
	for _, e := range entries {
		index[e.Key] = e.Value
	}

	var violations []SchemaViolation
	for _, field := range s.Fields {
		val, exists := index[field.Key]
		if !exists {
			if field.Required {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: "required key is missing"})
			}
			continue
		}
		if field.Required && val == "" {
			violations = append(violations, SchemaViolation{Key: field.Key, Message: "required key has empty value"})
		}
		if field.Pattern != "" {
			matched, err := regexp.MatchString(field.Pattern, val)
			if err != nil {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err)})
			} else if !matched {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("value does not match pattern %q", field.Pattern)})
			}
		}
	}
	return violations
}
