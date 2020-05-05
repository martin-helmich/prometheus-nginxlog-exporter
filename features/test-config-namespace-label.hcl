port = 4040

namespace "testone" {
  source_files = [".behave-sandbox/access.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  labels {
    test1 = "1-1"
    test2 = "1-2"
  }

  metrics_override = { prefix = "test" }
  namespace_label = "vhost"
}

namespace "testtwo" {
  source_files = [".behave-sandbox/access-two.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  labels {
    test1 = "2-1"
    test2 = "2-2"
  }
  
  metrics_override = { prefix = "test" }
  namespace_label = "vhost"
}
