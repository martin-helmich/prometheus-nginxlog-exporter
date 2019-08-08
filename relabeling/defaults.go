package relabeling

import "github.com/martin-helmich/prometheus-nginxlog-exporter/config"

// DefaultRelabelings are hardcoded relabeling configs that are always there
// and do not need to be explicitly configured

var DefaultRelabelings = []*Relabeling{
        {
                config.RelabelConfig{
                        TargetLabel: "method",
                        SourceValue: "request",
                        Split:       1,

                        WhitelistExists: true,
                        WhitelistMap: map[string]interface{}{
                                "CONNECT": nil,
                                "DELETE":  nil,
                                "GET":     nil,
                                "HEAD":    nil,
                                "OPTIONS": nil,
                                "PATCH":   nil,
                                "POST":    nil,
                                "PUT":     nil,
                                "TRACE":   nil,
                        },
                },
        },
        {
                config.RelabelConfig{
                        TargetLabel: "status",
                        SourceValue: "status",
                },
        },
}
