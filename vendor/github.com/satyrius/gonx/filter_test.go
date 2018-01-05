package gonx

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestFilter(t *testing.T) {
	Convey("Test Datetime filter", t, func() {
		start := time.Date(2015, time.February, 2, 2, 2, 2, 0, time.UTC)
		end := time.Date(2015, time.May, 5, 5, 5, 5, 0, time.UTC)

		jan := NewEntry(Fields{"timestamp": "2015-01-01T01:01:01Z", "foo": "12"})
		feb := NewEntry(Fields{"timestamp": "2015-02-02T02:02:02Z", "foo": "34"})
		mar := NewEntry(Fields{"timestamp": "2015-03-03T03:03:03Z", "foo": "56"})
		apr := NewEntry(Fields{"timestamp": "2015-04-04T04:04:04Z", "foo": "78"})
		may := NewEntry(Fields{"timestamp": "2015-05-05T05:05:05Z", "foo": "90"})

		Convey("Filter Entry", func() {
			Convey("Start and end", func() {
				filter := &Datetime{
					Field:  "timestamp",
					Format: time.RFC3339,
					Start:  start,
					End:    end,
				}

				// entries is out of datetime range
				So(filter.Filter(jan), ShouldBeNil)
				So(filter.Filter(may), ShouldBeNil)

				// entry's timestamp meets filter condition
				So(filter.Filter(feb), ShouldEqual, feb)
			})

			Convey("Start only", func() {
				filter := &Datetime{
					Field:  "timestamp",
					Format: time.RFC3339,
					Start:  start,
				}

				// entry is out of datetime range
				So(filter.Filter(jan), ShouldBeNil)

				// entry's timestamp meets filter condition
				So(filter.Filter(feb), ShouldEqual, feb)
			})

			Convey("End only", func() {
				filter := &Datetime{
					Field:  "timestamp",
					Format: time.RFC3339,
					End:    end,
				}

				// entry's timestamp meets filter condition
				So(filter.Filter(jan), ShouldEqual, jan)

				// entry is out of datetime range
				So(filter.Filter(may), ShouldBeNil)
			})
		})

		Convey("Deal with input channel", func() {
			input := make(chan *Entry, 5)
			input <- jan
			input <- feb
			input <- mar
			input <- apr
			input <- may
			close(input)

			filter := &Datetime{
				Field:  "timestamp",
				Format: time.RFC3339,
				Start:  start,
				End:    end,
			}
			output := make(chan *Entry, 5) // Make it buffered to avoid deadlock

			Convey("Reduce channel", func() {
				filter.Reduce(input, output)

				expected := []string{
					"'timestamp'=2015-02-02T02:02:02Z;'foo'=34",
					"'timestamp'=2015-03-03T03:03:03Z;'foo'=56",
					"'timestamp'=2015-04-04T04:04:04Z;'foo'=78",
				}
				results := []string{}

				for result := range output {
					results = append(
						results,
						result.FieldsHash([]string{"timestamp", "foo"}),
					)
				}
				So(results, ShouldResemble, expected)
			})

			Convey("Filter channel", func() {
				chain := NewChain(filter, &Avg{[]string{"foo"}}, &Count{})
				chain.Reduce(input, output)

				result, ok := <-output
				So(ok, ShouldBeTrue)

				value, err := result.FloatField("foo")
				So(err, ShouldBeNil)
				So(value, ShouldEqual, (34+56+78)/3)

				count, err := result.FloatField("count")
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 3)

				_, err = result.Field("bar")
				So(err, ShouldNotBeNil)
			})
		})
	})
}
