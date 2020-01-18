Feature: Various issues reported in the bug tracker remain solved

  Scenario: Issue 44: Patterns with numeric characters are possible
    Given a running exporter listening with configuration file "test-config-issue44.yaml"
    When the following HTTP request is logged to "access.log"
      """
      2018-07-19T23:21:56+03:00	RU	141.8.0.0	https	www.example.com	GET	"/robots.txt"	200	38	"-"	"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)"	-	0.000	0.001
      """
    Then the exporter should report value 1 for metric ht_router_http_response_count_total{method="",status="200"}

  Scenario: Issue 91: Unknown parse error
    Given a running exporter listening with configuration file "test-config-issue91.yaml"
    When the following HTTP request is logged to "access.log"
      """
      28.90.74.145 - - [17/Jan/2020:10:18:11 +0000] "GET /category/finance HTTP/1.1" 200 83 "/category/books?from=20" "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)" 5175 1477 8842
      """
    Then the exporter should report value 1 for metric test_http_response_count_total{method="GET",status="200"}
