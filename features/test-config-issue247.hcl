namespace "nginx_proxy_staging_api" {
  source = {
    files = [
      ".behave-sandbox/access.log"
    ]
  }
  format = "[$time_local] $upstream_cache_status $upstream_status $status - $request_method $scheme $host \"$request_uri\" [Client $remote_addr] [Length $body_bytes_sent] [Gzip $gzip_ratio] [Sent-to $server] \"$http_user_agent\" \"$http_referer\""

  labels {
    app = "acarat_staging_api"
  }

  relabel "method" {
    from = "request_method"
  }
}