port = 4040
enable_experimental = true

namespace "nginx" {
  source_files = [".behave-sandbox/access.log"]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  #routes = [
  #  "^/users/[0-9]+",
  #  "^/profile"
  #]

  relabel "request_uri" {
    from = "request"
    split = 2

    match "^/users/[0-9]+" {
      replacement = "/users/:id"
    }

    match "^/profile" {
      replacement = "/profile"
    }

    match {
      replacement = ""
    }
  }
}
