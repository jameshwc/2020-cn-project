package myhttp

import "net/url"

func (r *Request) ParseForm() (err error) {
	if r.PostForm == nil {
		r.PostForm, err = url.ParseQuery(r.Body)
	}
	return
}
