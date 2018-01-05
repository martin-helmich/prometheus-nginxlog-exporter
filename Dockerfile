FROM golang:1.9

COPY . /go/src/github.com/martin-helmich/prometheus-nginxlog-exporter
WORKDIR /go/src/github.com/martin-helmich/prometheus-nginxlog-exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prometheus-nginxlog-exporter

FROM scratch

COPY --from=0 /go/src/github.com/martin-helmich/prometheus-nginxlog-exporter/prometheus-nginxlog-exporter /prometheus-nginxlog-exporter

EXPOSE 4040
ENTRYPOINT ["/prometheus-nginxlog-exporter"]
