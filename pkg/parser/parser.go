package parser

import (
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/config"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/parser/jsonparser"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/parser/textparser"
)

// Parser parses a line of log to a map[string]string.
type Parser interface {
	ParseString(line string) (map[string]string, error)
}

// NewParser returns a Parser with the given config.NamespaceConfig.
func NewParser(nsCfg config.NamespaceConfig) Parser {
	switch nsCfg.Parser {
	case "text":
		return textparser.NewTextParser(nsCfg.Format)
	case "json":
		return jsonparser.NewJsonParser()
	default:
		return textparser.NewTextParser(nsCfg.Format)
	}
}
