Feature: Various issues reported in the bug tracker remain solved

  Scenario: Issue 44: Patterns with numeric characters are possible
    Given a running exporter listening with configuration file "test-config-issue44.yaml"
    When the following HTTP request is logged to "access.log"
      """
      2018-07-19T23:21:56+03:00	RU	141.8.0.0	https	www.example.com	GET	"/robots.txt"	200	38	"-"	"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)"	-	0.000	0.001
      """
    Then the exporter should report value 1 for metric ht_router_http_response_count_total{method="",status="200"}

  Scenario: Issue 90: Unknown parse error
    Given a running exporter listening with configuration file "test-config-issue90.yaml"
    When the following HTTP request is logged to "access.log"
      """
      10.rr.ii.yy - - [25/Dec/2019:08:06:42 +0300] "GET /offsets/topic/xxx-xxx/partition/0 HTTP/1.1" 200 96 "-" "Java/11.0.2" "-"
      10.rr.ii.yy - - [25/Dec/2019:08:06:42 +0300] "GET /api/v2/topics/xxx-xxxx/partitions/0/offsets HTTP/1.1" 200 68 "-" "Java/11.0.2" "-"
      10.rr.ii.yy - - [25/Dec/2019:08:06:42 +0300] "GET /api/v2/topics/xxxxxx/partitions/2/messages?offset=96&count=10 HTTP/1.1" 200 41 "-" "Java/11.0.2" "-"
      """
    Then the exporter should report value 3 for metric test_http_response_count_total{method="GET",status="200"}

  Scenario: Issue 91: Unknown parse error
    Given a running exporter listening with configuration file "test-config-issue91.yaml"
    When the following HTTP request is logged to "access.log"
      """
      28.90.74.145 - - [17/Jan/2020:10:18:11 +0000] "GET /category/finance HTTP/1.1" 200 83 "/category/books?from=20" "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)" 5175 1477 8842
      """
    Then the exporter should report value 1 for metric test_http_response_count_total{method="GET",status="200"}

  Scenario: Issue 217: Multiple response time values
    Given a running exporter listening with configuration file "test-config-issue217.yaml"
    When the following HTTP request is logged to "access.log"
      """
      www.example.com 10.0.0.1 - - [27/Oct/2021:06:00:14 +0200] "GET /websocket HTTP/1.1" 200 552 "-" "SomeUserAgentString" 0.000 "1.000, 0.250, 2.000, 0.750" . TLSv1.3/TLS_AES_256_GCM_SHA384 1234567890 www.example.com
      www.example.com 10.0.0.2 - - [27/Oct/2021:06:00:15 +0200] "GET /websocket HTTP/1.1" 200 552 "-" "AnotherUserAgent" 0.000 "1.000, 0.250, 2.000, 0.750" . TLSv1.3/TLS_AES_256_GCM_SHA384 1234567890 www.example.com
      www.example.com 10.0.0.3 - - [27/Oct/2021:06:00:15 +0200] "GET /websocket HTTP/1.1" 200 150 "-" "SomeUserAgentString" 0.001 "1.000, 0.250, 2.000, 0.750" . TLSv1.3/TLS_AES_256_GCM_SHA384 1234567890 www.example.com
      """
    Then the exporter should report value 12 for metric test_http_upstream_time_seconds_sum{method="GET",status="200"}