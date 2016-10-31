NGINX-to-Prometheus log file exporter
=====================================

Helper tool that continuously reads an NGINX log file and exports metrics to
[Prometheus](prom).

Usage
-----

You can either use a simple configuration, using command-line flags, or create
a configuration file with a more advanced configration.

Use the command-line:

    ./nginx-log-exporter \
      -format="<FORMAT>" \
      -listen-port=4040 \
      -namespace=nginx \
      [PATHS-TO-LOGFILES...]

Use the configuration file:

    ./nginx-log-exporter -config-file /path/to/config.hcl

Configuration file
------------------

You can specify a configuration file to read at startup. The configuration file
is expected to be in [HCL](hcl) format. Here's an example file:

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
      format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\" $upstream_response_time
      source_file = [
        "/var/log/nginx/app2/access.log"
      ]
    }

Credits
-------

- [tail](https://github.com/hpcloud/tail), MIT license
- [gonx](https://github.com/satyrius/gonx), MIT license
- [Prometheus Go client library](https://github.com/prometheus/client_golang), Apache License
- [HashiCorp configuration language](hcl), Mozilla Public License

[prom]: https://prometheus.io/
[hcl]: https://github.com/hashicorp/hcl
