package gonx

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// StringParser is the interface that wraps the ParseString method.
type StringParser interface {
	ParseString(line string) (entry *Entry, err error)
}

// Log record parser. Use specific constructors to initialize it.
type Parser struct {
	format string
	regexp *regexp.Regexp
}

// Returns a new Parser, use given log format to create its internal
// strings parsing regexp.
func NewParser(format string) *Parser {
	re := regexp.MustCompile(`\\\$([a-z_]+)(\\?(.))`).ReplaceAllString(
		regexp.QuoteMeta(format+" "), "(?P<$1>[^$3]*)$2")
	return &Parser{format, regexp.MustCompile(fmt.Sprintf("^%v$", strings.Trim(re, " ")))}
}

// Parse log file line using internal format regexp. If line do not match
// given format an error will be returned.
func (parser *Parser) ParseString(line string) (entry *Entry, err error) {
	re := parser.regexp
	fields := re.FindStringSubmatch(line)
	if fields == nil {
		err = fmt.Errorf("access log line '%v' does not match given format '%v'", line, re)
		return
	}

	// Iterate over subexp foung and fill the map record
	entry = NewEmptyEntry()
	for i, name := range re.SubexpNames() {
		if i == 0 {
			continue
		}
		entry.SetField(name, fields[i])
	}
	return
}

// NewNginxParser parse nginx conf file to find log_format with given name and
// returns parser for this format. It returns an error if cannot find the needle.
func NewNginxParser(conf io.Reader, name string) (parser *Parser, err error) {
	scanner := bufio.NewScanner(conf)
	re := regexp.MustCompile(fmt.Sprintf(`^\s*log_format\s+%v\s+(.+)\s*$`, name))
	found := false
	var format string
	for scanner.Scan() {
		var line string
		if !found {
			// Find a log_format definition
			line = scanner.Text()
			formatDef := re.FindStringSubmatch(line)
			if formatDef == nil {
				continue
			}
			found = true
			line = formatDef[1]
		} else {
			line = scanner.Text()
		}
		// Look for a definition end
		re = regexp.MustCompile(`^\s*(.*?)\s*(;|$)`)
		lineSplit := re.FindStringSubmatch(line)
		if l := len(lineSplit[1]); l > 2 {
			format += lineSplit[1][1 : l-1]
		}
		if lineSplit[2] == ";" {
			break
		}
	}
	if !found {
		err = fmt.Errorf("`log_format %v` not found in given config", name)
	} else {
		err = scanner.Err()
	}
	parser = NewParser(format)
	return
}
