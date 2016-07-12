NGINX-to-Prometheus log file exporter
=====================================

Usage
-----

    ./nginx-log-exporter -format="<FORMAT>" -listen-port=4040 -namespace=nginx [PATHS-TO-LOGFILES...]

Credits
-------

- [tail](https://github.com/hpcloud/tail), MIT license
- [gonx](https://github.com/satyrius/gonx), MIT license
- [Prometheus Go client library](https://github.com/prometheus/client_golang), Apache License
