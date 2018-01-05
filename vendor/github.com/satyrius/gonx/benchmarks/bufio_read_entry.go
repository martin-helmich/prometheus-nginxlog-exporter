// Example program that reads big nginx file from stdin line by line
// and measure reading time. The file should be big enough, at least 500K lines
package main

import (
	gonx ".."
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var count int
	format := `$remote_addr - $remote_user [$time_local] "$request" $status ` +
		`$body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for" ` +
		`"$cookie_uid" "$cookie_userid" "$request_time" "$http_host" "$is_ajax" ` +
		`"$uid_got/$uid_set" "$msec" "$geoip_country_code"`
	parser := gonx.NewParser(format)
	start := time.Now()
	for scanner.Scan() {
		// A dummy action, jest read line by line
		parser.ParseString(scanner.Text())
		count++
	}
	duration := time.Since(start)
	fmt.Printf("%v lines readed, it takes %v\n", count, duration)
}
