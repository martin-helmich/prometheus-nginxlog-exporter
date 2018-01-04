NGINX-to-Prometheus log file exporter
=====================================

Helper tool that continuously reads an NGINX log file and exports metrics to
[Prometheus](prom).

Usage
-----

You can either use a simple configuration, using command-line flags, or create
a configuration file with a more advanced configuration.

Use the command-line:

    ./nginx-log-exporter \
      -format="<FORMAT>" \
      -listen-port=4040 \
      -namespace=nginx \
      [PATHS-TO-LOGFILES...]

Use the configuration file:

    ./nginx-log-exporter -config-file /path/to/config.hcl

Collected metrics
-----------------

This exporter collects the following metrics. This collector can listen on
multiple log files at once and publish metrics in different namespaces. Each
metric uses the labels `method` (containing the HTTP request method) and
`status` (containing the HTTP status code).

- `<namespace>_http_response_count_total` - The total amount of processed HTTP requests/responses.
- `<namespace>_http_response_size_bytes` - The total amount of transferred content in bytes.
- `<namespace>_http_upstream_time_seconds` - A summary vector of the upstream
  response times in seconds. Logging these needs to be specifically enabled in
  NGINX using the `$upstream_response_time` variable in the log format.
- `<namespace>_http_response_time_seconds` - A summary vector of the total
  response times in seconds. Logging these needs to be specifically enabled in
  NGINX using the `$request_time` variable in the log format.

Additional labels can be configured in the configuration file (see below).

Configuration file
------------------

You can specify a configuration file to read at startup. The configuration file
is expected to be in [HCL](hcl) format. Here's an example file:

```hcl
listen {
  port = 4040
  address = "10.1.2.3"
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

namespace "app-1" {
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  source_files = [
    "/var/log/nginx/app1/access.log"
  ]
  labels {
    app = "application-one"
    environment = "production"
    foo = "bar"
  }
}

namespace "app-2" {
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\" $upstream_response_time"
  source_files = [
    "/var/log/nginx/app2/access.log"
  ]
}
```

Experimental features
---------------------

The exporter contains features that are currently experimental and may change without prior notice.
To use these features, either set the `-enable-experimental` flag or add a `enable_experimental` option
to your configuration file.

### Aggregation by request path

Collecting metrics by the requested resource path has been discussed in #14. Directly adding the requested path as a label is problematic since the set of possible values is infinitely large. For this reason, you can specify a set of `routes` in your configuration file, which is basically a list of regular expressions; if one of these matches a request path, the regular expression will be added as a label to the respective metric:

```hcl
namespace "app-1" {
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  source_files = [
    "/var/log/nginx/app1/access.log"
  ]
  routes = [
    "^/users/[0-9]+",
    "^/profile",
    "^/news"
  ]
}
``` 

Running the collector
---------------------

### Systemd

You can find an example unit file for this service [in this repository](systemd/prometheus-nginxlog-exporter.service). Simply copy the unit file to `/etc/systemd/system`:

    $ wget -O /etc/systemd/system/prometheus-nginxlog-exporter.service https://raw.githubusercontent.com/martin-helmich/prometheus-nginxlog-exporter/master/systemd/prometheus-nginxlog-exporter.service
    $ systemctl enable prometheus-nginxlog-exporter
    $ systemctl start prometheus-nginxlog-exporter

The shipped unit file expects the binary to be located in `/usr/local/bin/prometheus-nginxlog-exporter` and the configuration file in `/etc/prometheus-nginxlog-exporter.hcl`. Adjust to your own needs.

Credits
-------

- [tail](https://github.com/hpcloud/tail), MIT license
- [gonx](https://github.com/satyrius/gonx), MIT license
- [Prometheus Go client library](https://github.com/prometheus/client_golang), Apache License
- [HashiCorp configuration language](hcl), Mozilla Public License

[prom]: https://prometheus.io/
[hcl]: https://github.com/hashicorp/hcl
