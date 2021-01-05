Feature: Config file allows relabeling that only apply to the request counter

  Scenario: Labels are added request counter
    Given a running exporter listening with configuration file "test-configuration-only-counter-relabel.hcl"
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET /users HTTP/1.1" 200 612 "-" "curl/7.29.0" "-" 10 10
      172.17.0.1 - foo [23/Jun/2016:16:04:20 +0000] "GET /groups HTTP/1.1" 200 518 "-" "curl/7.29.0" "-" 10 10
      """
    Then the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",path="/groups",status="200",user="foo"}
    And the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",path="/users",status="200",user="other"}

  Scenario: Labels are not add added to size counter
    Given a running exporter listening with configuration file "test-configuration-only-counter-relabel.hcl"
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET /users HTTP/1.1" 200 612 "-" "curl/7.29.0" "-" 10 10
      172.17.0.1 - foo [23/Jun/2016:16:04:20 +0000] "GET /groups HTTP/1.1" 200 518 "-" "curl/7.29.0" "-" 10 10
      """
    Then the exporter should report value 518 for metric nginx_http_response_size_bytes{method="GET",status="200",user="foo"}
    And the exporter should report value 612 for metric nginx_http_response_size_bytes{method="GET",status="200",user="other"}

  Scenario: Labels are not add added to histograms or summaries
    Given a running exporter listening with configuration file "test-configuration-only-counter-relabel.hcl"
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET /users HTTP/1.1" 200 612 "-" "curl/7.29.0" "-" 10 10
      172.17.0.1 - foo [23/Jun/2016:16:04:20 +0000] "GET /groups HTTP/1.1" 200 518 "-" "curl/7.29.0" "-" 10 10
      """
    Then the exporter should report value 1 for metric nginx_http_upstream_time_seconds_hist_count{method="GET",status="200",user="foo"}
    And the exporter should report value 1 for metric nginx_http_upstream_time_seconds_hist_count{method="GET",status="200",user="other"}
    And the exporter should report value 1 for metric nginx_http_upstream_time_seconds_count{method="GET",status="200",user="foo"}
    And the exporter should report value 1 for metric nginx_http_upstream_time_seconds_count{method="GET",status="200",user="other"}
    And the exporter should report value 1 for metric nginx_http_response_time_seconds_hist_count{method="GET",status="200",user="foo"}
    And the exporter should report value 1 for metric nginx_http_response_time_seconds_hist_count{method="GET",status="200",user="other"}
    And the exporter should report value 1 for metric nginx_http_response_time_seconds_count{method="GET",status="200",user="foo"}
    And the exporter should report value 1 for metric nginx_http_response_time_seconds_count{method="GET",status="200",user="other"}