Feature: Config file allows multiple namespaces

  Scenario: Single request is counted
    Given a running exporter listening with default configuration file
    When the following HTTP request is logged to "access-1.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",status="200"}
