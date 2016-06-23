package gonx

import (
	"io"
)

// Log file reader. Use specific constructors to create it.
type Reader struct {
	file    io.Reader
	parser  *Parser
	entries chan *Entry
}

// Creates reader for custom log format.
func NewReader(logFile io.Reader, format string) *Reader {
	return &Reader{
		file:   logFile,
		parser: NewParser(format),
	}
}

// Creates reader for nginx log format. Nginx config parser will be used
// to get particular format from the conf file.
func NewNginxReader(logFile io.Reader, nginxConf io.Reader, formatName string) (reader *Reader, err error) {
	parser, err := NewNginxParser(nginxConf, formatName)
	if err != nil {
		return nil, err
	}
	reader = &Reader{
		file:   logFile,
		parser: parser,
	}
	return
}

// Get next parsed Entry from the log file. Return EOF if there is no Entries to read.
func (r *Reader) Read() (entry *Entry, err error) {
	if r.entries == nil {
		r.entries = MapReduce(r.file, r.parser, new(ReadAll))
	}
	entry, ok := <-r.entries
	if !ok {
		err = io.EOF
	}
	return
}
