package relabeling

import (
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
		if _, ok := r.WhitelistMap[sourceValue]; ok {
			return sourceValue, nil
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

	if r.MaxCardinality <= 0 || r.seenValues == nil {
		return sourceValue, nil
	}

	r.seenValuesMu.RLock()
	_, knownValue := r.seenValues[sourceValue]
	r.seenValuesMu.RUnlock()

	if knownValue {
		return sourceValue, nil
	}

	r.seenValuesMu.Lock()
	defer r.seenValuesMu.Unlock()

	if _, knownValue := r.seenValues[sourceValue]; knownValue {
		return sourceValue, nil
	}

	if len(r.seenValues) < r.MaxCardinality {
		r.seenValues[sourceValue] = struct{}{}
		return sourceValue, nil
	}

	if r.OverflowValue != "" {
		return r.OverflowValue, nil
	}

	return "other", nil
}
