FROM scratch

COPY prometheus-nginxlog-exporter /prometheus-nginxlog-exporter

EXPOSE 4040
ENTRYPOINT ["/prometheus-nginxlog-exporter"]
