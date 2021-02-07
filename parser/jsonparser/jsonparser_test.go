package jsonparser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestJsonParse(t *testing.T) {
	parser := NewJsonParser()
	line := `{"time_local":"2021-02-03T11:22:33+08:00","request_length":123,"request_method":"GET","request":"GET /order/2145 HTTP/1.1","body_bytes_sent":518,"status": 200,"request_time":0.544,"upstream_response_time":"0.543"}`

	got, err := parser.ParseString(line)
	if err != nil {
		t.Error(err)
	}
	want := map[string]string{
		"time_local":             "2021-02-03T11:22:33+08:00",
		"request_time":           "0.544",
		"request_length":         "123",
		"upstream_response_time": "0.543",
		"status":                 "200",
		"body_bytes_sent":        "518",
		"request":                "GET /order/2145 HTTP/1.1",
		"request_method":         "GET",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("JsonParser.Parse() = %v, want %v", got, want)
	}
}

func BenchmarkParseJson(b *testing.B) {
	parser := NewJsonParser()
	line := `{"time_local":"2021-02-03T11:22:33+08:00","request_length":123,"request_method":"GET","request":"GET /order/2145 HTTP/1.1","body_bytes_sent":518,"status": 200,"request_time":0.544,"upstream_response_time":"0.543"}`
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		res, err := parser.ParseString(line)
		if err != nil {
			b.Error(err)
		}
		_ = fmt.Sprintf("%v", res)
	}
}
