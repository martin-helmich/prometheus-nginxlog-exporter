package relabeling

import (
	"path/filepath"
	"strings"
)

// Map maps a sourceValue from the access log line according to the relabeling
// config (matching against whitelists, regular expressions etc.)
func (r *Relabeling) Map(sourceValue string) (string, error) {
	if r.Split > 0 {
		separator := r.Separator
		if separator == "" {
			separator = " "
		}

		values := strings.Split(sourceValue, separator)

		if len(values) >= r.Split {
			sourceValue = values[r.Split-1]
		} else {
			sourceValue = ""
		}
	}

	if r.WhitelistExists {
		for pattern := range r.WhitelistMap {
			if strings.Contains(pattern, "*") {
				matched, err := filepath.Match(pattern, sourceValue)
				if err != nil {
					return "other", nil // handle invalid pattern error if it occurs
				}
				if matched {
					return sourceValue, nil
				}
			} else {
				if pattern == sourceValue {
					return sourceValue, nil
				}
			}
		}

		return "other", nil
	}

	if len(r.Matches) > 0 {
		replacement := ""
		for i := range r.Matches {
			if r.Matches[i].CompiledRegexp.MatchString(sourceValue) {
				replacement = r.Matches[i].CompiledRegexp.ReplaceAllString(sourceValue, r.Matches[i].Replacement)
				break
			}
		}
		sourceValue = replacement
	}

	return sourceValue, nil
}
