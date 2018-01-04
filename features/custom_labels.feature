Feature: Config file allows custom labels

  Scenario: Labels are added
    Given a running exporter listening with configuration file "test-configuration-labels.hcl"
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 1 for metric nginx_http_response_count_total{foo="bar",method="GET",status="200"}

  Scenario: Labels are added correctly for multiple namespaces
    Given a running exporter listening with configuration file "test-configuration-labels-multi.hcl"
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    And the following HTTP request is logged to "access-two.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "POST / HTTP/1.1" 500 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 1 for metric testone_http_response_count_total{method="GET",status="200",test1="1-1",test2="1-2"}
    And the exporter should report value 1 for metric testtwo_http_response_count_total{method="POST",status="500",test1="2-1",test2="2-2"}

  Scenario: Labels are added for matches routes
    Given a running exporter listening with configuration file "test-configuration-labels-routes.hcl"
    When the following HTTP requests are logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET /users/1234 HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET /users/4321 HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET /profile HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET /foo HTTP/1.1" 200 612 "-" "curl/7.29.0" "-"
      """
    Then the exporter should report value 2 for metric nginx_http_response_count_total{method="GET",request_uri="/users/:id",status="200"}
    And the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",request_uri="/profile",status="200"}
    And the exporter should report value 1 for metric nginx_http_response_count_total{method="GET",request_uri="",status="200"}