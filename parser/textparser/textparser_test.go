package textparser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTextParse(t *testing.T) {
	parser := NewTextParser(`[$time_local] $request_method "$request" $request_length $body_bytes_sent $status $request_time $upstream_response_time`)
	line := `[03/Feb/2021:11:22:33 +0800] GET "GET /order/2145 HTTP/1.1" 123 518 200 0.544 0.543`

	got, err := parser.ParseString(line)
	require.NoError(t, err)

	want := map[string]string{
		"time_local":             "03/Feb/2021:11:22:33 +0800",
		"request_time":           "0.544",
		"request_length":         "123",
		"upstream_response_time": "0.543",
		"status":                 "200",
		"body_bytes_sent":        "518",
		"request":                "GET /order/2145 HTTP/1.1",
		"request_method":         "GET",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TextParser.Parse() = %v, want %v", got, want)
	}
}

func BenchmarkParseText(b *testing.B) {
	parser := NewTextParser(`[$time_local] $request_method "$request" $request_length $body_bytes_sent $status $request_time $upstream_response_time`)
	line := `[03/Feb/2021:11:22:33 +0800] GET "GET /order/2145 HTTP/1.1" 123 518 200 0.544 0.543`

	for i := 0; i < b.N; i++ {
		res, err := parser.ParseString(line)
		if err != nil {
			b.Error(err)
		}
		_ = fmt.Sprintf("%v", res)
	}
}
