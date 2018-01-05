package gonx

import (
	"testing"
)

func benchLogParsing(b *testing.B, format string, line string) {
	parser := NewParser(format)

	// Ensure the string is in valid format
	_, err := parser.ParseString(line)
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseString(line)
	}
}

func BenchmarkParseSimpleLogRecord(b *testing.B) {
	format := "$remote_addr [$time_local] \"$request\""
	line := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`
	benchLogParsing(b, format, line)
}

func BenchmarkParseLogRecord(b *testing.B) {
	format := `$remote_addr - $remote_user [$time_local] "$request" $status ` +
		`$body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for" ` +
		`"$cookie_uid" "$cookie_userid" "$request_time" "$http_host" "$is_ajax" ` +
		`"$uid_got/$uid_set" "$msec" "$geoip_country_code"`
	line := `**.***.**.*** - - [08/Nov/2013:13:39:18 +0000] ` +
		`"GET /api/internal/v2/item/1?lang=en HTTP/1.1" 200 142 "http://example.com" ` +
		`"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36" ` +
		`"-" "-" "-" "0.084" "example.com" "ajax" "-/-" "1383917958.587" "-"`
	benchLogParsing(b, format, line)
}
