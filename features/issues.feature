Feature: Various issues reported in the bug tracker remain solved

  Scenario: Patterns with numeric characters are possible
    Given a running exporter listening with configuration file "test-config-issue44.yaml"
    When the following HTTP request is logged to "access.log"
      """
      2018-07-19T23:21:56+03:00	RU	141.8.0.0	https	www.example.com	GET	"/robots.txt"	200	38	"-"	"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)"	-	0.000	0.001
      """
    Then the exporter should report value 1 for metric ht_router_http_response_count_total{method="",status="200"}
