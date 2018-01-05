package gonx

import (
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestReader(t *testing.T) {
	Convey("Test Reader", t, func() {
		format := "$remote_addr [$time_local] \"$request\""

		Convey("Test valid input", func() {
			file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
			reader := NewReader(file, format)
			So(reader.entries, ShouldBeNil)

			expected := NewEntry(Fields{
				"remote_addr": "89.234.89.123",
				"time_local":  "08/Nov/2013:13:39:18 +0000",
				"request":     "GET /api/foo/bar HTTP/1.1",
			})

			// Read entry from incoming channel
			entry, err := reader.Read()
			So(err, ShouldBeNil)
			So(entry, ShouldResemble, expected)

			// It was only one line, nothing to read
			_, err = reader.Read()
			So(err, ShouldEqual, io.EOF)
		})

		Convey("Test long line", func() {
			longStr := RandString(64 * 1024)
			file := strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET ` + longStr + ` HTTP/1.1"`)
			reader := NewReader(file, format)
			_, err := reader.Read()
			So(err, ShouldBeNil)
		})
	})
}
