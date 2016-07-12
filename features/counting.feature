Feature: Requests are counted

  Scenario: Single request is counted
    Given a running exporter listening on "access.log" with default format
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",status="200"}

  Scenario: Multiple requests are counted
    Given a running exporter listening on "access.log" with default format
    When the following HTTP requests are logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 3 for metric nginx_http_response_count_total{method="GET",status="200"}

  Scenario: Requests are grouped by method
    Given a running exporter listening on "access.log" with default format
    When the following HTTP requests are logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "POST / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 2 for metric nginx_http_response_count_total{method="GET",status="200"}
    And the exporter should report value 1 for metric nginx_http_response_count_total{method="POST",status="200"}

  Scenario: Requests are grouped by status
    Given a running exporter listening on "access.log" with default format
    When the following HTTP requests are logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 500 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 2 for metric nginx_http_response_count_total{method="GET",status="200"}
    And the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",status="500"}