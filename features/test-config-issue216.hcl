listen {
  port = 4040
}

namespace "nginx" {
  source = {
    files = [
      ".behave-sandbox/access.log"
    ]
  }

  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\" $request_time $upstream_response_time"

  relabel "status" {
    from = "status"
    only_counter = true
  }
  relabel "method" {
    from = "request"
    only_counter = true
    split = 1
    whitelist = [ "GET" ]
  }
}
