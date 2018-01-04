Feature: Config file allows relabel configurations

  Scenario: Labels are added
    Given a running exporter listening with configuration file "test-configuration-relabel.hcl"
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - foo [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",status="200",user="foo"}
    And the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",status="200",user="other"}
