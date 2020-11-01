package server

type ResponseWriter interface {
	Header() map[string][]string
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
