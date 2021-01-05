port = 4040
enable_experimental = true

namespace "nginx" {
  source {
    files = [".behave-sandbox/access.log"]
  }
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\" $request_time $upstream_response_time"

  relabel "user" {
    from = "remote_user"
    whitelist = ["foo", "bar"]
  }

  relabel "path" {
    from = "request"
    split = 2

    only_counter = true
  }
}
