package middleware

import "net/http"

func ChainHandlers(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler.ServeHTTP(w, r)
		}
	})
}

/*
TODO
chain 정보
*/
