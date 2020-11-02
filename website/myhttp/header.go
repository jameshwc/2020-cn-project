package myhttp

type Header map[string][]string

func (h Header) Set(key, val string) {
	h[key] = append(h[key], val)
}

func (h Header) Get(key string) string {
	if v, ok := h[key]; ok {
		return v[0]
	}
	return ""
}
