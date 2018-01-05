package gonx

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReducer(t *testing.T) {
	Convey("Test process input channel with reducers", t, func() {
		input := make(chan *Entry, 10)

		Convey("ReadAll reducer", func() {
			reducer := new(ReadAll)

			// Prepare import channel
			entry := NewEmptyEntry()
			input <- entry
			close(input)

			output := make(chan *Entry, 1) // Make it buffered to avoid deadlock
			reducer.Reduce(input, output)

			// ReadAll reducer writes input channel to the output
			result, ok := <-output
			So(ok, ShouldBeTrue)
			So(result, ShouldEqual, entry)
		})

		Convey("With filled input channel", func() {
			// Prepare import channel
			input <- NewEntry(Fields{
				"uri":  "/asd/fgh",
				"host": "alpha.example.com",
				"foo":  "1",
				"bar":  "2",
				"baz":  "3",
			})
			input <- NewEntry(Fields{
				"uri":  "/zxc/vbn",
				"host": "beta.example.com",
				"foo":  "4",
				"bar":  "5",
				"baz":  "6",
			})
			input <- NewEntry(Fields{
				"uri":  "/ijk/lmn",
				"host": "beta.example.com",
				"foo":  "7",
				"bar":  "8",
				"baz":  "9",
			})
			close(input)
			total := float64(len(input))

			output := make(chan *Entry, 10) // Make it buffered to avoid deadlock

			Convey("Count reducer", func() {
				reducer := new(Count)
				reducer.Reduce(input, output)

				result, ok := <-output
				So(ok, ShouldBeTrue)
				count, err := result.FloatField("count")
				So(err, ShouldBeNil)
				So(count, ShouldEqual, total)
			})

			Convey("Sum reducer", func() {
				reducer := &Sum{[]string{"foo", "bar"}}
				reducer.Reduce(input, output)

				result, ok := <-output
				So(ok, ShouldBeTrue)

				value, err := result.FloatField("foo")
				So(err, ShouldBeNil)
				So(value, ShouldEqual, 1+4+7)

				value, err = result.FloatField("bar")
				So(err, ShouldBeNil)
				So(value, ShouldEqual, 2+5+8)

				_, err = result.Field("buz")
				So(err, ShouldNotBeNil)
			})

			Convey("Avg reducer", func() {
				reducer := &Avg{[]string{"foo", "bar"}}
				reducer.Reduce(input, output)

				result, ok := <-output
				So(ok, ShouldBeTrue)

				value, err := result.FloatField("foo")
				So(err, ShouldBeNil)
				So(value, ShouldEqual, (1+4+7)/total)

				value, err = result.FloatField("bar")
				So(err, ShouldBeNil)
				So(value, ShouldEqual, (2+5+8)/total)

				_, err = result.Field("buz")
				So(err, ShouldNotBeNil)
			})

			Convey("Chain reducer", func() {
				reducer := NewChain(&Avg{[]string{"foo", "bar"}}, &Count{})
				So(len(reducer.reducers), ShouldEqual, 2)
				reducer.Reduce(input, output)

				result, ok := <-output
				So(ok, ShouldBeTrue)

				value, err := result.FloatField("foo")
				So(err, ShouldBeNil)
				So(value, ShouldEqual, (1+4+7)/total)

				value, err = result.FloatField("bar")
				So(err, ShouldBeNil)
				So(value, ShouldEqual, (2+5+8)/total)

				count, err := result.FloatField("count")
				So(err, ShouldBeNil)
				So(count, ShouldEqual, total)

				_, err = result.Field("buz")
				So(err, ShouldNotBeNil)
			})

			Convey("Group reducer", func() {
				reducer := NewGroupBy(
					// Fields to group by
					[]string{"host"},
					// Result reducers
					&Sum{[]string{"foo", "bar"}},
					new(Count),
				)
				So(len(reducer.reducers), ShouldEqual, 2)
				reducer.Reduce(input, output)

				// Collect result entries from output channel to the map, because reading
				// from channel can be in any order, it depends on each reducer processing
				resultMap := make(map[string]*Entry)
				for result := range output {
					value, err := result.Field("host")
					So(err, ShouldBeNil)
					resultMap[value] = result
				}
				So(len(resultMap), ShouldEqual, 2)

				// Read and assert first group result
				result := resultMap["alpha.example.com"]

				floatVal, err := result.FloatField("foo")
				So(err, ShouldBeNil)
				So(floatVal, ShouldEqual, 1)

				floatVal, err = result.FloatField("bar")
				So(err, ShouldBeNil)
				So(floatVal, ShouldEqual, 2)

				count, err := result.FloatField("count")
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)

				// Read and assert second group result
				result = resultMap["beta.example.com"]

				floatVal, err = result.FloatField("foo")
				So(err, ShouldBeNil)
				So(floatVal, ShouldEqual, 4+7)

				floatVal, err = result.FloatField("bar")
				So(err, ShouldBeNil)
				So(floatVal, ShouldEqual, 5+8)

				count, err = result.FloatField("count")
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 2)
			})
		})
	})
}
