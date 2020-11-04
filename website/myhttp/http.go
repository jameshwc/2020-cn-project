package myhttp

func Handle(conn Conn, r *Router) {
	defer conn.Close() // TODO: move it to method?
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return
	}
	request, err := parseRequest(buf)
	if err != nil {
		return
	}
	c := NewContext(conn, request)
	serve(c, r)
}

func serve(c *Context, r *Router) {
	beforeServe(c)
	f, ok := r.GetHandler(c.Request.Method, c.Request.URL.Path)
	if !ok {
		c.NotFound()
		return
	}
	f(c)
	afterServe(c)
}
