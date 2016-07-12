Feature: Upstream response times are summarized

  Scenario: Single request is summarized
    Given a running exporter listening on "access.log" with upstream-time format
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-" 10
      """
    Then the exporter should report value 10 for metric nginx_http_upstream_time_seconds{method="GET",status="200",quantile="0.5"}

  Scenario: .5 quantile is computed
    Given a running exporter listening on "access.log" with upstream-time format
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-" 10
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-" 20
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-" 30
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-" 40
      """
    Then the exporter should report value 20 for metric nginx_http_upstream_time_seconds{method="GET",status="200",quantile="0.5"}
