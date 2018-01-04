listen {
  port = 4040
}

consul {
  enable = true
  address = "localhost:8500"
  service {
    id = "nginx-exporter"
    name = "nginx-exporter"
    datacenter = "dc1"
    scheme = "http"
    token = ""
    tags = ["foo", "bar"]
  }
}

namespace "nginx" {
  source_files = [
    "test.log",
    "foo.log"
  ]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""

  labels {
    app = "magicapp"
    foo = "bar"
  }

  relabel "user" {
    from = "remote_user"
    // whitelist = ["-", "user1", "user2"]
  }

  routes = [
    "^/users/[0-9]+"
  ]
}
