port = 4040

namespace "nginx" {
  source_files = [".behave-sandbox/access-1.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
}

namespace "apache" {
  source_files = [".behave-sandbox/access-2.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
}
