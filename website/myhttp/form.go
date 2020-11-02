package myhttp

import "net/url"

func (r *Request) ParseFrom() (err error) {
	if r.PostForm == nil {
		r.PostForm, err = url.ParseQuery(r.Body)
	}
	return
}
