package myhttp

type Handler func(*Context)

type Router struct {
	get  map[string]Handler
	post map[string]Handler
	all  map[string]Handler
}

func NewRouter() *Router {
	get := make(map[string]Handler)
	post := make(map[string]Handler)
	all := make(map[string]Handler)
	return &Router{get, post, all}
}

func (r *Router) GET(path string, f Handler) {
	r.get[path] = f
}

func (r *Router) POST(path string, f Handler) {
	r.post[path] = f
}

func (r *Router) ALL(path string, f Handler) {
	r.all[path] = f
}

func (r *Router) GetHandler(method, path string) (f Handler, ok bool) {
	switch method {
	case "GET":
		if f, ok = r.get[path]; ok {
			return
		}
	case "POST":
		if f, ok = r.post[path]; ok {
			return
		}
	}
	f, ok = r.all[path]
	return
}
