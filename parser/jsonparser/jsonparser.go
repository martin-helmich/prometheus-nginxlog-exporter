package jsonparser

import (
	"encoding/json"
	"fmt"
)

// JsonParser parse a JSON string.
type JsonParser struct{}

// NewJsonParser returns a new json parser.
func NewJsonParser() *JsonParser {
	return &JsonParser{}
}

// ParseString implements the Parser interface.
// The value in the map is not necessarily a string, so it needs to be converted.
func (j *JsonParser) ParseString(line string) (map[string]string, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(line), &parsed)
	if err != nil {
		return nil, fmt.Errorf("json log parsing err: %w", err)
	}

	fields := make(map[string]string, len(parsed))
	for k, v := range parsed {
		if s, ok := v.(string); ok {
			fields[k] = s
		} else {
			fields[k] = fmt.Sprintf("%v", v)
		}
	}
	return fields, nil
}
