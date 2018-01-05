// Example program that reads big nginx file from stdin line by line
// and measure reading time. The file should be big enough, at least 500K lines
package main

import (
	gonx ".."
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

func init() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu + 1)
}

func main() {
	var count int
	format := `$remote_addr - $remote_user [$time_local] "$request" $status ` +
		`$body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for" ` +
		`"$cookie_uid" "$cookie_userid" "$request_time" "$http_host" "$is_ajax" ` +
		`"$uid_got/$uid_set" "$msec" "$geoip_country_code"`
	reader := gonx.NewReader(os.Stdin, format)
	start := time.Now()
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		count++
	}
	duration := time.Since(start)
	fmt.Printf("%v lines readed, it takes %v\n", count, duration)
}
