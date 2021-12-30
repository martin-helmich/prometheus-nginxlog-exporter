listen {
  port = 4040
  
  // "metrics_endpoint" can be used to configure an alternative metrics URL
  // path. Default value is "/metrics".
  //
  //   metrics_endpoint = "/metrics"
}

consul {
  enable = true
  address = "localhost:8500"
  datacenter = "dc1"
  scheme = "http"
  token = ""
  service {
    id = "nginx-exporter"
    name = "nginx-exporter"
    address = "192.168.3.1"
    tags = ["foo", "bar"]
  }
}

namespace "nginx" {
  source = {
    files = [
      "test.log",
      "foo.log",
    ]

    syslog {
      listen_address = "udp://0.0.0.0:5531"
      format = "rfc3164"
      tags = [
        "sometag"
      ]
    }
  }

  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""

  labels {
    app = "magicapp"
    foo = "bar"
  }

  histogram_buckets = [.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10]

  relabel "user" {
    from = "remote_user"
    // whitelist = ["-", "user1", "user2"]

    only_counter = true
  }

  relabel "request_uri" {
    from = "request"
    split = 2

    match "^/users/[0-9]+" {
      replacement = "/users/:id"
    }

    match ".*" {
      replacement = "unknown"
      ignore = "^/api"
    }
  }
}
