Feature: Message sizes are counted

  Scenario: Response body sizes are counted
    Given a running exporter listening with configuration file "test-config-message-sizes.yaml"
    When the following HTTP request is logged to "access.log"
      """
      $remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" $process_time $bytes_sent $request_length
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 1000 "-" "curl/7.29.0" 300 400
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 2000 "-" "curl/7.29.0" 300 500
      """
    Then the exporter should report value 3000 for metric test_http_response_size_bytes{method="GET",status="200"}

  Scenario: Request body sizes are counted
    Given a running exporter listening with configuration file "test-config-message-sizes.yaml"
    When the following HTTP request is logged to "access.log"
      """
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 1000 "-" "curl/7.29.0" 300 400
      172.17.0.1 - - [23/Jun/2016:16:04:20 +0000] "GET / HTTP/1.1" 200 2000 "-" "curl/7.29.0" 300 500
      """
    Then the exporter should report value 900 for metric test_http_request_size_bytes{method="GET",status="200"}
