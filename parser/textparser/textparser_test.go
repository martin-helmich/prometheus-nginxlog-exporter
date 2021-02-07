package textparser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTextParse(t *testing.T) {
	parser := NewTextParser(`[$time_local] $request_method "$request" $request_length $body_bytes_sent $status $request_time $upstream_response_time`)
	line := `[03/Feb/2021:09:58:57 +0800] GET "GET /gateway/worksheet/worksheet/getWorksheetById?id=16523 HTTP/1.1" 123 519 200 0.544 0.543`
	got, err := parser.ParseString(line)
	if err != nil {
		t.Error(err)
	}
	want := map[string]string{
		"time_local":             "03/Feb/2021:09:58:57 +0800",
		"request_time":           "0.544",
		"request_length":         "123",
		"upstream_response_time": "0.543",
		"status":                 "200",
		"body_bytes_sent":        "519",
		"request":                "GET /gateway/worksheet/worksheet/getWorksheetById?id=16523 HTTP/1.1",
		"request_method":         "GET",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TextParser.Parse() = %v, want %v", got, want)
	}
}

func BenchmarkParseText(b *testing.B) {
	parser := NewTextParser(`[$time_local] $request_method "$request" $request_length $body_bytes_sent $status $request_time $upstream_response_time`)
	line := `[03/Feb/2021:09:58:57 +0800] GET "GET /gateway/worksheet/worksheet/getWorksheetById?id=16523 HTTP/1.1" 123 519 200 0.544 0.543`
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		res, err := parser.ParseString(line)
		if err != nil {
			b.Error(err)
		}
		_ = fmt.Sprintf("%v", res)
	}
}
