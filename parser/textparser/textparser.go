package textparser

import (
	"fmt"

	"github.com/satyrius/gonx"
)

// TextParser parses variables patterns using config.NamespaceConfig.Format.
type TextParser struct {
	parser *gonx.Parser
}

// NewTextParser returns a new text parser.
func NewTextParser(format string) *TextParser {
	return &TextParser{
		parser: gonx.NewParser(format),
	}
}

// ParseString implements the Parser interface.
func (t *TextParser) ParseString(line string) (map[string]string, error) {
	entry, err := t.parser.ParseString(line)
	if err != nil {
		return nil, fmt.Errorf("text log parsing err: %w", err)
	}

	return entry.Fields(), nil
}
