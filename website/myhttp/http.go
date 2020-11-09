package myhttp

func Handle(conn Conn, r *Router) {
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
	defer c.Close()
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
