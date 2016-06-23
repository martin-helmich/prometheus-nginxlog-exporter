package gonx

import "time"

// Filter interface for Entries channel limiting.
//
// Filter method should accept *Entry and return *Entry if it meets
// filter condition, otherwise it returns nil.
type Filter interface {
	Reducer
	Filter(*Entry) *Entry
}

// Implements Filter interface to filter Entries with timestamp fields within
// the specified datetime interval.
type Datetime struct {
	Field  string
	Format string
	Start  time.Time
	End    time.Time
}

// Check field value to be in desired datetime range.
func (i *Datetime) Filter(entry *Entry) (validEntry *Entry) {
	val, err := entry.Field(i.Field)
	if err != nil {
		// TODO handle error
		return
	}
	t, err := time.Parse(i.Format, val)
	if err != nil {
		// TODO handle error
		return
	}
	if i.withinBounds(t) {
		validEntry = entry
	}
	return
}

// Reducer interface too. Go through input and apply Filter.
func (i *Datetime) Reduce(input chan *Entry, output chan *Entry) {
	for entry := range input {
		if valid := i.Filter(entry); valid != nil {
			output <- valid
		}
	}
	close(output)
}

func (i *Datetime) withinBounds(t time.Time) bool {
	if t.Equal(i.Start) {
		return true
	}
	if t.After(i.Start) && t.Before(i.End) {
		return true
	}
	return false
}
