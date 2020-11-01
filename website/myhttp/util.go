package myhttp

import "strings"

func header2string(m map[string][]string) (s string) {
	for k, v := range m {
		s += k + ":" + strings.Join(v, ",")
		s += "<br>"
	}
	return s
}
