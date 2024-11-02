FROM golang:1.23

COPY . /work
WORKDIR /work
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o prometheus-nginxlog-exporter

FROM scratch

COPY --from=0 /work/prometheus-nginxlog-exporter /prometheus-nginxlog-exporter

EXPOSE 4040
USER 65532
ENTRYPOINT ["/prometheus-nginxlog-exporter"]
