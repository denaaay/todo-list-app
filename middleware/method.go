package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"net/http"
)

func Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			error := model.ErrorResponse{
				Error: "Method is not allowed!",
			}

			jsonError, err := json.Marshal(error)
			if err != nil {
				return
			}
			w.Write([]byte(jsonError))
			return
		}

		next.ServeHTTP(w, r)
	}) // TODO: replace this
}

func Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			error := model.ErrorResponse{
				Error: "Method is not allowed!",
			}
			jsonError, err := json.Marshal(error)
			if err != nil {
				return
			}
			w.Write([]byte(jsonError))
			return
		}

		next.ServeHTTP(w, r)
	}) // TODO: replace this
}

func Delete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			w.WriteHeader(405)
			error := model.ErrorResponse{
				Error: "Method is not allowed!",
			}
			jsonError, err := json.Marshal(error)
			if err != nil {
				return
			}
			w.Write([]byte(jsonError))
			return
		}

		next.ServeHTTP(w, r)
	}) // TODO: replace this
}
