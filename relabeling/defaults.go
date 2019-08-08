package relabeling

import "github.com/martin-helmich/prometheus-nginxlog-exporter/config"

// DefaultRelabelings are hardcoded relabeling configs that are always there
// and do not need to be explicitly configured

func createMethodWhitelistMap() map[string]interface{} {
        whitelistMap := make(map[string]interface{})
        whitelist := []string{
                "HEAD",
                "POST",
                "GET",
                "PUT",
                "DELETE",
                "CONNECT",
                "OPTIONS",
                "TRACE",
                "PATCH",
        }

        for i := range whitelist {
                whitelistMap[whitelist[i]] = nil
        }

        return whitelistMap
}

var DefaultRelabelings = []*Relabeling{
        {
                config.RelabelConfig{
                        TargetLabel: "method",
                        SourceValue: "request",
                        Split:       1,

                        WhitelistExists: true,
                        WhitelistMap:    createMethodWhitelistMap(),
                },
        },
        {
                config.RelabelConfig{
                        TargetLabel: "status",
                        SourceValue: "status",
                },
        },
}
