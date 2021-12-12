package utils

import "net/http"

func SetWriterHeaderJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func AddSetJsonHeaderMw(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		SetWriterHeaderJson(w)
		handler(w, r)
	}
}

func ChainMiddlewares(root func(w http.ResponseWriter, r *http.Request), fns ...func(func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	ret := root
	for _, fn := range fns {
		ret = fn(root)
	}

	return ret
}
