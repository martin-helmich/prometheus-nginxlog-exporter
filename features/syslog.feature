Feature: Can read log entries from syslog

  Scenario: Read from syslog
    Given a running exporter listening with configuration file "test-config-syslog.yaml"
    When the following HTTP request is logged to syslog on port 1234
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",status="200"}
